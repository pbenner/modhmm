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

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"

import . "github.com/pbenner/modhmm/utility"

/* -------------------------------------------------------------------------- */

type BasicClassifier struct {
}

func (obj BasicClassifier) PeakSym_(x ConstMatrix, m, min, k0 int) float64 {
  _, n := x.Dims()
  r    := 0.0
  // pattern:
  //
  // 0   1   2   3   4
  //         *           OR
  //     *       *       OR
  // *               *
  //
  // k defines the positive region
  for k := k0; k < DivIntUp(n,2); k++ {
    t := 1.0
    for i := k0; i <= k; i++ {
      j := n-i-1
      if j - i + 1 < min {
        break
      }
      xi := x.ValueAt(m, i)
      xj := x.ValueAt(m, j)
      if i >= k {
        // positive
        if i == j {
          t *= xi
        } else {
          t *= xi
          t *= xj
        }
      } else {
        // negative
        if i == j {
          t *= 1.0 - xi
        } else {
          t *= xi*(1.0-xj) + (1.0-xi)*xj + (1.0-xi)*(1.0-xj)
        }
      }
    }
    r += t
  }
  return r
}

func (obj BasicClassifier) PeakSym(x ConstMatrix, m, min int) float64 {
  return obj.PeakSym_(x, m, min, 0)
}

func (obj BasicClassifier) PeakAny(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  t    :=     x.ValueAt(i, 0)
  r    := 1.0-x.ValueAt(i, 0)
  for k := 1; k < n; k++ {
    t +=   r*x.ValueAt(i, k)
    r *= 1.0-x.ValueAt(i, k)
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
  r    := 1.0
  for k := 0; k < n; k++ {
    r *= x.ValueAt(i, k)
  }
  return r
}

func (obj BasicClassifier) PeakRange(x ConstMatrix, i, k1, k2 int) float64 {
  r := 1.0
  for j := k1; j < k2; j++ {
    r *= x.ValueAt(i, j)
  }
  return r
}

func (obj BasicClassifier) NoPeakRange(x ConstMatrix, i, k1, k2 int) float64 {
  r := 1.0
  for j := k1; j < k2; j++ {
    r *= 1.0-x.ValueAt(i, j)
  }
  return r
}

func (obj BasicClassifier) NoPeakAt(x ConstMatrix, i, k int) float64 {
  return 1.0-x.ValueAt(i, k)
}

func (obj BasicClassifier) NoPeakAtCenter(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  return 1.0-x.ValueAt(i, n/2)
}

func (obj BasicClassifier) NoPeakAll(x ConstMatrix, i int) float64 {
  _, n := x.Dims()
  r    := 1.0
  for k := 0; k < n; k++ {
    r *= 1.0-x.ValueAt(i, k)
  }
  return r
}

/* -------------------------------------------------------------------------- */

type ClassifierPA struct {
  BasicClassifier
}

func (obj ClassifierPA) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // atac peak at the center
    r *= obj.PeakAtCenter(x, jOpen)
  }
  { // h3k27ac peak at any position
    r *= obj.PeakSym(x, jH3k27ac, 0)
  }
  { // h3k4me3 peak at any position
    r *= obj.PeakAny(x, jH3k4me3)
  }
  { // h3k4me3o1 peak at any position
    r *= obj.PeakAny(x, jH3k4me3o1)
  }
  { // no control peak at all positions
    r *= obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierPA) Dims() (int, int) {
  return 10, 5
}

func (ClassifierPA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierPA{}
}

/* -------------------------------------------------------------------------- */

type ClassifierEA struct {
  BasicClassifier
}

func (obj ClassifierEA) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // atac peak at the center
    r *= obj.PeakAtCenter(x, jOpen)
  }
  { // h3k27ac peak at any position
    r *= obj.PeakSym_(x, jH3k27ac, 0, 1)
  }
  { // h3k4me1 peak at any position
    r *= obj.PeakSym(x, jH3k4me1, 0)
  }
  { // no h3k4me3o1 peak at all positions
    r *= obj.NoPeakRange(x, jH3k4me3o1, 1, 6)
  }
  { // no control peak at all positions
    r *= obj.NoPeakRange(x, jControl, 1, 6)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierEA) Dims() (int, int) {
  return 10, 7
}

func (ClassifierEA) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierEA{}
}

/* -------------------------------------------------------------------------- */

type ClassifierBI struct {
  BasicClassifier
}

func (obj ClassifierBI) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // atac peak at the center
    //r *= obj.PeakAtCenter(x, jOpen)
  }
  { // h3k27me3 peak at any position
    r *= obj.PeakSym_(x, jH3k27me3, 0, 1)
  }
  { // symmetric jH3k4me1 peak or h3k4me3 peak at any position
    t1 := obj.PeakSym  (x, jH3k4me1, 0)
    t2 := obj.PeakRange(x, jH3k4me3, 1, 6)
    r  *= t1 + (1.0-t1)*t2
  }
  { // no control peak at all positions
    r *= obj.NoPeakRange(x, jControl, 1, 6)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierBI) Dims() (int, int) {
  return 10, 7
}

func (ClassifierBI) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierBI{}
}

/* -------------------------------------------------------------------------- */

type ClassifierPR struct {
  BasicClassifier
}

func (obj ClassifierPR) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // atac peak at the center
    r *= obj.PeakAtCenter(x, jOpen)
  }
  { // no h3k27ac peak
    r *= obj.NoPeakRange(x, jH3k27ac, 1, 6)
  }
  { // no h3k27me3 peak
    r *= obj.NoPeakRange(x, jH3k27me3, 1, 6)
  }
  { // symmetric jH3k4me1 peak or h3k4me3 peak at any position
    t1 := obj.PeakSym  (x, jH3k4me1, 0)
    t2 := obj.PeakRange(x, jH3k4me3, 1, 6)
    r  *= t1 + (1.0-t1)*t2
  }
  { // no control peak at all positions
    r *= obj.NoPeakRange(x, jControl, 1, 6)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierPR) Dims() (int, int) {
  return 10, 7
}

func (ClassifierPR) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierPR{}
}

/* -------------------------------------------------------------------------- */

type ClassifierTR struct {
  BasicClassifier
}

func (obj ClassifierTR) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // no atac peak at center
    r *= obj.NoPeakAll(x, jOpen)
  }
  { // no h3k4me1 peak at center
    //r *= obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at center
    r *= obj.NoPeakAll(x, jH3k4me3)
  }
  { // rna peak at center
    r *= obj.PeakAtCenter(x, jRna)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierTR) Dims() (int, int) {
  return 10, 1
}

func (ClassifierTR) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierTR{}
}

/* -------------------------------------------------------------------------- */

type ClassifierTL struct {
  BasicClassifier
}

func (obj ClassifierTL) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // no atac peak at center
    r *= obj.NoPeakAll(x, jOpen)
  }
  { // no h3k4me1 peak at center
    //r *= obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at center
    r *= obj.NoPeakAll(x, jH3k4me3)
  }
  { // rna peak at center
    r *= obj.PeakAtCenter(x, jRnaLow)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierTL) Dims() (int, int) {
  return 10, 1
}

func (ClassifierTL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierTL{}
}

/* -------------------------------------------------------------------------- */

type ClassifierR1 struct {
  BasicClassifier
}

func (obj ClassifierR1) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // h3k27me3 peak at any position
    r *= obj.PeakAny(x, jH3k27me3)
  }
  { // no h3k4me1 peak at all positions
    //r *= obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r *= obj.NoPeakAll(x, jH3k4me3)
  }
  { // no control peak at all positions
    r *= obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierR1) Dims() (int, int) {
  return 10, 1
}

func (ClassifierR1) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierR1{}
}

/* -------------------------------------------------------------------------- */

type ClassifierR2 struct {
  BasicClassifier
}

func (obj ClassifierR2) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // h3k9me3 peak at any position
    r *= obj.PeakAny(x, jH3k9me3)
  }
  { // no h3k4me1 peak at all positions
    //r *= obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r *= obj.NoPeakAll(x, jH3k4me3)
  }
  { // no control peak at all positions
    r *= obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierR2) Dims() (int, int) {
  return 10, 1
}

func (ClassifierR2) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierR2{}
}

/* -------------------------------------------------------------------------- */

type ClassifierNS struct {
  BasicClassifier
}

func (obj ClassifierNS) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // no atac peak at any position
    r *= obj.NoPeakAll(x, jOpen)
  }
  { // no h3k27ac peak at any position
    r *= obj.NoPeakAll(x, jH3k27ac)
  }
  { // no h3k27me3 peak at any position
    r *= obj.NoPeakAll(x, jH3k27me3)
  }
  { // no h3k9me3 peak at any position
    r *= obj.NoPeakAll(x, jH3k9me3)
  }
  { // no h3k4me1 peak at all positions
    //r *= obj.NoPeakAll(x, jH3k4me1)
  }
  { // no h3k4me3 peak at all positions
    r *= obj.NoPeakAll(x, jH3k4me3)
  }
  { // no rna peak at all positions
    r *= obj.NoPeakAll(x, jRna)
  }
  { // no rna-low peak at all positions
    r *= obj.NoPeakAll(x, jRnaLow)
  }
  { // no control peak at all positions
    r *= obj.NoPeakAll(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierNS) Dims() (int, int) {
  return 10, 1
}

func (ClassifierNS) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierNS{}
}

/* -------------------------------------------------------------------------- */

type ClassifierCL struct {
  BasicClassifier
}

func (obj ClassifierCL) Eval(s Scalar, x ConstMatrix) error {
  r := 1.0
  { // control peak at any position
    r *= obj.PeakAny(x, jControl)
  }
  s.SetValue(r)
  return nil
}

func (ClassifierCL) Dims() (int, int) {
  return 10, 1
}

func (ClassifierCL) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return ClassifierCL{}
}
