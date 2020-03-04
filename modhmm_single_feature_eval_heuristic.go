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
import   "github.com/pbenner/threadpool"

import . "github.com/pbenner/modhmm/config"

import . "github.com/pbenner/autodiff"
import   "github.com/pbenner/autodiff/algorithm/rprop"

/* -------------------------------------------------------------------------- */

func single_feature_import_heuristic(config ConfigModHmm, files SingleFeatureFiles) Track {
  if files.Feature == "h3k4me3o1" {
    config.BinSummaryStatistics = "discrete mean"
    config.BinOverlap = 1
    track1 := single_feature_import_and_normalize(config, files.SrcCoverage[0].Filename, files.SrcCoverageCnts[0].Filename, false)
    track2 := single_feature_import_and_normalize(config, files.SrcCoverage[1].Filename, files.SrcCoverageCnts[1].Filename, false)
    return single_feature_compute_h3k4me3o1(config, track1, track2)
  } else {
    config.BinSummaryStatistics = "discrete mean"
    return single_feature_import_and_normalize(config, files.Coverage.Filename, files.CoverageCnts.Filename, false)
  }
}

/* -------------------------------------------------------------------------- */

func compute_sigmoid_parameters(x1, x2, p1, p2 float64) (float64, float64) {

  sigmoid := func(r Scalar, x, a, b ConstScalar) {
    r.Mul(a, x)
    r.Add(r, b)
    r.Neg(r)
    r.Exp(r)
    r.Add(r, ConstReal(1.0))
    r.Div(ConstReal(1.0), r)
  }
  generator := func(x1, x2 float64) func(Vector) (Scalar, error) {
    f := func(a Vector) (Scalar, error) {
      r := NullDenseRealVector(2)
      sigmoid(r[0], ConstReal(x1), a.At(0), a.At(1))
      sigmoid(r[1], ConstReal(x2), a.At(0), a.At(1))
      r[0].Sub(r[0], ConstReal(p1))
      r[1].Sub(r[1], ConstReal(p2))
      r[0].Mul(r[0], r[0])
      r[1].Mul(r[1], r[1])
      r[0].Add(r[0], r[1])
      return r[0], nil
    }
    return f
  }
  objective := generator(x1, x2)

  if x, err := rprop.Run(objective, NewDenseRealVector([]float64{0.01,0.01}), 0.01, []float64{1.1,0.9}, rprop.Epsilon{1e-10}); err != nil {
    panic(err)
  } else {
    return x.ValueAt(0), x.ValueAt(1)
  }
}

/* -------------------------------------------------------------------------- */

func single_feature_eval_heuristic_loop(config ConfigModHmm, result MutableTrack, data Track, a, b float64) {
  pool  := threadpool.New(config.Threads, 10000)
  group := pool.NewJobGroup()

  for _, name := range data.GetSeqNames() {
    name := name
    pool.AddJob(group, func(pool threadpool.ThreadPool, erf func() error) error {
    
      seq1, err := data.GetSequence(name); if err != nil {
        log.Fatal(err)
      }
      seq2, err := result.GetSequence(name); if err != nil {
        log.Fatal(err)
      }
      nbins := seq2.NBins()

      // loop over sequence
      for i := 0; i < nbins; i++ {
        x := seq1.AtBin(i)
        // apply logistic function
        seq2.SetBin(i, 1.0/(1.0 + math.Exp(-a*x-b)))
      }
      return nil
    })
  }
  pool.Wait(group)
}

func single_feature_eval_heuristic_parameters(config ConfigModHmm, files SingleFeatureFiles, counts Counts) (float64, float64) {
  if files.Feature == "rna" {
    q  := 0.20
    p1 := 0.01
    p2 := 0.50
    m1 := counts.X[0]
    m2 := counts.Quantile(q)
    return compute_sigmoid_parameters(m1, m2, p1, p2)
  } else {
    // default parameters
    q  := 0.80
    p1 := 0.01
    p2 := 0.50
    // update parameters
    switch files.Feature {
    case "control" : q = 0.95
    case "h3k27me3": q = 0.90
    case "h3k9me3" : q = 0.90
    }
    m1 := counts.Quantile(q)
    m2 := counts.ThresholdedMean(m1)
    return compute_sigmoid_parameters(m1, m2, p1, p2)
  }
}

func single_feature_eval_heuristic(config ConfigModHmm, files SingleFeatureFiles) {
  data   := single_feature_import_heuristic(config, files)
  counts := compute_counts(config, data)
  a, b   := single_feature_eval_heuristic_parameters(config, files, counts)
  result := AllocSimpleTrack("classification", data.GetGenome(), data.GetBinSize())

  // compute probabilities
  single_feature_eval_heuristic_loop(config, result, data, a, b)

  // rna-low is a special case
  if files.Feature == "rna" {
    single_feature_eval_rna_low(config, result, data)
  }
  if err := ExportTrack(config.SessionConfig, result, files.Probabilities.Filename); err != nil {
    log.Fatal(err)
  }
}