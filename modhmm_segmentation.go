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
import . "github.com/pbenner/ngstat/estimation"
import . "github.com/pbenner/ngstat/track"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/matrixClassifier"
import   "github.com/pbenner/autodiff/statistics/matrixDistribution"
import   "github.com/pbenner/autodiff/statistics/matrixEstimator"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func estimate(config ConfigModHmm, trackFiles []string, model string) {
  var estimator  *matrixEstimator.HmmEstimator
  var stateNames []string

  switch model {
  case "default":
    estimator, stateNames = getModHmmDefaultEstimator(config)
  case "dense":
    estimator, stateNames = getModHmmDenseEstimator(config)
  default:
    log.Fatalf("ERROR: invalid model name `%s'", model)
  }

  if err := ImportAndEstimateOnMultiTrack(config.SessionConfig, estimator, trackFiles, true); err != nil {
    log.Fatalf("ERROR: %s", err)
  }
  modhmm := ModHmm{}
  modhmm.Hmm        = *estimator.GetEstimate().(*matrixDistribution.Hmm)
  modhmm.StateNames = stateNames

  printStderr(config, 1, "Exporting model to `%s'... ", config.Model)
  if err := ExportDistribution(config.Model, &modhmm); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("ERROR: %s", err)
  }
  printStderr(config, 1, "done\n")
}

/* -------------------------------------------------------------------------- */

func segment(config ConfigModHmm, trackFiles []string) {
  modhmm := ModHmm{}
  printStderr(config, 1, "Importing model from `%s'... ", config.Model)
  if err := ImportDistribution(config.Model, &modhmm, BareRealType); err != nil {
    log.Fatal(err)
    printStderr(config, 1, "failed\n")
  }
  printStderr(config, 1, "done\n")

  if result, err := ImportAndClassifyMultiTrack(config.SessionConfig, matrixClassifier.HmmClassifier{&modhmm.Hmm}, trackFiles, true); err != nil {
    log.Fatal(err)
  } else {
    var name, desc string
    if config.Description == "" {
      name = "ModHMM"
      desc = "Segmentation ModHMM"
    } else {
      name = fmt.Sprintf("ModHMM [%s]", config.Description)
      desc = fmt.Sprintf("Segmentation ModHMM [%s]", config.Description)
    }
    printStderr(config, 1, "Writing genome segmentation to `%s'... ", config.Segmentation)
    if err := ExportTrackSegmentation(config.SessionConfig, result, config.Segmentation, name, desc, true, modhmm.StateNames, nil); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_segmentation_dep(config ConfigModHmm) []string {
  files := make([]string, len(multiFeatureList))
  for i, state := range multiFeatureList {
    files[i] =  getFieldAsString(config.MultiFeatureProb, strings.ToUpper(state))
  }
  return files
}

func modhmm_segmentation(config ConfigModHmm, model string) {

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_segmentation_dep(config)...)

  trackFiles := modhmm_segmentation_dep(config)

  filenameModel        := config.Model
  filenameSegmentation := config.Segmentation

  if updateRequired(config, filenameModel, dependencies...) {
    modhmm_multi_feature_eval_all(config)

    printStderr(config, 1, "==> Estimating ModHmm transition parameters <==\n")
    estimate(config, trackFiles, model)
  }
  if updateRequired(config, filenameSegmentation, append(dependencies, filenameModel)...) {
    printStderr(config, 1, "==> Computing Segmentation <==\n")
    segment(config, trackFiles)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_segmentation_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s segmentation", os.Args[0]))
  options.SetParameters("<STATE>\n")

  optHelp  := options.   BoolLong("help",     'h',            "print help")
  optModel := options. StringLong("model",     0 , "default", "default, dense")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  // command arguments
  if len(options.Args()) > 0 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }

  modhmm_segmentation(config, *optModel)
}
