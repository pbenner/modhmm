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

import . "github.com/pbenner/autodiff/statistics"
import . "github.com/pbenner/gonetics"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func multi_feature_eval_mixture_weights(config ConfigModHmm) []float64 {
  checkModelFiles(config.SingleFeatureJson)
  checkModelFiles(config.SingleFeatureComp)
  pi := []float64{}
  for _, feature := range singleFeatureList {
    filenameModel := getFieldAsString(config.SingleFeatureJson, feature)
    filenameComp  := getFieldAsString(config.SingleFeatureComp, feature)
    p, q := ImportMixtureWeights(config, filenameModel, filenameComp)
    pi = append(pi, p, q)
  }
  return pi
}

/* -------------------------------------------------------------------------- */

func get_multi_feature_model(config ConfigModHmm, state string) MatrixBatchClassifier {
  switch config.Type {
  case "": fallthrough
  case "likelihood":
    pi := multi_feature_eval_mixture_weights(config)
    switch strings.ToLower(state) {
    case "pa": return ModelPA{BasicMultiFeatureModel{pi}}
    case "pb": return ModelPB{BasicMultiFeatureModel{pi}}
    case "ea": return ModelEA{BasicMultiFeatureModel{pi}}
    case "ep": return ModelEP{BasicMultiFeatureModel{pi}}
    case "tr": return ModelTR{BasicMultiFeatureModel{pi}}
    case "tl": return ModelTL{BasicMultiFeatureModel{pi}}
    case "r1": return ModelR1{BasicMultiFeatureModel{pi}}
    case "r2": return ModelR2{BasicMultiFeatureModel{pi}}
    case "ns": return ModelNS{BasicMultiFeatureModel{pi}}
    case "cl": return ModelCL{BasicMultiFeatureModel{pi}}
    default:
      log.Fatalf("unknown state: %s", state)
    }
  case "posterior":
    switch strings.ToLower(state) {
    case "pa": return ClassifierPA{}
    case "pb": return ClassifierPB{}
    case "ea": return ClassifierEA{}
    case "ep": return ClassifierEP{}
    case "tr": return ClassifierTR{}
    case "tl": return ClassifierTL{}
    case "r1": return ClassifierR1{}
    case "r2": return ClassifierR2{}
    case "ns": return ClassifierNS{}
    case "cl": return ClassifierCL{}
    default:
      log.Fatalf("unknown state: %s", state)
    }
  default:
    log.Fatal("invalid model type `%s'", config.Type)
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func multi_feature_eval(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, tracks []Track, filenameResult string) []Track {
  if len(tracks) != len(trackFiles) {
    tracks = make([]Track, len(trackFiles))
    for i, filename := range trackFiles {
      if t, err := ImportTrack(config.SessionConfig, filename); err != nil {
        log.Fatal(err)
      } else {
        tracks[i] = t
      }
    }
  }
  result, err := BatchClassifyMultiTrack(config.SessionConfig, classifier, tracks, false); if err != nil {
    log.Fatal(err)
  }
  if err := ExportTrack(config.SessionConfig, result, filenameResult); err != nil {
    log.Fatal(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_dep(config ConfigModHmm) []string {
  files := []string{}
  for _, feature := range singleFeatureList {
    files = append(files, getFieldAsString(config.SingleFeatureFg, feature))
    files = append(files, getFieldAsString(config.SingleFeatureBg, feature))
  }
  return files
}

func modhmm_multi_feature_eval(config ConfigModHmm, state string, tracks []Track) []Track {

  if !multiFeatureList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  localConfig := config
  localConfig.BinSummaryStatistics = "mean"

  dependencies   := []string{}
  dependencies    = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies    = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  trackFiles     := modhmm_multi_feature_eval_dep(config)
  filenameResult := getFieldAsString(config.MultiFeatureProb, strings.ToUpper(state))

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_single_feature_eval_all(config)
    printStderr(config, 1, "==> Evaluating Multi-Feature Model (%s) <==\n", strings.ToUpper(state))
    classifier := get_multi_feature_model(config, state)
    tracks = multi_feature_eval(localConfig, classifier, trackFiles, tracks, filenameResult)
  }
  return tracks
}

func modhmm_multi_feature_eval_all(config ConfigModHmm) {
  var tracks []Track
  for _, feature := range multiFeatureList {
    tracks = modhmm_multi_feature_eval(config, feature, tracks)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-multi-feature", os.Args[0]))
  options.SetParameters("[STATE]\n")

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
    modhmm_multi_feature_eval_all(config)
  } else {
    modhmm_multi_feature_eval(config, options.Args()[0], nil)
  }
}
