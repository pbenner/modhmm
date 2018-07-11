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
import   "testing"

import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

func Test1(t *testing.T) {
  pi := []float64{math.Log(0.3), math.Log(0.7)}
  m  := BasicMultiFeatureModel{pi}
  x  := NewMatrix(BareRealType, 2, 3, []float64{
    math.Log(0.04), math.Log(0.09), math.Log(0.01),
    math.Log(0.14), math.Log(0.03), math.Log(0.18) })

  if math.Abs(math.Exp(m.PeakSym(x, 0)) - 0.001057537) > 1e-8 {
    t.Error("test failed")
  }
}

func Test2(t *testing.T) {
  pi := []float64{math.Log(0.3), math.Log(0.7)}
  m  := BasicMultiFeatureModel{pi}
  x  := NewMatrix(BareRealType, 2, 4, []float64{
    math.Log(0.04), math.Log(0.09), math.Log(0.01), math.Log(0.04),
    math.Log(0.14), math.Log(0.03), math.Log(0.18), math.Log(0.03) })

  if math.Abs(math.Exp(m.PeakSym(x, 0)) - 6.829634e-06) > 1e-12 {
    t.Error("test failed")
  }
}

func Test3(t *testing.T) {
  pi := []float64{math.Log(0.3), math.Log(0.7)}
  m  := BasicMultiFeatureModel{pi}
  x  := NewMatrix(BareRealType, 2, 5, []float64{
    math.Log(0.04), math.Log(0.09), math.Log(0.01), math.Log(0.04), math.Log(0.25),
    math.Log(0.14), math.Log(0.03), math.Log(0.18), math.Log(0.03), math.Log(0.05) })

  if math.Abs(math.Exp(m.PeakSym(x, 0)) - 1.6519249161373184e-06) > 1e-12 {
    t.Error("test failed")
  }
}
