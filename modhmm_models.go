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
import   "strings"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/generic"
import   "github.com/pbenner/autodiff/statistics/vectorEstimator"
import   "github.com/pbenner/autodiff/statistics/matrixDistribution"
import   "github.com/pbenner/autodiff/statistics/matrixEstimator"

import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

func init() {
  VectorPdfRegistry["vector:probability distribution"] = new(EmissionDistribution)
  MatrixPdfRegistry["ModHmm"] = new(ModHmm)
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
  r.Set(x.ConstAt(obj.i))
  return nil
}

func (obj *EmissionDistribution) Dim() int {
  return obj.n
}

func (obj *EmissionDistribution) ScalarType() ScalarType {
  return BareRealType
}

func (obj *EmissionDistribution) GetParameters() Vector {
  p := NullVector(BareRealType, 2)
  p.At(0).SetValue(float64(obj.i))
  p.At(1).SetValue(float64(obj.n))
  return p
}

func (obj *EmissionDistribution) SetParameters(parameters Vector) error {
  obj.i = int(parameters.At(0).GetValue())
  obj.n = int(parameters.At(1).GetValue())
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
    "PA", "PB", "EA", "EP", "TR", "TL", "R1", "R2", "CL", "NS"}

  n := 10

  pi := NullVector(BareRealType, n)
  tr := NullMatrix(BareRealType, n, n)
  pi.Map(func(x Scalar) { x.SetValue(1.0) })
  tr.Map(func(x Scalar) { x.SetValue(1.0) })

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
  const jEP   =  1 // enhancer poised
  const jTL   =  2 // transcribed (low)
  const jR1   =  3 // repressed h3k27me3
  const jR2   =  4 // repressed h3k9me3
  const jNS   =  5 // no signal
  const jCL   =  6 // control
  const jPA1  =  7 // promoter active
  const jPA2  =  8 // promoter active
  const jPB   =  9 // promoter bivalent
  const jT1   = 10 // transcribed
  const jT2   = 11 // transcribed
  const jEAt1 = 12 // enhancer active
  const jEPt1 = 13 // enhancer poised
  const jEAt2 = 14 // enhancer active
  const jEPt2 = 15 // enhancer poised
  const jPBt1 = 16 // promoter bivalent
  const jPBt2 = 17 // promoter bivalent

  stateNames := []string{
    "EA", "EP", "TL", "R1", "R2", "NS", "CL", "PA", "PA", "PB", "TR", "TR", "EA:tr", "EP:tr", "EA:tr", "EP:tr", "PB:tr", "PB:tr"}

  n := 10
  m := 18

  stateMap := make([]int, m)
  stateMap[jEA]   = iEA
  stateMap[jEP]   = iEP
  stateMap[jEAt1] = iEA
  stateMap[jEPt1] = iEP
  stateMap[jEAt2] = iEA
  stateMap[jEPt2] = iEP
  stateMap[jPA1]  = iPA
  stateMap[jPA2]  = iPA
  stateMap[jPB]   = iPB
  stateMap[jPBt1] = iPB
  stateMap[jPBt2] = iPB
  stateMap[jT1]   = iTR
  stateMap[jT2]   = iTR
  stateMap[jTL]   = iTL
  stateMap[jR1]   = iR1
  stateMap[jR2]   = iR2
  stateMap[jNS]   = iNS
  stateMap[jCL]   = iCL

  pi := NullVector(BareRealType, m)
  tr := NullMatrix(BareRealType, m, m)
  pi.Map(func(x Scalar) { x.SetValue(1.0) })

  // allow self-transitions for all states
  for i := 0; i < m; i++ {
    tr.At(i,i).SetValue(1.0)
  }
  // enhancer active
  tr.At(jEA  ,jCL  ).SetValue(1.0)
  tr.At(jEA  ,jNS  ).SetValue(1.0)
  tr.At(jEA  ,jR1  ).SetValue(1.0)
  tr.At(jEA  ,jR2  ).SetValue(1.0)
  tr.At(jEA  ,jTL  ).SetValue(1.0)
  // enhancer poised
  tr.At(jEP  ,jCL  ).SetValue(1.0)
  tr.At(jEP  ,jNS  ).SetValue(1.0)
  tr.At(jEP  ,jR1  ).SetValue(1.0)
  tr.At(jEP  ,jR2  ).SetValue(1.0)
  tr.At(jEP  ,jTL  ).SetValue(1.0)
  // promoter bivalent
  tr.At(jPB  ,jCL  ).SetValue(1.0)
  tr.At(jPB  ,jNS  ).SetValue(1.0)
  tr.At(jPB  ,jR1  ).SetValue(1.0)
  tr.At(jPB  ,jR2  ).SetValue(1.0)
  tr.At(jPB  ,jTL  ).SetValue(1.0)
  // transcribed (low)
  tr.At(jTL  ,jCL  ).SetValue(1.0)
  tr.At(jTL  ,jEA  ).SetValue(1.0)
  tr.At(jTL  ,jEP  ).SetValue(1.0)
  tr.At(jTL  ,jNS  ).SetValue(1.0)
  tr.At(jTL  ,jR1  ).SetValue(1.0)
  tr.At(jTL  ,jR2  ).SetValue(1.0)
  tr.At(jTL  ,jPA1 ).SetValue(1.0)
  tr.At(jTL  ,jPB  ).SetValue(1.0)
  tr.At(jTL  ,jT1  ).SetValue(1.0)
  // no signal
  tr.At(jNS  ,jCL  ).SetValue(1.0)
  tr.At(jNS  ,jEA  ).SetValue(1.0)
  tr.At(jNS  ,jEP  ).SetValue(1.0)
  tr.At(jNS  ,jR1  ).SetValue(1.0)
  tr.At(jNS  ,jR2  ).SetValue(1.0)
  tr.At(jNS  ,jPA1 ).SetValue(1.0)
  tr.At(jNS  ,jPB  ).SetValue(1.0)
  tr.At(jNS  ,jT1  ).SetValue(1.0)
  tr.At(jNS  ,jTL  ).SetValue(1.0)
  // control
  tr.At(jCL  ,jEA  ).SetValue(1.0)
  tr.At(jCL  ,jEP  ).SetValue(1.0)
  tr.At(jCL  ,jNS  ).SetValue(1.0)
  tr.At(jCL  ,jR1  ).SetValue(1.0)
  tr.At(jCL  ,jR2  ).SetValue(1.0)
  tr.At(jCL  ,jPA1 ).SetValue(1.0)
  tr.At(jCL  ,jPB  ).SetValue(1.0)
  tr.At(jCL  ,jT1  ).SetValue(1.0)
  tr.At(jCL  ,jTL  ).SetValue(1.0)
  // repressed 1
  tr.At(jR1  ,jCL  ).SetValue(1.0)
  tr.At(jR1  ,jEA  ).SetValue(1.0)
  tr.At(jR1  ,jEP  ).SetValue(1.0)
  tr.At(jR1  ,jNS  ).SetValue(1.0)
  tr.At(jR1  ,jR2  ).SetValue(1.0)
  tr.At(jR1  ,jPA1 ).SetValue(1.0)
  tr.At(jR1  ,jPB  ).SetValue(1.0)
  tr.At(jR1  ,jT1  ).SetValue(1.0)
  tr.At(jR1  ,jTL  ).SetValue(1.0)
  // repressed 2
  tr.At(jR2  ,jCL  ).SetValue(1.0)
  tr.At(jR2  ,jEA  ).SetValue(1.0)
  tr.At(jR2  ,jEP  ).SetValue(1.0)
  tr.At(jR2  ,jNS  ).SetValue(1.0)
  tr.At(jR2  ,jR1  ).SetValue(1.0)
  tr.At(jR2  ,jPA1 ).SetValue(1.0)
  tr.At(jR2  ,jPB  ).SetValue(1.0)
  tr.At(jR2  ,jT1  ).SetValue(1.0)
  tr.At(jR2  ,jTL  ).SetValue(1.0)
  // promoter active 1
  tr.At(jPA1 ,jT2  ).SetValue(1.0)
  // promoter active 2
  tr.At(jPA2 ,jT2  ).SetValue(1.0)
  tr.At(jPA2 ,jCL  ).SetValue(1.0)
  tr.At(jPA2 ,jNS  ).SetValue(1.0)
  tr.At(jPA2 ,jR1  ).SetValue(1.0)
  tr.At(jPA2 ,jR2  ).SetValue(1.0)
  tr.At(jPA2 ,jTL  ).SetValue(1.0)
  // transcribed 1
  tr.At(jT1  ,jPA2 ).SetValue(1.0)
  tr.At(jT1  ,jEAt1).SetValue(1.0)
  tr.At(jT1  ,jEPt1).SetValue(1.0)
  tr.At(jT1  ,jPBt1).SetValue(1.0)
  // transcribed 2
  tr.At(jT2  ,jPA2 ).SetValue(1.0)
  tr.At(jT2  ,jEAt2).SetValue(1.0)
  tr.At(jT2  ,jEPt2).SetValue(1.0)
  tr.At(jT2  ,jPBt2).SetValue(1.0)
  tr.At(jT2  ,jCL  ).SetValue(1.0)
  tr.At(jT2  ,jNS  ).SetValue(1.0)
  tr.At(jT2  ,jR1  ).SetValue(1.0)
  tr.At(jT2  ,jR2  ).SetValue(1.0)
  tr.At(jT2  ,jTL  ).SetValue(1.0)
  // ea/ep/pb transcribed
  tr.At(jEAt1,jT1  ).SetValue(1.0)
  tr.At(jEPt1,jT1  ).SetValue(1.0)
  tr.At(jPBt1,jT1  ).SetValue(1.0)
  // ea/ep/pb transcribed
  tr.At(jEAt2,jT2  ).SetValue(1.0)
  tr.At(jEPt2,jT2  ).SetValue(1.0)
  tr.At(jPBt2,jT2  ).SetValue(1.0)

  constraints := make([]generic.EqualityConstraint, m)
  switch strings.ToLower(config.Type) {
  case "likelihood":
    printStderr(config, 2, "Implementing constraints for modhmm:likelihood\n")
    // constrain self-transitions
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jPA1, jPA1}, [2]int{jPA2, jPA2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jEA, jEA}, [2]int{jEAt1, jEAt1}, [2]int{jEAt2, jEAt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jEP, jEP}, [2]int{jEPt1, jEPt1}, [2]int{jEPt2, jEPt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jPB, jPB}, [2]int{jPBt1, jPBt1}, [2]int{jPBt2, jPBt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jT1}, [2]int{jT2, jT2}})
    // transition into active enhancers
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jR1, jEA}, [2]int{jR2, jEA}, [2]int{jTL, jEA}, [2]int{jNS, jEA}, [2]int{jCL, jEA}, [2]int{jT1, jEAt1}, [2]int{jT2, jEAt2}})
    // transition into poised enhancers
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jR1, jEP}, [2]int{jR2, jEP}, [2]int{jTL, jEP}, [2]int{jNS, jEP}, [2]int{jCL, jEP}, [2]int{jT1, jEPt1}, [2]int{jT2, jEPt2}})
    // transition into bivalend promoter
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jR1, jPB}, [2]int{jR2, jPB}, [2]int{jTL, jPB}, [2]int{jNS, jPB}, [2]int{jCL, jPB}, [2]int{jT1, jPBt1}, [2]int{jT2, jPBt2}})
  case "posterior":
    printStderr(config, 2, "Implementing constraints for modhmm:posterior\n")
    for i := 0; i < m; i++ {
      constraint := generic.EqualityConstraint{}
      for j := 0; j < m; j++ {
        if i == j {
          continue
        }
        if tr.ConstAt(i, j).GetValue() != 0 {
          constraint = append(constraint, [2]int{i,j})
        }
      }
      constraints = append(constraints, constraint)
    }
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
