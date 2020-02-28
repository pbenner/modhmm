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

//import   "fmt"
import   "log"
import   "math"

import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/ngstat/track"
import   "github.com/pbenner/ngstat/utility"
import   "github.com/pbenner/threadpool"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

/* -------------------------------------------------------------------------- */

func single_feature_eval_heuristic_loop(config ConfigModHmm, files SingleFeatureFiles) MutableTrack {
  config.BinSummaryStatistics = "max"
  data   := single_feature_import(config, files, false)
  result := AllocSimpleTrack("classification", data.GetGenome(), data.GetBinSize())

  // counter
  l := 0
  // total track length
  L := 0
  for _, length := range data.GetGenome().Lengths {
    L += length/config.BinSize
  }
  window_size := 10000
  offset1 := DivIntUp  (window_size, 2)
  offset2 := DivIntDown(window_size, 2)

  //w_n := float64(window_size)
  // c_n := 1.0

  pool  := threadpool.New(config.Threads, 10000)
  group := pool.NewJobGroup()

  for _, name := range data.GetSeqNames() {
    name := name
    pool.AddJob(group, func(pool threadpool.ThreadPool, erf func() error) error {
    
      slidingMed := NewSlidingMedian()
      slidingStd := NewSlidingMedian()
      seq1, err := data.GetSequence(name); if err != nil {
        log.Fatal(err)
      }
      seq2, err := result.GetSequence(name); if err != nil {
        log.Fatal(err)
      }
      nbins := seq2.NBins()

      // compute initial mean
      for i := 0; i < window_size; i++ {
        if i >= nbins {
          break
        }
        slidingMed.Insert(seq1.AtBin(i))
      }
      {
        // compute initial variance
        m := slidingMed.Median()
        for i := 0; i < window_size; i++ {
          if i >= nbins {
            break
          }
          if x := seq1.AtBin(i); x > m {
            slidingStd.Insert(x-m)
          }
        }
      }
      // loop over sequence
      for i := 0; i < nbins; i++ {
        // update mean and variance
        if i > offset1 && i < nbins-offset2 {
          x1 := seq1.AtBin(i-offset1)
          x2 := seq1.AtBin(i+offset2)
          m1 := slidingMed.Median()
          if x1 != x2 {
            slidingMed.Remove(x1)
            slidingMed.Insert(x2)
          }
          m2 := slidingMed.Median()
          if m1 != m2 {
            // recompute variance since median changed
            slidingStd = NewSlidingMedian()
            for j := i-offset1; j <= i+offset2; j++ {
              if x := seq1.AtBin(j); x > m2 {
                slidingStd.Insert(x-m2)
              }
            }
          } else {
            // update variance
            if x1 != x2 {
              slidingStd.Remove(x1-m2)
              if x2 > m1 {
                slidingStd.Insert(x2-m2)
              }
            }
          }
        }
        x := seq1.AtBin(i)
        m := slidingMed.Median()
        v := slidingStd.Median()
        if x > m {
          seq2.SetBin(i, (x - m)/(3.0*v))
          seq2.SetBin(i, 2.0/(1.0 + math.Exp(-seq2.AtBin(i))) - 1.0)
        } else {
          seq2.SetBin(i, 1e-8)
        }
        seq2.SetBin(i, math.Log(seq2.AtBin(i)))
      }
      if config.Verbose > 0 {
        l += nbins
        utility.NewProgress(L, L).PrintStderr(l)
      }
      return nil
    })
  }
  pool.Wait(group)

  return result
}

func single_feature_eval_heuristic(config ConfigModHmm, files SingleFeatureFiles, logScale bool) {

  result := single_feature_eval_heuristic_loop(config, files)

  if !logScale {
    if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
      return math.Exp(value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result, files.Foreground.Filename); err != nil {
    log.Fatal(err)
  }
  if !logScale {
    if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
      return 1.0 - value
    }); err != nil {
      log.Fatal(err)
    }
  } else {
    if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
      return utility.LogSub(0.0, value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result, files.Background.Filename); err != nil {
    log.Fatal(err)
  }
}
