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
    r += obj.PeakAtCenter(x, 0)
  }
  { // h3k27ac peak at any position
    r += obj.PeakAny(x, 2)
  }
  { // h3k4me3 peak at any position
    r += obj.PeakAny(x, 4)
  }
  { // h3k4me3o1 peak at any position
    r += obj.PeakAny(x, 6)
  }
  { // no control peak at all positions
    r += obj.NoPeakAll(x, 8)
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
