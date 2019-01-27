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

//import   "fmt"
import   "log"
import   "math"
import   "strings"

import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/track"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/logarithmetic"
import . "github.com/pbenner/autodiff/statistics"

import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

type normalizationClassifier struct {
  k, n     int
  logScale bool
}

func (obj normalizationClassifier) Eval(s Scalar, x ConstMatrix) error {
  r := math.Inf(-1)

  for i := 0; i < obj.n; i++ {
    r = LogAdd(r, x.ValueAt(i, 0))
  }
  r = x.ValueAt(obj.k, 0) - r
  if !obj.logScale {
    r = math.Exp(r)
  }

  s.SetValue(r); return nil
}

func (obj normalizationClassifier) Dims() (int, int) {
  return obj.n, 1
}

func (obj normalizationClassifier) CloneMatrixBatchClassifier() MatrixBatchClassifier {
  return normalizationClassifier{obj.k, obj.n, obj.logScale}
}

/* -------------------------------------------------------------------------- */

func multi_feature_eval_norm(config ConfigModHmm, state string, trackFiles []string, tracks []Track, filenameResult string, logScale bool) []Track {
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

  n := len(tracks)
  i := MultiFeatureList.Index(strings.ToLower(state))

  if result, err := BatchClassifyMultiTrack(config.SessionConfig, normalizationClassifier{i, n, logScale}, tracks, false); err != nil {
    log.Fatal(err)
  } else {
    if err := ExportTrack(config.SessionConfig, result, filenameResult); err != nil {
      log.Fatal(err)
    }
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_norm_dep(config ConfigModHmm) []string {
  return modhmm_segmentation_dep(config)
}

func modhmm_multi_feature_eval_norm(config ConfigModHmm, state string, tracks []Track, logScale bool) []Track {

  if !MultiFeatureList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_norm_dep(config)...)

  trackFiles := modhmm_multi_feature_eval_norm_dep(config)
  filenameResult := TargetFile{}
  if logScale {
    filenameResult = config.MultiFeatureProbNorm.GetTargetFile(state)
  } else {
    filenameResult = config.MultiFeatureProbNormExp.GetTargetFile(state)
  }

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_multi_feature_eval_all(config, true)

    printStderr(config, 1, "==> Evaluating Normalized Multi-Feature Model (%s) <==\n", strings.ToUpper(state))
    tracks = multi_feature_eval_norm(config, state, trackFiles, tracks, filenameResult.Filename, logScale)
  }
  return tracks
}

func modhmm_multi_feature_eval_norm_loop(config ConfigModHmm, states []string, logScale bool) {
  var tracks []Track
  for _, state := range states {
    tracks = modhmm_multi_feature_eval_norm(config, state, tracks, logScale)
  }
}

func modhmm_multi_feature_eval_norm_all(config ConfigModHmm, logScale bool) {
  modhmm_multi_feature_eval_norm_loop(config, MultiFeatureList, logScale)
}
