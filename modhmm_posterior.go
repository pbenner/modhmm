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
import   "github.com/pbenner/autodiff/statistics/matrixClassifier"

import . "github.com/pbenner/gonetics"
import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func getStateIndices(modhmm ModHmm, state string) []int {
  stateMap := modhmm.Hmm.StateMap
  iState   := multiFeatureList.Index(strings.ToLower(state))
  result   := []int{}
  for i := 0; i < len(stateMap); i++ {
    if stateMap[i] == iState {
      result = append(result, i)
    }
  }
  return result
}

/* -------------------------------------------------------------------------- */

func posterior(config ConfigModHmm, state string, trackFiles []string, tracks []Track, filenameResult string) []Track {
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
  modhmm := ModHmm{}
  printStderr(config, 1, "Importing model from `%s'... ", config.Model)
  if err := ImportDistribution(config.Model, &modhmm, BareRealType); err != nil {
    log.Fatal(err)
    printStderr(config, 1, "failed\n")
  }
  printStderr(config, 1, "done\n")

  states := getStateIndices(modhmm, state)
  printStderr(config, 1, "State %s with maps to state indices %v\n", state, states)

  result, err := ClassifyMultiTrack(config.SessionConfig, matrixClassifier.HmmPosterior{&modhmm.Hmm, states, false}, tracks, true); if err != nil {
    panic(err)
  }
  err = ExportTrack(config.SessionConfig, result, filenameResult); if err != nil {
    panic(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_posterior_tracks(config ConfigModHmm) []string {
  files := make([]string, len(multiFeatureList))
  for i, state := range multiFeatureList {
    files[i] = getFieldAsString(config.MultiFeatureProb, strings.ToUpper(state))
  }
  return files
}

func modhmm_posterior(config ConfigModHmm, state string, tracks []Track) []Track {

  if !multiFeatureList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_segmentation_dep(config)...)
  dependencies  = append(dependencies, config.Model)

  trackFiles := modhmm_posterior_tracks(config)
  filenameResult := getFieldAsString(config.Posterior, strings.ToUpper(state))

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_multi_feature_eval_all(config)
    modhmm_segmentation(config, "default")

    printStderr(config, 1, "==> Evaluating Posterior Marginals (%s) <==\n", strings.ToUpper(state))
    tracks = posterior(config, state, trackFiles, tracks, filenameResult)
  }
  return tracks
}

func modhmm_posterior_all(config ConfigModHmm) {
  var tracks []Track
  for _, state := range multiFeatureList {
    tracks = modhmm_posterior(config, state, tracks)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_posterior_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s posterior-marginals", os.Args[0]))
  options.SetParameters("<STATE>\n")

  optHelp  := options.BoolLong("help", 'h', "print help")

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
    modhmm_posterior_all(config)
  } else {
    modhmm_posterior(config, options.Args()[0], nil)
  }
}
