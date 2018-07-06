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

const jAtac      =  0
const jH3k27ac   =  2
const jH3k27me3  =  4
const jH3k9me1   =  6
const jH3k4me1   =  8
const jH3k4me3   = 10
const jH3k4me3o1 = 12
const jRna       = 14
const jRnaLow    = 16
const jControl   = 18

/* -------------------------------------------------------------------------- */

type BasicClassifier struct {
}

func (obj BasicClassifier) PeakSym(x ConstMatrix, m int) float64 {
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

func (obj BasicClassifier) PeakAny(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  t    := x.ConstAt(i,   0).GetValue()
  r    := x.ConstAt(i+1, 0).GetValue()
  for k := 1; k < n; k++ {
    t  = LogAdd(t, r + x.ConstAt(i, k).GetValue())
    r += x.ConstAt(i+1, k).GetValue()
  }
  return t
}

func (obj BasicClassifier) PeakAt(x ConstMatrix, i, k int) float64 {
  return x.ConstAt(i, k).GetValue()
}

func (obj BasicClassifier) PeakAtCenter(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  return x.ConstAt(i, n/2).GetValue()
}

func (obj BasicClassifier) NoPeakAt(x ConstMatrix, i, k int) float64 {
  return x.ConstAt(i+1, k).GetValue()
}

func (obj BasicClassifier) NoPeakAll(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := x.ConstAt(i+1, 0).GetValue()
  for k := 1; k < n; k++ {
    r += x.ConstAt(i+1, k).GetValue()
  }
  return r
}

/* -------------------------------------------------------------------------- */

type ClassifierPA struct {
  BasicClassifier
}

func (obj ClassifierPA) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierPA) Dims() (int, int) {
  return 20, 5
}

func (ClassifierPA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierPA{}
}

/* -------------------------------------------------------------------------- */

type ClassifierPB struct {
  BasicClassifier
}

func (obj ClassifierPB) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierPB) Dims() (int, int) {
  return 20, 5
}

func (ClassifierPB) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierPB{}
}

/* -------------------------------------------------------------------------- */

type ClassifierEA struct {
  BasicClassifier
}

func (obj ClassifierEA) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // atac peak at the center
    r += obj.PeakAtCenter(x, jAtac)
  }
  { // h3k27ac peak at any position
    r += obj.PeakAny(x, jH3k27ac)
  }
  { // h3k4me1 peak at any position
    r += obj.PeakSym(x, jH3k4me1)
  }
  { // no h3k4me3o1 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me3o1)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierEA) Dims() (int, int) {
  return 20, 5
}

func (ClassifierEA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierEA{}
}

/* -------------------------------------------------------------------------- */

type ClassifierEP struct {
  BasicClassifier
}

func (obj ClassifierEP) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // atac peak at the center
    //r += obj.PeakAtCenter(x, jAtac)
  }
  { // h3k27me3 peak at any position
    r += obj.PeakAny(x, jH3k27me3)
  }
  { // h3k4me1 peak at any position
    r += obj.PeakSym(x, jH3k4me1)
  }
  { // no h3k4me3o1 peak at all positions
    r += obj.NoPeakAll(x, jH3k4me3o1)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierEP) Dims() (int, int) {
  return 20, 5
}

func (ClassifierEP) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierEP{}
}

/* -------------------------------------------------------------------------- */

type ClassifierTR struct {
  BasicClassifier
}

func (obj ClassifierTR) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierTR) Dims() (int, int) {
  return 20, 1
}

func (ClassifierTR) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierTR{}
}

/* -------------------------------------------------------------------------- */

type ClassifierTL struct {
  BasicClassifier
}

func (obj ClassifierTL) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierTL) Dims() (int, int) {
  return 20, 1
}

func (ClassifierTL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierTL{}
}
