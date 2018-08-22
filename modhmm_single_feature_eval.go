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

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarClassifier"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"
import   "github.com/pbenner/autodiff/statistics/vectorClassifier"

import . "github.com/pbenner/gonetics"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func checkModelFiles(filenames []string) {
  for _, filename := range filenames {
    if !fileExists(filename) {
      log.Fatalf(
          "ERROR: Model file `%s' required for enrichment analysis does not exist.\n" +
          "       Please download the respective file or estimate a model with the\n" +
          "       `estimate-single-feature` subcommand", filename)
    }
  }
}

/* -------------------------------------------------------------------------- */

func single_feature_eval(config ConfigModHmm, filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1, filenameResult2 string, logScale bool) {
  mixture := &scalarDistribution.Mixture{}
  counts  := Counts{}

  printStderr(config, 1, "Importing mixture model from `%s'... ", filenameModel)
  if err := ImportDistribution(filenameModel, mixture, BareRealType); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  k := ImportComponents(config, filenameComp, mixture.NComponents())
  r := Components(k).Invert(mixture.NComponents())

  printStderr(config, 1, "Importing reference counts from `%s'... ", filenameCnts)
  if err := counts.ImportFile(filenameCnts); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  var scalarClassifier1 ScalarBatchClassifier
  var scalarClassifier2 ScalarBatchClassifier

  switch strings.ToLower(config.ModelType) {
  case "likelihood":
    scalarClassifier1 = scalarClassifier.MixtureLikelihood{mixture, k}
    scalarClassifier2 = scalarClassifier.MixtureLikelihood{mixture, r}
  case "posterior":
    scalarClassifier1 = scalarClassifier.MixturePosterior{mixture, k}
    scalarClassifier2 = scalarClassifier.MixturePosterior{mixture, r}
  default:
    log.Fatalf("invalid model type `%s'", config.ModelType)
  }
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

  if !singleFeatureList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }

  localConfig := config
  localConfig.BinSummaryStatistics = "discrete mean"
  filenameModel   := ""
  filenameComp    := ""
  filenameData    := ""
  filenameCnts    := ""
  filenameResult1 := TargetFile{}
  filenameResult2 := TargetFile{}

  switch strings.ToLower(feature) {
  case "rna-low":
    filenameData    = config.Coverage.Rna.Filename
    filenameCnts    = config.CoverageCnts.Rna.Filename
    filenameModel   = config.SingleFeatureModel.Rna_low.Filename
    filenameComp    = config.SingleFeatureComp.Rna_low.Filename
    if logScale {
      filenameResult1 = config.SingleFeatureFg.Rna_low
      filenameResult2 = config.SingleFeatureBg.Rna_low
    } else {
      filenameResult1 = config.SingleFeatureFgExp.Rna_low
      filenameResult2 = config.SingleFeatureBgExp.Rna_low
    }
  default:
    filenameData    = config.Coverage          .GetTargetFile(feature).Filename
    filenameCnts    = config.CoverageCnts      .GetTargetFile(feature).Filename
    filenameModel   = config.SingleFeatureModel.GetTargetFile(feature).Filename
    filenameComp    = config.SingleFeatureComp .GetTargetFile(feature).Filename
    if logScale {
      filenameResult1 = config.SingleFeatureFg.GetTargetFile(feature)
      filenameResult2 = config.SingleFeatureBg.GetTargetFile(feature)
    } else {
      filenameResult1 = config.SingleFeatureFgExp.GetTargetFile(feature)
      filenameResult2 = config.SingleFeatureBgExp.GetTargetFile(feature)
    }
  }
  if updateRequired(config, filenameResult1, filenameData, filenameCnts, filenameModel, filenameComp) ||
    (updateRequired(config, filenameResult2, filenameData, filenameCnts, filenameModel, filenameComp)) {
    checkModelFiles(config.SingleFeatureModel.GetFilenames())
    checkModelFiles(config.SingleFeatureComp .GetFilenames())
    checkModelFiles(config.CoverageCnts      .GetFilenames())

    modhmm_coverage_all(config)
    printStderr(config, 1, "==> Evaluating Single-Feature Model (%s) <==\n", feature)
    single_feature_eval(localConfig, filenameModel, filenameComp, filenameData, filenameCnts, filenameResult1.Filename, filenameResult2.Filename, logScale)
  }
}

func modhmm_single_feature_eval_loop(config ConfigModHmm, states []string, logScale bool) {
  for _, feature := range states {
    modhmm_single_feature_eval(config, feature, logScale)
  }
}

func modhmm_single_feature_eval_all(config ConfigModHmm, logScale bool) {
  modhmm_single_feature_eval_loop(config, singleFeatureList, logScale)
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
