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
import . "github.com/pbenner/autodiff/logarithmetic"
import . "github.com/pbenner/autodiff/statistics"

import . "github.com/pbenner/gonetics"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

type normalizationClassifier struct {
  k, n int
}

func (obj normalizationClassifier) Eval(s Scalar, x ConstMatrix) error {
  r := math.Inf(-1)

  for i := 0; i < obj.n; i++ {
    r = LogAdd(r, x.ValueAt(i, 0))
  }
  r = x.ValueAt(obj.k, 0) - r
  r = math.Exp(r)

  s.SetValue(r); return nil
}

func (obj normalizationClassifier) Dims() (int, int) {
  return obj.n, 1
}

func (obj normalizationClassifier) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return normalizationClassifier{obj.k, obj.n}
}

/* -------------------------------------------------------------------------- */

func multi_feature_eval_norm_dep(config ConfigModHmm, state string, trackFiles []string, tracks []Track, filenameResult string) []Track {
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
  classifier := normalizationClassifier{multiFeatureList.Index(state), len(tracks)}

  result, err := BatchClassifyMultiTrack(config.SessionConfig, classifier, tracks, false); if err != nil {
    log.Fatal(err)
  }
  if err := ExportTrack(config.SessionConfig, result, filenameResult); err != nil {
    log.Fatal(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_norm_dep(config ConfigModHmm) []string {
  return modhmm_segmentation_dep(config)
}

func modhmm_multi_feature_eval_norm(config ConfigModHmm, state string, tracks []Track) []Track {

  if !multiFeatureList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_norm_dep(config)...)

  trackFiles := modhmm_multi_feature_eval_norm_dep(config)
  filenameResult := getFieldAsString(config.MultiFeatureProbNorm, strings.ToUpper(state))

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_multi_feature_eval_all(config)

    printStderr(config, 1, "==> Evaluating Normalized Multi-Feature Model (%s) <==\n", strings.ToUpper(state))
    tracks = multi_feature_eval_norm_dep(config, state, trackFiles, tracks, filenameResult)
  }
  return tracks
}

func modhmm_multi_feature_eval_norm_all(config ConfigModHmm) {
  var tracks []Track
  for _, state := range multiFeatureList {
    tracks = modhmm_multi_feature_eval_norm(config, state, tracks)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_norm_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s multi-feature-eval-norm", os.Args[0]))
  options.SetParameters("[STATE]\n")

  optHelp  := options.   BoolLong("help",     'h',            "print help")

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
  if len(options.Args()) == 0 {
    modhmm_multi_feature_eval_norm_all(config)
  } else {
    modhmm_multi_feature_eval_norm(config, options.Args()[0], nil)
  }
}
