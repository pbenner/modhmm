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

import   "github.com/pbenner/autodiff/statistics/matrixClassifier"

import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func getStateIndices(modhmm ModHmm, state string) []int {
  stateMap := modhmm.Hmm.StateMap
  iState   := ChromatinStateList.Index(strings.ToLower(state))
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
  modhmm := ImportHMM(config)

  states := getStateIndices(modhmm, state)
  printStderr(config, 2, "State %s maps to state indices %v\n", strings.ToUpper(state), states)
  tracks  = import_chromatin_state_tracks(config, tracks, trackFiles)

  result, err := ClassifyMultiTrack(config.SessionConfig, matrixClassifier.HmmPosterior{&modhmm.Hmm, states, false}, tracks, true, ChromatinStateFilterZeros{}); if err != nil {
    panic(err)
  }
  err = ExportTrack(config.SessionConfig, result, filenameResult); if err != nil {
    panic(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_posterior_tracks(config ConfigModHmm) []string {
  files := make([]string, len(ChromatinStateList))
  for i, state := range ChromatinStateList {
    files[i] = config.ChromatinStateProb.GetTargetFile(state).Filename
  }
  return files
}

func modhmm_posterior(config ConfigModHmm, state string, tracks []Track) []Track {

  if !ChromatinStateList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  dependencies := []string{}
  if config.ModelEstimate {
    dependencies  = append(dependencies, config.Model.Filename)
  }
  dependencies  = append(dependencies, modhmm_segmentation_dep(config)...)
  dependencies  = append(dependencies, modhmm_chromatin_state_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_enrichment_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_coverage_dep(config)...)

  trackFiles     := modhmm_posterior_tracks(config)
  filenameResult := config.PosteriorProb.GetTargetFile(state)

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_chromatin_state_eval_all(config)
    modhmm_segmentation(config, "default")

    printStderr(config, 1, "==> Evaluating Posterior Marginals (%s) <==\n", strings.ToUpper(state))
    tracks = posterior(config, state, trackFiles, tracks, filenameResult.Filename)
  }
  return tracks
}

func modhmm_posterior_loop(config ConfigModHmm, states []string) {
  var tracks []Track
  for _, state := range states {
    tracks = modhmm_posterior(config, state, tracks)
  }
}

func modhmm_posterior_all(config ConfigModHmm) {
  modhmm_posterior_loop(config, ChromatinStateList)
}

/* -------------------------------------------------------------------------- */

func modhmm_posterior_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-posterior-marginals", os.Args[0]))
  options.SetParameters("[STATE]...\n")

  optHelp := options.BoolLong("help",      'h', "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_posterior_all(config)
  } else {
    modhmm_posterior_loop(config, options.Args())
  }
}
