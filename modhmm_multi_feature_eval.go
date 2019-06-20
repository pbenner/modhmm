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

import . "github.com/pbenner/autodiff/statistics"
import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func get_multi_feature_model(config ConfigModHmm, state string) MatrixBatchClassifier {
  switch strings.ToLower(state) {
  case "pa": return ClassifierPA{}
  case "ea": return ClassifierEA{}
  case "bi": return ClassifierBI{}
  case "pr": return ClassifierPR{}
  case "tr": return ClassifierTR{}
  case "tl": return ClassifierTL{}
  case "r1": return ClassifierR1{}
  case "r2": return ClassifierR2{}
  case "ns": return ClassifierNS{}
  case "cl": return ClassifierCL{}
  default:
    log.Fatalf("unknown state: %s", state)
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func multi_feature_eval(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, tracks []Track, filenameResult string, logScale bool) []Track {
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
  if !logScale {
    if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
      return math.Exp(value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result, filenameResult); err != nil {
    log.Fatal(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_dep(config ConfigModHmm) []string {
  files := []string{}
  for _, feature := range SingleFeatureList {
    files = append(files, config.SingleFeatureFg.GetTargetFile(feature).Filename)
    files = append(files, config.SingleFeatureBg.GetTargetFile(feature).Filename)
  }
  return files
}

func modhmm_multi_feature_eval(config ConfigModHmm, state string, tracks []Track, logScale bool) []Track {

  if !MultiFeatureList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  localConfig := config
  localConfig.BinSummaryStatistics = "mean"

  dependencies   := []string{}
  dependencies    = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  dependencies    = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  trackFiles     := modhmm_multi_feature_eval_dep(config)
  filenameResult := TargetFile{}
  if logScale {
    filenameResult = config.MultiFeatureProb.GetTargetFile(state)
  } else {
    filenameResult = config.MultiFeatureProbExp.GetTargetFile(state)
  }

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_single_feature_eval_all(config, true)
    printStderr(config, 1, "==> Evaluating Multi-Feature Model (%s) <==\n", strings.ToUpper(state))
    classifier := get_multi_feature_model(config, state)
    tracks = multi_feature_eval(localConfig, classifier, trackFiles, tracks, filenameResult.Filename, logScale)
  }
  return tracks
}

func modhmm_multi_feature_eval_loop(config ConfigModHmm, states []string, logScale bool) {
  var tracks []Track
  for _, feature := range states {
    tracks = modhmm_multi_feature_eval(config, feature, tracks, logScale)
  }
}

func modhmm_multi_feature_eval_all(config ConfigModHmm, logScale bool) {
  modhmm_multi_feature_eval_loop(config, MultiFeatureList, logScale)
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-multi-feature", os.Args[0]))
  options.SetParameters("[STATE]...\n")

  optNormalize := options.BoolLong("normalize",  0 ,  "normalize multi-feature likelihoods")
  optStdScale  := options.BoolLong("std-scale",  0 ,  "multi-feature output on standard scale")
  optHelp      := options.BoolLong("help",      'h',  "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if *optNormalize {
    if len(options.Args()) == 0 {
      modhmm_multi_feature_eval_norm_all(config, !*optStdScale)
    } else {
      modhmm_multi_feature_eval_norm_loop(config, options.Args(), !*optStdScale)
    }
  } else {
    if len(options.Args()) == 0 {
      modhmm_multi_feature_eval_all(config, !*optStdScale)
    } else {
      modhmm_multi_feature_eval_loop(config, options.Args(), !*optStdScale)
    }
  }
}
