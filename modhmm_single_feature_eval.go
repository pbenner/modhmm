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
import   "os"
import   "strings"

import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/track"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarClassifier"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"
import   "github.com/pbenner/autodiff/statistics/vectorClassifier"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func checkModelFiles(config interface{}) {
  for _, filename := range collectStrings(config) {
    if !fileExists(filename) {
      log.Fatalf(
          "ERROR: Model file `%s' required for enrichment analysis does not exist.\n" +
          "       Please download the respective file or estimate a model with the\n" +
          "       `estimate-single-feature-mixture` subcommand", filename)
    }
  }
}

/* -------------------------------------------------------------------------- */

func single_feature_eval(config ConfigModHmm, filenameModel, filenameComp, filenameData, filenameResult1, filenameResult2 string) {
  mixture := &scalarDistribution.Mixture{}

  printStderr(config, 1, "Importing mixture model from `%s'... ", filenameModel)
  if err := ImportDistribution(filenameModel, mixture, BareRealType); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  k := ImportComponents(config, filenameComp, mixture.NComponents())
  r := Components(k).Invert(mixture.NComponents())

  var scalarClassifier1 ScalarBatchClassifier
  var scalarClassifier2 ScalarBatchClassifier

  switch strings.ToLower(config.Type) {
  case "":
  case "likelihood":
    scalarClassifier1 = scalarClassifier.MixtureLikelihood{mixture, k}
    scalarClassifier2 = scalarClassifier.MixtureLikelihood{mixture, r}
  case "posterior":
    scalarClassifier1 = scalarClassifier.MixturePosterior{mixture, k}
    scalarClassifier2 = scalarClassifier.MixturePosterior{mixture, r}
  }
  vectorClassifier1 := vectorClassifier.ScalarBatchIid{scalarClassifier1, 1}
  vectorClassifier2 := vectorClassifier.ScalarBatchIid{scalarClassifier2, 1}

  if data, err := ImportTrack(config.SessionConfig, filenameData); err != nil {
    log.Fatal(err)
  } else {
    result1, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier1, data); if err != nil {
      log.Fatal(err)
    }
    if err := ExportTrack(config.SessionConfig, result1, filenameResult1); err != nil {
      log.Fatal(err)
    }

    result2, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier2, data); if err != nil {
      log.Fatal(err)
    }
    if err := ExportTrack(config.SessionConfig, result2, filenameResult2); err != nil {
      log.Fatal(err)
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_dep(config ConfigModHmm) []string {
  r := []string{}
  r  = append(r, collectStrings(config.SingleFeatureData)...)
  r  = append(r, collectStrings(config.SingleFeatureJson)...)
  r  = append(r, collectStrings(config.SingleFeatureComp)...)
  return r
}

func modhmm_single_feature_eval(config ConfigModHmm, feature string) {

  dependencies := modhmm_single_feature_eval_dep(config)

  localConfig := config
  localConfig.BinSummaryStatistics = "discrete mean"
  filenameModel   := ""
  filenameComp    := ""
  filenameData    := ""
  filenameResult1 := ""
  filenameResult2 := ""

  switch strings.ToLower(feature) {
  case "atac":
    filenameData    = config.SingleFeatureData.Atac
    filenameModel   = config.SingleFeatureJson.Atac
    filenameComp    = config.SingleFeatureComp.Atac
    filenameResult1 = config.SingleFeatureFg.Atac
    filenameResult2 = config.SingleFeatureBg.Atac
  case "h3k27ac":
    filenameData    = config.SingleFeatureData.H3k27ac
    filenameModel   = config.SingleFeatureJson.H3k27ac
    filenameComp    = config.SingleFeatureComp.H3k27ac
    filenameResult1 = config.SingleFeatureFg.H3k27ac
    filenameResult2 = config.SingleFeatureBg.H3k27ac
  case "h3k27me3":
    filenameData    = config.SingleFeatureData.H3k27me3
    filenameModel   = config.SingleFeatureJson.H3k27me3
    filenameComp    = config.SingleFeatureComp.H3k27me3
    filenameResult1 = config.SingleFeatureFg.H3k27me3
    filenameResult2 = config.SingleFeatureBg.H3k27me3
  case "h3k9me3":
    filenameData    = config.SingleFeatureData.H3k9me3
    filenameModel   = config.SingleFeatureJson.H3k9me3
    filenameComp    = config.SingleFeatureComp.H3k9me3
    filenameResult1 = config.SingleFeatureFg.H3k9me3
    filenameResult2 = config.SingleFeatureBg.H3k9me3
  case "h3k4me1":
    filenameData    = config.SingleFeatureData.H3k4me1
    filenameModel   = config.SingleFeatureJson.H3k4me1
    filenameComp    = config.SingleFeatureComp.H3k4me1
    filenameResult1 = config.SingleFeatureFg.H3k4me1
    filenameResult2 = config.SingleFeatureBg.H3k4me1
  case "h3k4me3":
    filenameData    = config.SingleFeatureData.H3k4me3
    filenameModel   = config.SingleFeatureJson.H3k4me3
    filenameComp    = config.SingleFeatureComp.H3k4me3
    filenameResult1 = config.SingleFeatureFg.H3k4me3
    filenameResult2 = config.SingleFeatureBg.H3k4me3
  case "h3k4me3o1":
    filenameData    = config.SingleFeatureData.H3k4me3o1
    filenameModel   = config.SingleFeatureJson.H3k4me3o1
    filenameComp    = config.SingleFeatureComp.H3k4me3o1
    filenameResult1 = config.SingleFeatureFg.H3k4me3o1
    filenameResult2 = config.SingleFeatureBg.H3k4me3o1
    localConfig.BinSummaryStatistics = "mean"
  case "rna":
    filenameData    = config.SingleFeatureData.Rna
    filenameModel   = config.SingleFeatureJson.Rna
    filenameComp    = config.SingleFeatureComp.Rna
    filenameResult1 = config.SingleFeatureFg.Rna
    filenameResult2 = config.SingleFeatureBg.Rna
  case "rnalow":
    filenameData    = config.SingleFeatureData.Rna
    filenameModel   = config.SingleFeatureJson.RnaLow
    filenameComp    = config.SingleFeatureComp.RnaLow
    filenameResult1 = config.SingleFeatureFg.RnaLow
    filenameResult2 = config.SingleFeatureBg.RnaLow
  case "control":
    filenameData    = config.SingleFeatureData.Control
    filenameModel   = config.SingleFeatureJson.Control
    filenameComp    = config.SingleFeatureComp.Control
    filenameResult1 = config.SingleFeatureFg.Control
    filenameResult2 = config.SingleFeatureBg.Control
  default:
    log.Fatalf("unknown feature: %s", feature)
  }
  if updateRequired(config, filenameResult1, dependencies...) ||
    (updateRequired(config, filenameResult2, dependencies...)) {
    checkModelFiles(config.SingleFeatureJson)
    checkModelFiles(config.SingleFeatureComp)

    modhmm_single_feature_coverage_all(config)
    printStderr(config, 1, "==> Computing Single-Feature Classification (%s) <==\n", feature)
    single_feature_eval(localConfig, filenameModel, filenameComp, filenameData, filenameResult1, filenameResult2)
  }
}

func modhmm_single_feature_eval_all(config ConfigModHmm) {
  for _, feature := range []string{"atac", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "rnaLow", "control"} {
    modhmm_single_feature_eval(config, feature)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-single-feature", os.Args[0]))
  options.SetParameters("<FEATURE>\n")

  optHelp := options.   BoolLong("help",     'h',     "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  // command arguments
  if len(options.Args()) > 1 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }

  if len(options.Args()) == 0 {
    modhmm_single_feature_eval_all(config)
  } else {
    modhmm_single_feature_eval(config, options.Args()[0])
  }
}