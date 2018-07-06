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

import   "fmt"
import   "log"
import   "math/rand"
import   "os"
import   "strconv"
import   "strings"

import . "github.com/pbenner/ngstat/estimation"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"
import   "github.com/pbenner/autodiff/statistics/scalarEstimator"
import   "github.com/pbenner/autodiff/statistics/vectorDistribution"
import   "github.com/pbenner/autodiff/statistics/vectorEstimator"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func newEstimator(config ConfigModHmm, n_delta, n_poisson, n_geometric int) VectorEstimator {
  components := []ScalarEstimator{}
  for i := 0; i < n_delta; i++ {
    if delta, err := scalarEstimator.NewDeltaEstimator(float64(i)); err != nil {
      log.Fatal(err)
    } else {
      components = append(components, delta)
    }
  }
  for i := 0; i < n_poisson; i++ {
    if poisson, err := scalarEstimator.NewPoissonEstimator(rand.Float64()); err != nil {
      log.Fatal(err)
    } else {
      if t, err := scalarEstimator.NewTranslationEstimator(poisson, -float64(n_delta)); err != nil {
        log.Fatal(err)
      } else {
        components = append(components, t)
      }
    }
  }
  for i := 0; i < n_geometric; i++ {
    if geometric, err := scalarEstimator.NewGeometricEstimator(0.01*rand.Float64()); err != nil {
      log.Fatal(err)
    } else {
      components = append(components, geometric)
    }
  }
  if mixture, err := scalarEstimator.NewDiscreteMixtureEstimator(nil, components, 1e-8, -1); err != nil {
    log.Fatal(err)
  } else {
    if estimator, err := vectorEstimator.NewScalarIid(mixture, -1); err != nil {
      log.Fatal(err)
    } else {
      return estimator
    }
  }
  return nil
}

func newContinuousEstimator(config ConfigModHmm, n_lognormal, n_exponential int) VectorEstimator {
  components := []ScalarEstimator{}
  for i := 0; i < n_lognormal; i++ {
    if normal, err := scalarEstimator.NewNormalEstimator(rand.Float64(), 1.0, 1e-8); err != nil {
      log.Fatal(err)
    } else {
      if t, err := scalarEstimator.NewLogTransformEstimator(normal, 1.0); err != nil {
        log.Fatal(err)
      } else {
        components = append(components, t)
      }
    }
  }
  for i := 0; i < n_exponential; i++ {
    if exponential, err := scalarEstimator.NewExponentialEstimator(0.04, 1e4); err != nil {
      log.Fatal(err)
    } else {
      components = append(components, exponential)
    }
  }
  if mixture, err := scalarEstimator.NewDiscreteMixtureEstimator(nil, components, 1e-8, -1); err != nil {
    log.Fatal(err)
  } else {
    // set options
    mixture.Verbose = config.Verbose
    if estimator, err := vectorEstimator.NewScalarIid(mixture, -1); err != nil {
      log.Fatal(err)
    } else {
      return estimator
    }
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func estimate(config ConfigModHmm, feature string, n []int) {
  var estimator VectorEstimator

  filenameIn  := ""
  filenameOut := ""
  discrete    := true

  switch strings.ToLower(feature) {
  case "atac":
    filenameIn  = config.SingleFeatureData.atac
    filenameOut = config.SingleFeatureJson.atac
  case "h3k27ac":
    filenameIn  = config.SingleFeatureData.h3k27ac
    filenameOut = config.SingleFeatureJson.h3k27ac
  case "h3k27me3":
    filenameIn  = config.SingleFeatureData.h3k27me3
    filenameOut = config.SingleFeatureJson.h3k27me3
  case "h3k9me3":
    filenameIn  = config.SingleFeatureData.h3k9me3
    filenameOut = config.SingleFeatureJson.h3k9me3
  case "h3k4me1":
    filenameIn  = config.SingleFeatureData.h3k4me1
    filenameOut = config.SingleFeatureJson.h3k4me1
  case "h3k4me3":
    filenameIn  = config.SingleFeatureData.h3k4me3
    filenameOut = config.SingleFeatureJson.h3k4me3
  case "h3k4me3o1":
    filenameIn  = config.SingleFeatureData.h3k4me3o1
    filenameOut = config.SingleFeatureJson.h3k4me3o1
    discrete    = false
  default:
    log.Fatal("unknown feature: %s", feature)
  }

  if  discrete && len(n) != 3 {
    log.Fatalf("feature `%s' is discrete and requires three integer arguments", feature)
  }
  if !discrete && len(n) != 2 {
    log.Fatalf("feature `%s' is continuous and requires two integer arguments", feature)
  }
  if discrete {
    estimator = newEstimator(config, n[0], n[1], n[2])
  } else {
    estimator = newContinuousEstimator(config, n[0], n[1])
  }

  localConfig := config.SessionConfig
  localConfig.BinSummaryStatistics = "discrete mean"

  if err := ImportAndEstimateOnSingleTrack(localConfig, estimator, filenameIn); err != nil {
    log.Fatal(err)
  }
  result := estimator.GetEstimate().(*vectorDistribution.ScalarIid).Distribution.(*scalarDistribution.Mixture)

  ExportDistribution(filenameOut, result)
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_estimate_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s estimate-single-feature-mixture", os.Args[0]))
  options.SetParameters("<FEATURE> <N_DELTA> <N_POISSON> <N_GEOMETRIC>\n")
  options.SetParameters("<FEATURE> <N_LOGNORMAL> <N_EXPONENTIAL>\n")

  optHelp := options.   BoolLong("help",     'h',     "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  // command arguments
  if len(options.Args()) != 3 && len(options.Args()) != 4 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }

  feature := options.Args()[0]
  n       := []int{}

  if m, err := strconv.ParseInt(options.Args()[1], 10, 64); err != nil {
    log.Fatal(err)
  } else {
    n = append(n, int(m))
  }
  if m, err := strconv.ParseInt(options.Args()[2], 10, 64); err != nil {
    log.Fatal(err)
  } else {
    n = append(n, int(m))
  }
  if len(options.Args()) == 4 {
    if m, err := strconv.ParseInt(options.Args()[3], 10, 64); err != nil {
      log.Fatal(err)
    } else {
      n = append(n, int(m))
    }
  }

  estimate(config, feature, n)
}
