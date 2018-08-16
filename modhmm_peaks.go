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

import . "github.com/pbenner/gonetics"

/* -------------------------------------------------------------------------- */

func positive(sequence TrackSequence, threshold float64, i int) bool {
  if math.IsNaN(sequence.AtBin(i)) || sequence.AtBin(i) < threshold {
    return false
  }
  return true
}

func getPeaks(track Track, threshold float64) (GRanges, error) {
  seqnames := []string{}
  from     := []int{}
  to       := []int{}
  strand   := []byte{}
  test     := []float64{}

  for _, name := range track.GetSeqNames() {
    s, err := track.GetSequence(name); if err != nil {
      return GRanges{}, err
    }
    binsize   := s.GetBinSize()
    seqlen    := s.NBins()
    for i := 0; i < seqlen; i++ {
      if positive(s, threshold, i) {
        // peak begins here
        i_from := i
        // maximum value
        v_max  := s.AtBin(i)
        // increment until either the sequence ended or
        // the value drops below the threshold
        for i < seqlen && positive(s, threshold, i) {
          if v := s.AtBin(i); v > v_max {
            // update maximum position and value
            v_max = v
          }
          i += 1
        }
        // save peak
        seqnames = append(seqnames, name)
        from     = append(from,     i_from*binsize)
        to       = append(to,       i     *binsize)
        test     = append(test,     math.Exp(v_max))
      }
    }
  }
  peaks := NewGRanges(seqnames, from, to, strand)
  peaks.AddMeta("test", test)
  peaks, _ = peaks.Sort("test", true)

  return peaks, nil
}
