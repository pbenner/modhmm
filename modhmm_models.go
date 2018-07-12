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

  if estimator, err := matrixEstimator.NewHmmEstimator(pi, tr, nil, nil, nil, estimators, 1e-2, -1); err != nil {
    panic(err)
  } else {
    estimator.ChunkSize = 10000
    estimator.OptimizeEmissions = false
    estimator.Verbose = config.Verbose-1
    return estimator, stateNames
  }
}

/* -------------------------------------------------------------------------- */

func getModHmmDefaultEstimator(config ConfigModHmm) (*matrixEstimator.HmmEstimator, []string) {
  const iPA =  0 // promoter active
  const iPB =  1 // promoter bivalent
  const iEA =  2 // enhancer active
  const iEP =  3 // enhancer poised
  const iTR =  4 // transcribed
  const iTL =  5 // transcribed (low)
  const iR1 =  6 // repressed h3k27me3
  const iR2 =  7 // repressed h3k9me3
  const iCL =  8 // control
  const iNS =  9 // no signal

  const jEA   =  0 // enhancer active
  const jEP   =  1 // enhancer poised
  const jTL   =  2 // transcribed (low)
  const jR1   =  3 // repressed h3k27me3
  const jR2   =  4 // repressed h3k9me3
  const jNS   =  5 // no signal
  const jCL   =  6 // control
  const jPA   =  7 // promoter active
  const jPB   =  8 // promoter bivalent
  const jT1   =  9 // transcribed
  const jT2   = 10 // transcribed
  const jEAt1 = 11 // enhancer active
  const jEPt1 = 12 // enhancer poised
  const jEAt2 = 13 // enhancer active
  const jEPt2 = 14 // enhancer poised

  stateNames := []string{
    "EA", "EP", "TL", "R1", "R2", "NS", "CL", "PA", "PB", "TR", "TR", "EA:tr", "EP:tr", "EA:tr", "EP:tr"}

  n := 10
  m := 15

  stateMap := make([]int, m)
  stateMap[jEA]   = iEA
  stateMap[jEP]   = iEP
  stateMap[jEAt1] = iEA
  stateMap[jEPt1] = iEP
  stateMap[jEAt2] = iEA
  stateMap[jEPt2] = iEP
  stateMap[jPA]   = iPA
  stateMap[jPB]   = iPB
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
  tr.At(jPA  ,jCL  ).SetValue(1.0)
  tr.At(jPA  ,jNS  ).SetValue(1.0)
  tr.At(jPA  ,jR1  ).SetValue(1.0)
  tr.At(jPA  ,jR2  ).SetValue(1.0)
  tr.At(jPA  ,jT2  ).SetValue(1.0)
  tr.At(jPA  ,jTL  ).SetValue(1.0)
  // promoter bivalent
  tr.At(jPB  ,jCL  ).SetValue(1.0)
  tr.At(jPB  ,jNS  ).SetValue(1.0)
  tr.At(jPB  ,jR1  ).SetValue(1.0)
  tr.At(jPB  ,jR2  ).SetValue(1.0)
  tr.At(jPB  ,jT2  ).SetValue(1.0)
  tr.At(jPB  ,jTL  ).SetValue(1.0)
  // no signal
  tr.At(jNS  ,jCL  ).SetValue(1.0)
  tr.At(jNS  ,jEA  ).SetValue(1.0)
  tr.At(jNS  ,jEP  ).SetValue(1.0)
  tr.At(jNS  ,jR1  ).SetValue(1.0)
  tr.At(jNS  ,jR2  ).SetValue(1.0)
  tr.At(jNS  ,jPA  ).SetValue(1.0)
  tr.At(jNS  ,jPB  ).SetValue(1.0)
  tr.At(jNS  ,jT1  ).SetValue(1.0)
  tr.At(jNS  ,jTL  ).SetValue(1.0)
  // control
  tr.At(jCL  ,jEA  ).SetValue(1.0)
  tr.At(jCL  ,jEP  ).SetValue(1.0)
  tr.At(jCL  ,jNS  ).SetValue(1.0)
  tr.At(jCL  ,jR1  ).SetValue(1.0)
  tr.At(jCL  ,jR2  ).SetValue(1.0)
  tr.At(jCL  ,jPA  ).SetValue(1.0)
  tr.At(jCL  ,jPB  ).SetValue(1.0)
  tr.At(jCL  ,jT1  ).SetValue(1.0)
  tr.At(jCL  ,jTL  ).SetValue(1.0)
  // repressed 1
  tr.At(jR1  ,jCL  ).SetValue(1.0)
  tr.At(jR1  ,jEA  ).SetValue(1.0)
  tr.At(jR1  ,jEP  ).SetValue(1.0)
  tr.At(jR1  ,jNS  ).SetValue(1.0)
  tr.At(jR1  ,jR2  ).SetValue(1.0)
  tr.At(jR1  ,jPA  ).SetValue(1.0)
  tr.At(jR1  ,jPB  ).SetValue(1.0)
  tr.At(jR1  ,jT1  ).SetValue(1.0)
  tr.At(jR1  ,jTL  ).SetValue(1.0)
  // repressed 2
  tr.At(jR2  ,jCL  ).SetValue(1.0)
  tr.At(jR2  ,jEA  ).SetValue(1.0)
  tr.At(jR2  ,jEP  ).SetValue(1.0)
  tr.At(jR2  ,jNS  ).SetValue(1.0)
  tr.At(jR2  ,jR1  ).SetValue(1.0)
  tr.At(jR2  ,jPA  ).SetValue(1.0)
  tr.At(jR2  ,jPB  ).SetValue(1.0)
  tr.At(jR2  ,jT1  ).SetValue(1.0)
  tr.At(jR2  ,jTL  ).SetValue(1.0)
  // transcribed 1
  tr.At(jT1  ,jEAt1).SetValue(1.0)
  tr.At(jT1  ,jEPt1).SetValue(1.0)
  tr.At(jT1  ,jPA  ).SetValue(1.0)
  tr.At(jT1  ,jPB  ).SetValue(1.0)
  // transcribed 2
  tr.At(jT2  ,jCL  ).SetValue(1.0)
  tr.At(jT2  ,jEAt2).SetValue(1.0)
  tr.At(jT2  ,jEPt2).SetValue(1.0)
  tr.At(jT2  ,jNS  ).SetValue(1.0)
  tr.At(jT2  ,jR1  ).SetValue(1.0)
  tr.At(jT2  ,jR2  ).SetValue(1.0)
  tr.At(jT2  ,jPA  ).SetValue(1.0)
  tr.At(jT2  ,jPB  ).SetValue(1.0)
  tr.At(jT2  ,jTL  ).SetValue(1.0)
  // transcribed (low)
  tr.At(jTL  ,jCL  ).SetValue(1.0)
  tr.At(jTL  ,jEA  ).SetValue(1.0)
  tr.At(jTL  ,jEP  ).SetValue(1.0)
  tr.At(jTL  ,jNS  ).SetValue(1.0)
  tr.At(jTL  ,jR1  ).SetValue(1.0)
  tr.At(jTL  ,jR2  ).SetValue(1.0)
  tr.At(jTL  ,jPA  ).SetValue(1.0)
  tr.At(jTL  ,jPB  ).SetValue(1.0)
  tr.At(jTL  ,jT1  ).SetValue(1.0)
  // enhancer active transcribed
  tr.At(jEAt1,jT1  ).SetValue(1.0)
  tr.At(jEAt2,jT2  ).SetValue(1.0)
  // enhancer poised transcribed
  tr.At(jEPt1,jT1  ).SetValue(1.0)
  tr.At(jEPt2,jT2  ).SetValue(1.0)

  constraints := make([]generic.EqualityConstraint, m)
  switch strings.ToLower(config.Type) {
  case "": fallthrough
  case "likelihood":
    printStderr(config, 2, "Implementing constraints for modhmm:likelihood\n")
    // constrain self-transitions
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jEA, jEA}, [2]int{jEAt1, jEAt1}, [2]int{jEAt2, jEAt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jEP, jEP}, [2]int{jEPt1, jEPt1}, [2]int{jEPt2, jEPt2}})
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jT1}, [2]int{jT2, jT2}})
    // constrain transitions transcribed -> active enhancers
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jEAt1}, [2]int{jT2, jEAt2}})
    // constrain transitions transcribed -> poised enhancers
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jEPt1}, [2]int{jT2, jEPt2}})
    // constrain transitions transcribed -> active promoters
    constraints = append(constraints, generic.EqualityConstraint{
      [2]int{jT1, jPA}, [2]int{jT2, jPA}})
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

  if estimator, err := matrixEstimator.NewConstrainedHmmEstimator(pi, tr, stateMap, nil, nil, constraints, estimators, 1e-2, -1); err != nil {
    panic(err)
  } else {
    estimator.ChunkSize = 10000
    estimator.OptimizeEmissions = false
    estimator.Verbose = config.Verbose-1
    return estimator, stateNames
  }
}
