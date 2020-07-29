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
  x  := NewDenseFloat64Matrix([]float64{
    0.7, 0.4, 0.8,
    0.3, 0.6, 0.2 }, 2, 3)

  if math.Abs(m.PeakSym(x, 0, 0) - 0.736) > 1e-2 {
    t.Error("test failed")
  }
}

func TestCl2(t *testing.T) {
  m  := BasicClassifier{}
  x  := NewDenseFloat64Matrix([]float64{
    0.7, 0.4, 0.8, 0.3,
    0.3, 0.6, 0.2, 0.7 }, 2, 4)

  if math.Abs(m.PeakSym(x, 0, 0) - 0.4628) > 1e-3 {
    t.Error("test failed")
  }
}

func TestCl3(t *testing.T) {
  m  := BasicClassifier{}
  x  := NewDenseFloat64Matrix([]float64{
    0.7, 0.4, 0.8, 0.3, 0.1,
    0.3, 0.6, 0.2, 0.7, 0.9 }, 2, 5)

  if math.Abs(m.PeakSym(x, 0, 0) - 0.83632) > 1e-4 {
    t.Error("test failed")
  }
}
