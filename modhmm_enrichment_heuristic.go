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

func enrichment_import_heuristic(config ConfigModHmm, files EnrichmentFiles) Track {
  config.BinSummaryStatistics = "discrete mean"
  return enrichment_import_and_normalize(config, files.Coverage.Filename, files.CoverageCnts.Filename, false)
}

/* -------------------------------------------------------------------------- */

func compute_sigmoid_parameters(x1, x2, p1, p2 float64) (float64, float64) {

  sigmoid := func(r Scalar, x, a, b ConstScalar) {
    r.Mul(a, x)
    r.Add(r, b)
    r.Neg(r)
    r.Exp(r)
    r.Add(r, ConstFloat64(1.0))
    r.Div(ConstFloat64(1.0), r)
  }
  generator := func(x1, x2 float64) func(ConstVector) (MagicScalar, error) {
    f := func(a ConstVector) (MagicScalar, error) {
      r := NullDenseReal64Vector(2)
      sigmoid(r.At(0), ConstFloat64(x1), a.ConstAt(0), a.ConstAt(1))
      sigmoid(r.At(1), ConstFloat64(x2), a.ConstAt(0), a.ConstAt(1))
      r.At(0).Sub(r.ConstAt(0), ConstFloat64(p1))
      r.At(1).Sub(r.ConstAt(1), ConstFloat64(p2))
      r.At(0).Mul(r.ConstAt(0), r.ConstAt(0))
      r.At(1).Mul(r.ConstAt(1), r.ConstAt(1))
      r.At(0).Add(r.ConstAt(0), r.ConstAt(1))
      return r.MagicAt(0), nil
    }
    return f
  }
  objective := generator(x1, x2)

  if x, err := rprop.Run(objective, NewDenseFloat64Vector([]float64{0.01,0.01}), 0.01, []float64{1.05,0.95}, rprop.Epsilon{1e-10}); err != nil {
    panic(err)
  } else {
    return x.Float64At(0), x.Float64At(1)
  }
}

/* -------------------------------------------------------------------------- */

func enrichment_eval_heuristic_loop(config ConfigModHmm, result MutableTrack, data Track, a, b float64) {
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

func enrichment_eval_heuristic_parameters(config ConfigModHmm, files EnrichmentFiles, counts Counts) (float64, float64) {
  if files.Feature == "rna" {
    q := config.EnrichmentParameters.GetParameters(files.Feature)[0]
    p := config.EnrichmentParameters.GetParameters(files.Feature)[1]
    m1 := counts.Quantile(0.0)
    m2 := counts.Quantile(q)
    return compute_sigmoid_parameters(m1, m2, 0.01, p)
  } else {
    q  := config.EnrichmentParameters.GetParameters(files.Feature)[0]
    p1 := config.EnrichmentParameters.GetParameters(files.Feature)[1]
    p2 := config.EnrichmentParameters.GetParameters(files.Feature)[2]
    m1 := counts.Quantile(q)
    m2 := counts.ThresholdedMean(m1)
    return compute_sigmoid_parameters(m1, m2, p1, p2)
  }
}

func enrichment_eval_heuristic(config ConfigModHmm, files EnrichmentFiles) {
  data   := enrichment_import_heuristic(config, files)
  counts := compute_counts(config, data)
  result := AllocSimpleTrack("classification", data.GetGenome(), data.GetBinSize())
  a, b   := enrichment_eval_heuristic_parameters(config, files, counts)

  enrichment_eval_heuristic_loop(config, result, data, a, b)

  if err := ExportTrack(config.SessionConfig, result, files.Probabilities.Filename); err != nil {
    log.Fatal(err)
  }
}
