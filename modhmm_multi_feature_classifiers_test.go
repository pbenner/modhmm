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

func TestCl1(t *testing.T) {
  m  := BasicClassifier{}
  x  := NewMatrix(BareRealType, 2, 3, []float64{
    math.Log(0.7), math.Log(0.4), math.Log(0.8),
    math.Log(0.3), math.Log(0.6), math.Log(0.2) })

  if math.Abs(math.Exp(m.PeakSym(x, 0, 0)) - 0.736) > 1e-2 {
    t.Error("test failed")
  }
}

func TestCl2(t *testing.T) {
  m  := BasicClassifier{}
  x  := NewMatrix(BareRealType, 2, 4, []float64{
    math.Log(0.7), math.Log(0.4), math.Log(0.8), math.Log(0.3),
    math.Log(0.3), math.Log(0.6), math.Log(0.2), math.Log(0.7) })

  if math.Abs(math.Exp(m.PeakSym(x, 0, 0)) - 0.4628) > 1e-3 {
    t.Error("test failed")
  }
}

func TestCl3(t *testing.T) {
  m  := BasicClassifier{}
  x  := NewMatrix(BareRealType, 2, 5, []float64{
    math.Log(0.7), math.Log(0.4), math.Log(0.8), math.Log(0.3), math.Log(0.1),
    math.Log(0.3), math.Log(0.6), math.Log(0.2), math.Log(0.7), math.Log(0.9) })

  if math.Abs(math.Exp(m.PeakSym(x, 0, 0)) - 0.83632) > 1e-4 {
    t.Error("test failed")
  }
}
