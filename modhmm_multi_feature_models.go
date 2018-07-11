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

func (obj BasicMultiFeatureModel) PeakSym(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := math.Inf(-1)
  s    := math.Inf(-1)
  // pattern:
  //
  // 0   1   2   3   4
  //         *           OR
  //     *       *       OR
  // *               *
  //
  c := 0.0
  d := 0.0
  // j1 defines the positive region
  for j1 := 0; j1 < divIntUp(n,2); j1++ {
    j2 := n-j1-1
    if j1 == j2 {
      r = LogAdd(r, c + x.ValueAt(i, j1) + obj.pi[i])
      s = LogAdd(s, d + obj.pi[i])
    } else {
      // peak at (j1, j2)
      t1 := x.ValueAt(i, j1) + obj.pi[i] + x.ValueAt(i, j2) + obj.pi[i]
      t2 := 2.0*obj.pi[i]
      // no peak at (j1, j2)
      s1 := LogAdd(x.ValueAt(i+0, j1) + obj.pi[i+0] + x.ValueAt(i+1, j2) + obj.pi[i+1],  // p(x_j1,    peak at j1) p(x_j2, no peak at j2)
            LogAdd(x.ValueAt(i+1, j1) + obj.pi[i+1] + x.ValueAt(i+0, j2) + obj.pi[i+0],  // p(x_j1, no peak at j1) p(x_j2,    peak at j2)
                   x.ValueAt(i+1, j1) + obj.pi[i+1] + x.ValueAt(i+1, j2) + obj.pi[i+1])) // p(x_j1, no peak at j1) p(x_j2, no peak at j2)
      s2 := LogAdd(obj.pi[i+0] + obj.pi[i+1],  // p(   peak at j1) p(no peak at j2)
            LogAdd(obj.pi[i+1] + obj.pi[i+0],  // p(no peak at j1) p(   peak at j2)
                   obj.pi[i+1] + obj.pi[i+1])) // p(no peak at j1) p(no peak at j2)
      // peak at (j1, j2) => t1 * p(x_{j1+1}) ... p(x_{j2-1})
      for k := j1+1; k < j2; k++ {
        t1 += LogAdd(x.ValueAt(i+0, k) + obj.pi[i+0],
                     x.ValueAt(i+1, k) + obj.pi[i+1])
      }
      r = LogAdd(r, c+t1)
      s = LogAdd(s, d+t2)
      // no peak at (j1, j2)
      c += s1
      d += s2
    }
  }
  return r - s
}

func (obj BasicMultiFeatureModel) PeakAny(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := math.Inf(-1)
  for j1 := 0; j1 < n; j1++ {
    t := 0.0
    for j2 := 0; j2 < n; j2++ {
      if j1 == j2 {
        t += x.ValueAt(i+0, j2) + obj.pi[i  ]
      } else
      if j1 >  j2 {
        t += x.ValueAt(i+1, j2) + obj.pi[i+1]
      } else {
        t += LogAdd(
          x.ValueAt(i  , j2) + obj.pi[i  ],
          x.ValueAt(i+1, j2) + obj.pi[i+1])
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
      r += x.ValueAt(i, j)
    } else {
      r += LogAdd(
        x.ValueAt(i  , j) + obj.pi[i  ],
        x.ValueAt(i+1, j) + obj.pi[i+1])
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
      r += x.ValueAt(i+1, j)
    } else {
      r += LogAdd(
        x.ValueAt(i  , j) + obj.pi[i  ],
        x.ValueAt(i+1, j) + obj.pi[i+1])
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
    r += x.ValueAt(i+1, j)
  }
  return r
}

func (obj BasicMultiFeatureModel) Nil(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := 0.0
  for j := 0; j < n; j++ {
    r += LogAdd(
      x.ValueAt(i  , j) + obj.pi[i  ],
      x.ValueAt(i+1, j) + obj.pi[i+1])
  }
  return r
}

/* -------------------------------------------------------------------------- */

type ModelPA struct {
  BasicMultiFeatureModel
}

func (obj ModelPA) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  r += obj.PeakAtCenter(x, jAtac)
  r += obj.PeakAny     (x, jH3k27ac)
  r += obj.Nil         (x, jH3k27me3)
  r += obj.Nil         (x, jH3k9me3)
  r += obj.Nil         (x, jH3k4me1)
  r += obj.PeakAny     (x, jH3k4me3)
  r += obj.PeakAny     (x, jH3k4me3o1)
  r += obj.Nil         (x, jRna)
  r += obj.Nil         (x, jRnaLow)
  r += obj.NoPeakAll   (x, jControl)

  s.SetValue(r); return nil
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
  r += obj.Nil      (x, jAtac)
  r += obj.Nil      (x, jH3k27ac)
  r += obj.PeakAny  (x, jH3k27me3)
  r += obj.Nil      (x, jH3k9me3)
  r += obj.Nil      (x, jH3k4me1)
  r += obj.PeakAny  (x, jH3k4me3)
  r += obj.PeakAny  (x, jH3k4me3o1)
  r += obj.Nil      (x, jRna)
  r += obj.Nil      (x, jRnaLow)
  r += obj.NoPeakAll(x, jControl)

  s.SetValue(r); return nil
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
  r += obj.PeakAtCenter(x, jAtac)
  r += obj.PeakAny     (x, jH3k27ac)
  r += obj.Nil         (x, jH3k27me3)
  r += obj.Nil         (x, jH3k9me3)
  r += obj.PeakSym     (x, jH3k4me1)
  r += obj.Nil         (x, jH3k4me3)
  r += obj.NoPeakAll   (x, jH3k4me3o1)
  r += obj.Nil         (x, jRna)
  r += obj.Nil         (x, jRnaLow)
  r += obj.NoPeakAll   (x, jControl)

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
  r += obj.Nil      (x, jAtac)
  r += obj.Nil      (x, jH3k27ac)
  r += obj.PeakAny  (x, jH3k27me3)
  r += obj.Nil      (x, jH3k9me3)
  r += obj.PeakSym  (x, jH3k4me1)
  r += obj.Nil      (x, jH3k4me3)
  r += obj.NoPeakAll(x, jH3k4me3o1)
  r += obj.Nil      (x, jRna)
  r += obj.Nil      (x, jRnaLow)
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
  r += obj.NoPeakAtCenter(x, jAtac)
  r += obj.Nil           (x, jH3k27ac)
  r += obj.Nil           (x, jH3k27me3)
  r += obj.Nil           (x, jH3k9me3)
  r += obj.NoPeakAtCenter(x, jH3k4me1)
  r += obj.NoPeakAtCenter(x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.PeakAtCenter  (x, jRna)
  r += obj.Nil           (x, jRnaLow)
  r += obj.Nil           (x, jControl)

  s.SetValue(r); return nil
}

func (ModelTR) Dims() (int, int) {
  return 20, 5
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
  r += obj.NoPeakAtCenter(x, jAtac)
  r += obj.Nil           (x, jH3k27ac)
  r += obj.Nil           (x, jH3k27me3)
  r += obj.Nil           (x, jH3k9me3)
  r += obj.NoPeakAtCenter(x, jH3k4me1)
  r += obj.NoPeakAtCenter(x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.Nil           (x, jRna)
  r += obj.PeakAtCenter  (x, jRnaLow)
  r += obj.Nil           (x, jControl)

  s.SetValue(r); return nil
}

func (ModelTL) Dims() (int, int) {
  return 20, 5
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
  r += obj.Nil           (x, jAtac)
  r += obj.Nil           (x, jH3k27ac)
  r += obj.PeakAtCenter  (x, jH3k27me3)
  r += obj.Nil           (x, jH3k9me3)
  r += obj.NoPeakAtCenter(x, jH3k4me1)
  r += obj.NoPeakAtCenter(x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.Nil           (x, jRna)
  r += obj.Nil           (x, jRnaLow)
  r += obj.NoPeakAtCenter(x, jControl)
  s.SetValue(r); return nil
}

func (ModelR1) Dims() (int, int) {
  return 20, 5
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
  r += obj.Nil           (x, jAtac)
  r += obj.Nil           (x, jH3k27ac)
  r += obj.Nil           (x, jH3k27me3)
  r += obj.PeakAtCenter  (x, jH3k9me3)
  r += obj.NoPeakAtCenter(x, jH3k4me1)
  r += obj.NoPeakAtCenter(x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.Nil           (x, jRna)
  r += obj.Nil           (x, jRnaLow)
  r += obj.NoPeakAtCenter(x, jControl)

  s.SetValue(r); return nil
}

func (ModelR2) Dims() (int, int) {
  return 20, 5
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
  r += obj.NoPeakAtCenter(x, jAtac)
  r += obj.NoPeakAtCenter(x, jH3k27ac)
  r += obj.NoPeakAtCenter(x, jH3k27me3)
  r += obj.NoPeakAtCenter(x, jH3k9me3)
  r += obj.NoPeakAtCenter(x, jH3k4me1)
  r += obj.NoPeakAtCenter(x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.NoPeakAtCenter(x, jRna)
  r += obj.Nil           (x, jRnaLow)
  r += obj.NoPeakAtCenter(x, jControl)

  s.SetValue(r); return nil
}

func (ModelNS) Dims() (int, int) {
  return 20, 5
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
  r += obj.Nil           (x, jAtac)
  r += obj.Nil           (x, jH3k27ac)
  r += obj.Nil           (x, jH3k27me3)
  r += obj.Nil           (x, jH3k9me3)
  r += obj.Nil           (x, jH3k4me1)
  r += obj.Nil           (x, jH3k4me3)
  r += obj.Nil           (x, jH3k4me3o1)
  r += obj.Nil           (x, jRna)
  r += obj.Nil           (x, jRnaLow)
  r += obj.NoPeakAtCenter(x, jControl)

  s.SetValue(r); return nil
}

func (ModelCL) Dims() (int, int) {
  return 20, 5
}

func (obj ModelCL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ModelCL{BasicMultiFeatureModel{obj.pi}}
}