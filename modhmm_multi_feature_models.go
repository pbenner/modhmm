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

//import   "fmt"
import   "math"


import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import . "github.com/pbenner/autodiff/logarithmetic"

/* -------------------------------------------------------------------------- */

type BasicMultiFeatureModel struct {
  pi []float64
}

func (obj BasicMultiFeatureModel) PeakSym(x ConstMatrix, m int) float64 {
  _, n := x.Dims()
  r    := math.Inf(-1)
  // pattern:
  //
  // 0   1   2   3   4
  //         *           OR
  //     *       *       OR
  // *               *
  //
  // k defines the positive region
  for k := 0; k <= n/2; k++ {
    t := 0.0
    for i := 0; i <= k; i++ {
      j := n-i-1
      if i >= k {
        // positive
        if i == j {
          t += x.ConstAt(m, i).GetValue()
        } else {
          t += x.ConstAt(m, i).GetValue()
          t += x.ConstAt(m, j).GetValue()
        }
      } else {
        // negative
        if i == j {
          t += x.ConstAt(m+1, i).GetValue()
        } else {
          t += x.ConstAt(m+1, i).GetValue()
          t += x.ConstAt(m+1, j).GetValue()
        }
      }
    }
    r = LogAdd(r, t)
  }
  return r
}

func (obj BasicMultiFeatureModel) PeakAny(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := math.Inf(-1)
  for j1 := 0; j1 < n; j1++ {
    t := 0.0
    for j2 := 0; j2 < n; j2++ {
      if j1 == j2 {
        t += x.ConstAt(i+0, j2).GetValue()
      } else
      if j1 >  j2 {
        t += x.ConstAt(i+1, j2).GetValue()
      } else {
        t += LogAdd(
          x.ConstAt(i  , j2).GetValue() + obj.pi[i  ],
          x.ConstAt(i+1, j2).GetValue() + obj.pi[i+1])
      }
    }
    r = LogAdd(r, t)
  }
  return r - LogSub(0, float64(n)*obj.pi[i])
}

func (obj BasicMultiFeatureModel) PeakAt(x ConstMatrix, i, k int) float64 {
     r := 0.0
  _, n := x.Dims()
  for j := 0; j < n; j++ {
    if j == k {
      r += x.ConstAt(i, j).GetValue()
    } else {
      r += LogAdd(
        x.ConstAt(i  , j).GetValue() + obj.pi[i  ],
        x.ConstAt(i+1, j).GetValue() + obj.pi[i+1])
    }
  }
  return r
}

func (obj BasicMultiFeatureModel) PeakAtCenter(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  return obj.PeakAt(x, i, n/2)
}

func (obj BasicMultiFeatureModel) NoPeakAt(x ConstMatrix, i, k int) float64 {
     r := 0.0
  _, n := x.Dims()
  for j := 0; j < n; j++ {
    if j == k {
      r += x.ConstAt(i+1, j).GetValue()
    } else {
      r += LogAdd(
        x.ConstAt(i  , j).GetValue() + obj.pi[i  ],
        x.ConstAt(i+1, j).GetValue() + obj.pi[i+1])
    }
  }
  return r
}

func (obj BasicMultiFeatureModel) NoPeakAtCenter(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  return obj.NoPeakAt(x, i, n/2)
}

func (obj BasicMultiFeatureModel) NoPeakAll(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := 0.0
  for j := 0; j < n; j++ {
    r += x.ConstAt(i+1, j).GetValue()
  }
  return r
}

func (obj BasicMultiFeatureModel) Nil(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := 0.0
  for j := 0; j < n; j++ {
    r += LogAdd(
      x.ConstAt(i  , j).GetValue() + obj.pi[i  ],
      x.ConstAt(i+1, j).GetValue() + obj.pi[i+1])
  }
  return r
}

/* -------------------------------------------------------------------------- */

type ModelPA struct {
  BasicMultiFeatureModel
}

func (obj ModelPA) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // atac peak at the center
    r += obj.PeakAtCenter(x, jAtac)
  }
  { // h3k27ac peak at any position
    r += obj.PeakAny(x, jH3k27ac)
  }
  { // h3k4me3 peak at any position
    r += obj.PeakAny(x, jH3k4me3)
  }
  { // h3k4me3o1 peak at any position
    r += obj.PeakAny(x, jH3k4me3o1)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelPA) Dims() (int, int) {
  return 20, 5
}

func (obj ModelPA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelPA{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelPB struct {
  BasicMultiFeatureModel
}

func (obj ModelPB) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // atac peak at the center
    //r += obj.PeakAtCenter(x, jAtac)
  }
  { // h3k27me3 peak at any position
    r += obj.PeakAny(x, jH3k27me3)
  }
  { // h3k4me3 peak at any position
    r += obj.PeakAny(x, jH3k4me3)
  }
  { // h3k4me3o1 peak at any position
    r += obj.PeakAny(x, jH3k4me3o1)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelPB) Dims() (int, int) {
  return 20, 5
}

func (obj ModelPB) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelPB{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelEA struct {
  BasicMultiFeatureModel
}

func (obj ModelEA) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  // atac peak at the center
  r += obj.PeakAtCenter(x, jAtac)
  // h3k27ac peak at any position
  r += obj.PeakAny(x, jH3k27ac)
  // h3k27me3 nil
  r += obj.Nil(x, jH3k27me3)
  // h3k9me3 nil
  r += obj.Nil(x, jH3k9me3)
  // h3k4me1 peak at any position
  r += obj.PeakSym(x, jH3k4me1)
  // h3k4me3 nil
  r += obj.Nil(x, jH3k4me3)
  // no h3k4me3o1 peak at all positions
  r += obj.NoPeakAll(x, jH3k4me3o1)
  // rna nil
  r += obj.Nil(x, jRna)
  // rna-low nil
  r += obj.Nil(x, jRnaLow)
  // no control peak at all positions
  r += obj.NoPeakAll(x, jControl)

  s.SetValue(r); return nil
}

func (ModelEA) Dims() (int, int) {
  return 20, 5
}

func (obj ModelEA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelEA{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelEP struct {
  BasicMultiFeatureModel
}

func (obj ModelEP) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  // atac nil
  r += obj.Nil(x, jAtac)
  // h3k27ac nil
  r += obj.Nil(x, jH3k27ac)
  // h3k27me3 peak at any position
  r += obj.PeakAny(x, jH3k27me3)
  // h3k9me3 peak at any position
  r += obj.Nil(x, jH3k9me3)
  // h3k4me1 peak at any position
  r += obj.PeakSym(x, jH3k4me1)
  // h3k4me3 peak at any position
  r += obj.Nil(x, jH3k4me3)
  // no h3k4me3o1 peak at all positions
  r += obj.NoPeakAll(x, jH3k4me3o1)
  // rna nil
  r += obj.Nil(x, jRna)
  // rna-low nil
  r += obj.Nil(x, jRnaLow)
  // no control peak at all positions
  r += obj.NoPeakAll(x, jControl)

  s.SetValue(r); return nil
}

func (ModelEP) Dims() (int, int) {
  return 20, 5
}

func (obj ModelEP) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelEP{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelTR struct {
  BasicMultiFeatureModel
}

func (obj ModelTR) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // no atac peak at center
    r += obj.NoPeakAll(x, jAtac)
  }
  { // no h3k4me1 peak at center
    r += obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at center
    r += obj.NoPeakAll(x, jH3k4me3)
  }
  { // rna peak at center
    r += obj.PeakAtCenter(x, jRna)
  }
  s.SetValue(r)
  return nil
}

func (ModelTR) Dims() (int, int) {
  return 20, 1
}

func (obj ModelTR) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelTR{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelTL struct {
  BasicMultiFeatureModel
}

func (obj ModelTL) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // no atac peak at center
    r += obj.NoPeakAll(x, jAtac)
  }
  { // no h3k4me1 peak at center
    r += obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at center
    r += obj.NoPeakAll(x, jH3k4me3)
  }
  { // rna peak at center
    r += obj.PeakAtCenter(x, jRnaLow)
  }
  s.SetValue(r)
  return nil
}

func (ModelTL) Dims() (int, int) {
  return 20, 1
}

func (obj ModelTL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelTL{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelR1 struct {
  BasicMultiFeatureModel
}

func (obj ModelR1) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // h3k27me3 peak at any position
    r += obj.PeakAny(x, jH3k27me3)
  }
  { // no h3k4me1 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me3)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelR1) Dims() (int, int) {
  return 20, 1
}

func (obj ModelR1) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelR1{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelR2 struct {
  BasicMultiFeatureModel
}

func (obj ModelR2) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // h3k9me3 peak at any position
    r += obj.PeakAny(x, jH3k9me3)
  }
  { // no h3k4me1 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me3)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelR2) Dims() (int, int) {
  return 20, 1
}

func (obj ModelR2) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelR2{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelNS struct {
  BasicMultiFeatureModel
}

func (obj ModelNS) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // no atac peak at any position
    r += obj.NoPeakAll(x, jAtac)
  }
  { // no h3k27ac peak at any position
    r += obj.NoPeakAll(x, jH3k27ac)
  }
  { // no h3k27me3 peak at any position
    r += obj.NoPeakAll(x, jH3k27me3)
  }
  { // no h3k9me3 peak at any position
    r += obj.NoPeakAll(x, jH3k9me3)
  }
  { // no h3k4me1 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me3)
  }
  { // no rna peak at all positions
    r += obj.NoPeakAll(x, jRna)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelNS) Dims() (int, int) {
  return 20, 1
}

func (obj ModelNS) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelNS{BasicMultiFeatureModel{obj.pi}}
}

/* -------------------------------------------------------------------------- */

type ModelCL struct {
  BasicMultiFeatureModel
}

func (obj ModelCL) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // control peak at any position
    r += obj.PeakAny(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ModelCL) Dims() (int, int) {
  return 20, 1
}

func (obj ModelCL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelCL{BasicMultiFeatureModel{obj.pi}}
}
