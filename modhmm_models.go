/* Copyright (C) 2018 Philipp Benner
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

/* -------------------------------------------------------------------------- */

import   "fmt"
import   "math"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/generic"
import   "github.com/pbenner/autodiff/statistics/vectorEstimator"
import   "github.com/pbenner/autodiff/statistics/matrixDistribution"
import   "github.com/pbenner/autodiff/statistics/matrixEstimator"

import . "github.com/pbenner/autodiff"

import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

func init() {
  VectorPdfRegistry["vector:probability distribution"] = new(EmissionDistribution)
  MatrixPdfRegistry["ModHmm"] = new(ModHmm)
}

/* -------------------------------------------------------------------------- */

func getRGBMap() map[string]string {
  m := make(map[string]string)
  m["PA"   ] = "0,100,0"
  m["EA"   ] = "30,144,255"
  m["EA:tr"] = "30,144,255"
  m["PR"   ] = "112,128,144"
  m["PR:tr"] = "112,128,144"
  m["BI"   ] = "178,34,34"
  m["BI:tr"] = "178,34,34"
  m["R1"   ] = "255,69,0"
  m["R2"   ] = "255,69,0"
  m["TR"   ] = "255,215,0"
  m["CL"   ] = "255,0,255"
  m["NS"   ] = ""
  return m
}

/* -------------------------------------------------------------------------- */

type ModHmm struct {
  matrixDistribution.Hmm
  StateNames []string
}

func (obj *ModHmm) ImportConfig(config ConfigDistribution, t ScalarType) error {
  if len(config.Distributions) != 1 {
    return fmt.Errorf("invalid config")
  }
  if err := obj.Hmm.ImportConfig(config.Distributions[0], t); err != nil {
    return err
  }
  if s, ok := config.GetNamedParametersAsStrings("StateNames"); !ok {
    return fmt.Errorf("invalid config")
  } else {
    obj.StateNames = s
  }
  return nil
}

func (obj *ModHmm) ExportConfig() ConfigDistribution {

  parameters := struct{StateNames []string}{}
  parameters.StateNames = obj.StateNames

  return NewConfigDistribution("ModHmm", parameters, obj.Hmm.ExportConfig())
}

/* emission distribution
 * -------------------------------------------------------------------------- */

type EmissionDistribution struct {
  i int
  n int
}

func (obj *EmissionDistribution) CloneVectorPdf() VectorPdf {
  return &EmissionDistribution{obj.i, obj.n}
}

func (obj *EmissionDistribution) LogPdf(r Scalar, x ConstVector) error {
  r.SetFloat64(math.Log(x.Float64At(obj.i)))
  if math.IsNaN(r.GetFloat64()) {
    panic("internal error")
  }
  return nil
}

func (obj *EmissionDistribution) Dim() int {
  return obj.n
}

func (obj *EmissionDistribution) ScalarType() ScalarType {
  return Float64Type
}

func (obj *EmissionDistribution) GetParameters() Vector {
  p := NullDenseVector(Float64Type, 2)
  p.At(0).SetFloat64(float64(obj.i))
  p.At(1).SetFloat64(float64(obj.n))
  return p
}

func (obj *EmissionDistribution) SetParameters(parameters Vector) error {
  obj.i = int(parameters.At(0).GetFloat64())
  obj.n = int(parameters.At(1).GetFloat64())
  return nil
}

func (obj *EmissionDistribution) ImportConfig(config ConfigDistribution, t ScalarType) error {
  if parameters, ok := config.GetParametersAsFloats(); !ok {
    return fmt.Errorf("invalid config file")
  } else {
    obj.i = int(parameters[0])
    obj.n = int(parameters[1])
    return nil
  }
}

func (obj *EmissionDistribution) ExportConfig() ConfigDistribution {
  return NewConfigDistribution("vector:probability distribution", obj.GetParameters())
}

/* -------------------------------------------------------------------------- */

func getModHmmDenseEstimator(config ConfigModHmm) (*matrixEstimator.HmmEstimator, []string) {
  stateNames := []string{
    "PA", "EA", "BI", "TR", "R1", "R2", "CL", "NS"}

  n := 8

  pi := NullDenseVector(Float64Type, n)
  tr := NullDenseMatrix(Float64Type, n, n)
  pi.Map(func(x Scalar) { x.SetFloat64(1.0) })
  tr.Map(func(x Scalar) { x.SetFloat64(1.0) })

  // emissions
  estimators := make([]VectorEstimator, n)
  for i := 0; i < n; i++ {
    estimators[i] = vectorEstimator.NilEstimator{&EmissionDistribution{i, n}}
  }

  if estimator, err := matrixEstimator.NewHmmEstimator(pi, tr, nil, nil, nil, estimators, 1e-0, -1); err != nil {
    panic(err)
  } else {
    estimator.ChunkSize = 10000
    estimator.OptimizeEmissions = false
    switch config.Verbose {
    case 0 : estimator.Verbose = 0
    case 1 : estimator.Verbose = 1
    case 2 : estimator.Verbose = 1
    default: estimator.Verbose = 2
    }
    return estimator, stateNames
  }
}

/* -------------------------------------------------------------------------- */

func getModHmmDefaultEstimator(config ConfigModHmm) (*matrixEstimator.HmmEstimator, []string) {
  const jEA   =  0 // enhancer active
  const jPR   =  1 // enhancer active
  const jT3   =  2 // transcribed
  const jR1   =  3 // repressed h3k27me3
  const jR2   =  4 // repressed h3k9me3
  const jNS   =  5 // no signal
  const jCL   =  6 // control
  const jPA1  =  7 // promoter active
  const jPA2  =  8 // promoter active
  const jBI   =  9 // bivalent
  const jT1   = 10 // transcribed
  const jT2   = 11 // transcribed
  const jEAt1 = 12 // enhancer active
  const jEAt2 = 13 // enhancer active
  const jBIt1 = 14 // bivalent
  const jBIt2 = 15 // bivalent
  const jPRt1 = 16 // bivalent
  const jPRt2 = 17 // bivalent

  stateNames := []string{
    "EA", "PR", "TR", "R1", "R2", "NS", "CL", "PA", "PA", "BI", "TR", "TR", "EA:tr", "EA:tr", "BI:tr", "BI:tr", "PR:tr", "PR:tr"}

  n :=  9
  m := 18

  stateMap := make([]int, m)
  stateMap[jEA]   = iEA
  stateMap[jEAt1] = iEA
  stateMap[jEAt2] = iEA
  stateMap[jPR]   = iPR
  stateMap[jPRt1] = iPR
  stateMap[jPRt2] = iPR
  stateMap[jPA1]  = iPA
  stateMap[jPA2]  = iPA
  stateMap[jBI]   = iBI
  stateMap[jBIt1] = iBI
  stateMap[jBIt2] = iBI
  stateMap[jT1]   = iTR
  stateMap[jT2]   = iTR
  stateMap[jT3]   = iTR
  stateMap[jR1]   = iR1
  stateMap[jR2]   = iR2
  stateMap[jNS]   = iNS
  stateMap[jCL]   = iCL

  pi := NullDenseVector(Float64Type, m)
  tr := NullDenseMatrix(Float64Type, m, m)
  pi.Map(func(x Scalar) { x.SetFloat64(1.0) })

  // allow self-transitions for all states
  for i := 0; i < m; i++ {
    tr.At(i,i).SetFloat64(1.0)
  }
  // enhancer active
  tr.At(jEA  ,jCL  ).SetFloat64(1.0)
  tr.At(jEA  ,jNS  ).SetFloat64(1.0)
  tr.At(jEA  ,jR1  ).SetFloat64(1.0)
  tr.At(jEA  ,jR2  ).SetFloat64(1.0)
  // bivalent
  tr.At(jBI  ,jCL  ).SetFloat64(1.0)
  tr.At(jBI  ,jNS  ).SetFloat64(1.0)
  tr.At(jBI  ,jR1  ).SetFloat64(1.0)
  tr.At(jBI  ,jR2  ).SetFloat64(1.0)
  // primed
  tr.At(jPR  ,jCL  ).SetFloat64(1.0)
  tr.At(jPR  ,jNS  ).SetFloat64(1.0)
  tr.At(jPR  ,jR1  ).SetFloat64(1.0)
  tr.At(jPR  ,jR2  ).SetFloat64(1.0)
  // transcribed (low)
  tr.At(jT3  ,jCL  ).SetFloat64(1.0)
  tr.At(jT3  ,jNS  ).SetFloat64(1.0)
  tr.At(jT3  ,jR1  ).SetFloat64(1.0)
  tr.At(jT3  ,jR2  ).SetFloat64(1.0)
  // no signal
  tr.At(jNS  ,jCL  ).SetFloat64(1.0)
  tr.At(jNS  ,jEA  ).SetFloat64(1.0)
  tr.At(jNS  ,jR1  ).SetFloat64(1.0)
  tr.At(jNS  ,jR2  ).SetFloat64(1.0)
  tr.At(jNS  ,jPA1 ).SetFloat64(1.0)
  tr.At(jNS  ,jBI  ).SetFloat64(1.0)
  tr.At(jNS  ,jPR  ).SetFloat64(1.0)
  tr.At(jNS  ,jT1  ).SetFloat64(1.0)
  tr.At(jNS  ,jT3  ).SetFloat64(1.0)
  // control
  tr.At(jCL  ,jEA  ).SetFloat64(1.0)
  tr.At(jCL  ,jNS  ).SetFloat64(1.0)
  tr.At(jCL  ,jR1  ).SetFloat64(1.0)
  tr.At(jCL  ,jR2  ).SetFloat64(1.0)
  tr.At(jCL  ,jPA1 ).SetFloat64(1.0)
  tr.At(jCL  ,jBI  ).SetFloat64(1.0)
  tr.At(jCL  ,jPR  ).SetFloat64(1.0)
  tr.At(jCL  ,jT1  ).SetFloat64(1.0)
  tr.At(jCL  ,jT3  ).SetFloat64(1.0)
  // repressed 1
  tr.At(jR1  ,jCL  ).SetFloat64(1.0)
  tr.At(jR1  ,jEA  ).SetFloat64(1.0)
  tr.At(jR1  ,jNS  ).SetFloat64(1.0)
  tr.At(jR1  ,jR2  ).SetFloat64(1.0)
  tr.At(jR1  ,jPA1 ).SetFloat64(1.0)
  tr.At(jR1  ,jBI  ).SetFloat64(1.0)
  tr.At(jR1  ,jPR  ).SetFloat64(1.0)
  tr.At(jR1  ,jT1  ).SetFloat64(1.0)
  tr.At(jR1  ,jT3  ).SetFloat64(1.0)
  // repressed 2
  tr.At(jR2  ,jCL  ).SetFloat64(1.0)
  tr.At(jR2  ,jEA  ).SetFloat64(1.0)
  tr.At(jR2  ,jNS  ).SetFloat64(1.0)
  tr.At(jR2  ,jR1  ).SetFloat64(1.0)
  tr.At(jR2  ,jPA1 ).SetFloat64(1.0)
  tr.At(jR2  ,jBI  ).SetFloat64(1.0)
  tr.At(jR2  ,jPR  ).SetFloat64(1.0)
  tr.At(jR2  ,jT1  ).SetFloat64(1.0)
  tr.At(jR2  ,jT3  ).SetFloat64(1.0)
  // promoter active 1
  tr.At(jPA1 ,jT2  ).SetFloat64(1.0)
  // promoter active 2
  tr.At(jPA2 ,jT2  ).SetFloat64(1.0)
  tr.At(jPA2 ,jCL  ).SetFloat64(1.0)
  tr.At(jPA2 ,jNS  ).SetFloat64(1.0)
  tr.At(jPA2 ,jR1  ).SetFloat64(1.0)
  tr.At(jPA2 ,jR2  ).SetFloat64(1.0)
  // transcribed 1
  tr.At(jT1  ,jPA2 ).SetFloat64(1.0)
  tr.At(jT1  ,jEAt1).SetFloat64(1.0)
  tr.At(jT1  ,jBIt1).SetFloat64(1.0)
  tr.At(jT1  ,jPRt1).SetFloat64(1.0)
  // transcribed 2
  tr.At(jT2  ,jPA2 ).SetFloat64(1.0)
  tr.At(jT2  ,jEAt2).SetFloat64(1.0)
  tr.At(jT2  ,jBIt2).SetFloat64(1.0)
  tr.At(jT2  ,jPRt2).SetFloat64(1.0)
  tr.At(jT2  ,jCL  ).SetFloat64(1.0)
  tr.At(jT2  ,jNS  ).SetFloat64(1.0)
  tr.At(jT2  ,jR1  ).SetFloat64(1.0)
  tr.At(jT2  ,jR2  ).SetFloat64(1.0)
  // ea/bi/pr transcribed
  tr.At(jEAt1,jT1  ).SetFloat64(1.0)
  tr.At(jBIt1,jT1  ).SetFloat64(1.0)
  tr.At(jPRt1,jT1  ).SetFloat64(1.0)
  // ea/bi/pr transcribed
  tr.At(jEAt2,jT2  ).SetFloat64(1.0)
  tr.At(jBIt2,jT2  ).SetFloat64(1.0)
  tr.At(jPRt2,jT2  ).SetFloat64(1.0)

  constraints := []generic.EqualityConstraint{}
  if config.ModelUnconstrained {
    printStderr(config, 2, "Implementing default model with unconstrained transition matrix\n")
  } else {
    printStderr(config, 2, "Implementing default model with constrained transition matrix\n")
    for i := 0; i < m; i++ {
      constraint := generic.EqualityConstraint{}
      for j := 0; j < m; j++ {
        if i == j {
          continue
        }
        if tr.ConstAt(i, j).GetFloat64() != 0 {
          constraint = append(constraint, [2]int{i,j})
        }
      }
      constraints = append(constraints, constraint)
    }
    // constrain self-transitions
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jPA1, jPA1}, [2]int{jPA2, jPA2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jEA, jEA}, [2]int{jEAt1, jEAt1}, [2]int{jEAt2, jEAt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jBI, jBI}, [2]int{jBIt1, jBIt1}, [2]int{jBIt2, jBIt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jPR, jPR}, [2]int{jPRt1, jPRt1}, [2]int{jPRt2, jPRt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jT1}, [2]int{jT2, jT2}, [2]int{jT3, jT3}})
  }
  // emissions
  estimators := make([]VectorEstimator, n)
  for i := 0; i < n; i++ {
    estimators[i] = vectorEstimator.NilEstimator{&EmissionDistribution{i, n}}
  }

  if estimator, err := matrixEstimator.NewConstrainedHmmEstimator(pi, tr, stateMap, nil, nil, constraints, estimators, 1e-0, -1); err != nil {
    panic(err)
  } else {
    estimator.ChunkSize = 10000
    estimator.OptimizeEmissions = false
    switch config.Verbose {
    case 0 : estimator.Verbose = 0
    case 1 : estimator.Verbose = 1
    case 2 : estimator.Verbose = 1
    default: estimator.Verbose = 2
    }
    return estimator, stateNames
  }
}
