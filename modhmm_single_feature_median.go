/* Copyright (C) 2019 Philipp Benner
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

//import "fmt"
//import "bytes"
import "math"
//import "io"
import "sort"

/* -------------------------------------------------------------------------- */

type SlidingMedian struct {
  SortedSlice
}

func NewSlidingMedian() *SlidingMedian {
  return &SlidingMedian{}
}

func (obj *SlidingMedian) Median() float64 {
  if n := len(obj.SortedSlice); n == 0 {
    return math.NaN()
  } else {
    return obj.SortedSlice[n/2]
  }
}

/* -------------------------------------------------------------------------- */

type SortedSlice []float64

func (s_ *SortedSlice) Insert(value float64) {
  s := *s_
  s  = append(s, 0)
  i := sort.SearchFloat64s(s, value)
  if i == len(s) {
    s[i-1] = value
  } else {
    copy(s[i+1:], s[i:])
    s[i] = value
  }
  *s_  = s
}

func (s_ *SortedSlice) Remove(value float64) {
  s := *s_
  i := sort.SearchFloat64s(s, value)
  if i >= len(s) || s[i] != value {
    return
  }
  if i < len(s)-1 {
    copy(s[i:], s[i+1:])
  }
  s[len(s)-1] = math.NaN()
  s = s[:len(s)-1]
  *s_  = s
}
