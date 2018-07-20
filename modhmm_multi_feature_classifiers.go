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

type BasicClassifier struct {
}

func (obj BasicClassifier) PeakSym(x ConstMatrix, m, min int) float64 {
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
  for k := 0; k < divIntUp(n,2); k++ {
    t := 0.0
    for i := 0; i <= k; i++ {
      j := n-i-1
      if j - i + 1 < min {
        break
      }
      if i >= k {
        // positive
        if i == j {
          t += x.ValueAt(m, i)
        } else {
          t += x.ValueAt(m, i)
          t += x.ValueAt(m, j)
        }
      } else {
        // negative
        if i == j {
          t += x.ValueAt(m+1, i)
        } else {
          t += LogAdd(LogAdd(
            x.ValueAt(m+0, i)+x.ValueAt(m+1, j),
            x.ValueAt(m+1, i)+x.ValueAt(m+0, j)),
            x.ValueAt(m+1, i)+x.ValueAt(m+1, j))
        }
      }
    }
    r = LogAdd(r, t)
  }
  return r
}

func (obj BasicClassifier) PeakAny(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  t    := x.ValueAt(i,   0)
  r    := x.ValueAt(i+1, 0)
  for k := 1; k < n; k++ {
    t  = LogAdd(t, r + x.ValueAt(i, k))
    r += x.ValueAt(i+1, k)
  }
  return t
}

func (obj BasicClassifier) PeakAt(x ConstMatrix, i, k int) float64 {
  return x.ValueAt(i, k)
}

func (obj BasicClassifier) PeakAtCenter(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  return x.ValueAt(i, n/2)
}

func (obj BasicClassifier) PeakAll(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := x.ValueAt(i, 0)
  for k := 1; k < n; k++ {
    r += x.ValueAt(i, k)
  }
  return r
}

func (obj BasicClassifier) PeakRange(x ConstMatrix, i, k1, k2 int) float64 {
  r := 0.0
  for j := k1; j < k2; j++ {
    r += x.ValueAt(i, j)
  }
  return r
}

func (obj BasicClassifier) NoPeakRange(x ConstMatrix, i, k1, k2 int) float64 {
  r := 0.0
  for j := k1; j < k2; j++ {
    r += x.ValueAt(i+1, j)
  }
  return r
}

func (obj BasicClassifier) NoPeakAt(x ConstMatrix, i, k int) float64 {
  return x.ValueAt(i+1, k)
}

func (obj BasicClassifier) NoPeakAll(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := x.ValueAt(i+1, 0)
  for k := 1; k < n; k++ {
    r += x.ValueAt(i+1, k)
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
    r += obj.PeakSym(x, jH3k27ac, 0)
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
    r += obj.PeakSym(x, jH3k27me3, 0)
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
    r += obj.PeakSym(x, jH3k27ac, 0)
  }
  { // h3k4me1 peak at any position
    r += obj.PeakSym(x, jH3k4me1, 2)
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
    r += obj.PeakSym(x, jH3k27me3, 0)
  }
  { // h3k4me1 peak at any position
    r += obj.PeakSym(x, jH3k4me1, 2)
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

/* -------------------------------------------------------------------------- */

type ClassifierR1 struct {
  BasicClassifier
}

func (obj ClassifierR1) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierR1) Dims() (int, int) {
  return 20, 1
}

func (ClassifierR1) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierR1{}
}

/* -------------------------------------------------------------------------- */

type ClassifierR2 struct {
  BasicClassifier
}

func (obj ClassifierR2) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierR2) Dims() (int, int) {
  return 20, 1
}

func (ClassifierR2) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierR2{}
}

/* -------------------------------------------------------------------------- */

type ClassifierNS struct {
  BasicClassifier
}

func (obj ClassifierNS) Eval(s Scalar, x ConstMatrix) error {
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

func (ClassifierNS) Dims() (int, int) {
  return 20, 1
}

func (ClassifierNS) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierNS{}
}

/* -------------------------------------------------------------------------- */

type ClassifierCL struct {
  BasicClassifier
}

func (obj ClassifierCL) Eval(s Scalar, x ConstMatrix) error {
  r := 0.0
  { // control peak at any position
    r += obj.PeakAny(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierCL) Dims() (int, int) {
  return 20, 1
}

func (ClassifierCL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierCL{}
}
