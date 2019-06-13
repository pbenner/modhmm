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
import   "math"
import   "os"
import   "strings"

import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/track"

import   "github.com/pbenner/autodiff/statistics/scalarClassifier"
import   "github.com/pbenner/autodiff/statistics/vectorClassifier"

import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

/* default single-feature-model
 * -------------------------------------------------------------------------- */

//go:generate go run modhmm_single_feature_default.gen.go

/* -------------------------------------------------------------------------- */

func single_feature_eval(config ConfigModHmm, filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1, filenameResult2 string, logScale bool) {
  mixture := ImportMixtureDistribution(config, filenameModel)
  counts  := Counts{}

  k, r := ImportComponents(config, filenameComp, mixture.NComponents())

  printStderr(config, 1, "Importing reference counts from `%s'... ", filenameCnts)
  if err := counts.ImportFile(filenameCnts); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  scalarClassifier1 := scalarClassifier.MixturePosterior{mixture, k}
  scalarClassifier2 := scalarClassifier.MixturePosterior{mixture, r}
  vectorClassifier1 := vectorClassifier.ScalarBatchIid{scalarClassifier1, 1}
  vectorClassifier2 := vectorClassifier.ScalarBatchIid{scalarClassifier2, 1}

  if data, err := ImportTrack(config.SessionConfig, filenameData); err != nil {
    log.Fatal(err)
  } else {
    printStderr(config, 1, "Quantile normalizing track to reference distribution... ")
    if err := (GenericMutableTrack{data}).QuantileNormalizeToCounts(counts.X, counts.Y); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")

    result1, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier1, data); if err != nil {
      log.Fatal(err)
    }
    if !logScale {
      if err := (GenericMutableTrack{result1}).Map(result1, func(seqname string, position int, value float64) float64 {
        return math.Exp(value)
      }); err != nil {
        log.Fatal(err)
      }
    }
    if err := ExportTrack(config.SessionConfig, result1, filenameResult1); err != nil {
      log.Fatal(err)
    }

    result2, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier2, data); if err != nil {
      log.Fatal(err)
    }
    if !logScale {
      if err := (GenericMutableTrack{result2}).Map(result2, func(seqname string, position int, value float64) float64 {
        return math.Exp(value)
      }); err != nil {
        log.Fatal(err)
      }
    }
    if err := ExportTrack(config.SessionConfig, result2, filenameResult2); err != nil {
      log.Fatal(err)
    }
  }
}

func single_feature_files(config ConfigModHmm, feature string, logScale bool) (TargetFile, TargetFile, TargetFile, TargetFile, TargetFile, TargetFile) {

  if !SingleFeatureList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
  filenameModel   := TargetFile{}
  filenameComp    := TargetFile{}
  filenameData    := TargetFile{}
  filenameCnts    := TargetFile{}
  filenameResult1 := TargetFile{}
  filenameResult2 := TargetFile{}

  switch strings.ToLower(feature) {
  case "rna-low":
    filenameData    = config.Coverage          .Rna
    filenameCnts    = config.CoverageCnts      .Rna
    filenameModel   = config.SingleFeatureModel.Rna
    filenameComp    = config.SingleFeatureComp .Rna_low
    if logScale {
      filenameResult1 = config.SingleFeatureFg.Rna_low
      filenameResult2 = config.SingleFeatureBg.Rna_low
    } else {
      filenameResult1 = config.SingleFeatureFgExp.Rna_low
      filenameResult2 = config.SingleFeatureBgExp.Rna_low
    }
  default:
    filenameData    = config.Coverage          .GetTargetFile(feature)
    filenameCnts    = config.CoverageCnts      .GetTargetFile(feature)
    filenameModel   = config.SingleFeatureModel.GetTargetFile(feature)
    filenameComp    = config.SingleFeatureComp .GetTargetFile(feature)
    if logScale {
      filenameResult1 = config.SingleFeatureFg.GetTargetFile(feature)
      filenameResult2 = config.SingleFeatureBg.GetTargetFile(feature)
    } else {
      filenameResult1 = config.SingleFeatureFgExp.GetTargetFile(feature)
      filenameResult2 = config.SingleFeatureBgExp.GetTargetFile(feature)
    }
  }
  return filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1, filenameResult2
}

func single_feature_filter_update(config ConfigModHmm, features []string, logScale bool) []string {
  r := []string{}
  for _, feature := range features {
    feature = config.CoerceOpenChromatinAssay(feature)
    filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1, filenameResult2 :=
      single_feature_files(config, feature, logScale)
    if updateRequired(config, filenameResult1, filenameData.Filename, filenameCnts.Filename, filenameModel.Filename, filenameComp.Filename) ||
      (updateRequired(config, filenameResult2, filenameData.Filename, filenameCnts.Filename, filenameModel.Filename, filenameComp.Filename)) {
      r = append(r, feature)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_dep(config ConfigModHmm) []string {
  r := []string{}
  r  = append(r, config.Coverage.GetFilenames()...)
  r  = append(r, config.SingleFeatureModel.GetFilenames()...)
  r  = append(r, config.SingleFeatureComp.GetFilenames()...)
  r  = append(r, config.CoverageCnts.GetFilenames()...)
  return r
}

func modhmm_single_feature_eval(config ConfigModHmm, feature string, logScale bool) {

  if !SingleFeatureList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
  filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1, filenameResult2 :=
    single_feature_files(config, feature, logScale)

  localConfig := config
  localConfig.BinSummaryStatistics = "discrete mean"
  if updateRequired(config, filenameResult1, filenameData.Filename, filenameCnts.Filename, filenameModel.Filename, filenameComp.Filename) ||
    (updateRequired(config, filenameResult2, filenameData.Filename, filenameCnts.Filename, filenameModel.Filename, filenameComp.Filename)) {

    // estimate single feature model if required
    // modhmm_single_feature_estimate_default(config, feature)
    // if CoverageList.Contains(strings.ToLower(feature)) {
    //   modhmm_compute_counts(config, feature)
    // }
    printStderr(config, 1, "==> Evaluating Single-Feature Model (%s) <==\n", feature)
    single_feature_eval(localConfig, filenameModel.Filename, filenameComp.Filename, filenameData.Filename, filenameCnts.Filename, filenameResult1.Filename, filenameResult2.Filename, logScale)
  }
}

func modhmm_single_feature_eval_loop(config ConfigModHmm, features []string, logScale bool) {
  // reduce list of features to those that require an update
  features = single_feature_filter_update(config, features, logScale)
  // compute coverages here to make use of multi-threading
  modhmm_coverage_loop(config, InsensitiveStringList(features).Intersection(CoverageList))
  // eval single features
  for _, feature := range features {
    modhmm_single_feature_eval(config, feature, logScale)
  }
}

func modhmm_single_feature_eval_all(config ConfigModHmm, logScale bool) {
  modhmm_single_feature_eval_loop(config, SingleFeatureList, logScale)
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optStdScale := options.BoolLong("std-scale",  0 ,  "single-feature output on standard scale")
  optHelp     := options.BoolLong("help",      'h',  "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_single_feature_eval_all(config, !*optStdScale)
  } else {
    modhmm_single_feature_eval_loop(config, options.Args(), !*optStdScale)
  }
}
