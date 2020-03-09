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

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func get_chromatin_state_model(config ConfigModHmm, state string) MatrixBatchClassifier {
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

func chromatin_state_eval(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, tracks []Track, filenameResult string) []Track {
  if len(tracks) != len(trackFiles) {
    tracks  = make([]Track, len(trackFiles))
    genome := Genome{}
    empty  := []int{}
    for i, filename := range trackFiles {
      if !FileExists(filename) && EnrichmentIsOptional(EnrichmentList[i]) {
        empty = append(empty, i)
        continue
      }
      if t, err := ImportTrack(config.SessionConfig, filename); err != nil {
        log.Fatal(err)
      } else {
        tracks[i] = t
        genome    = t.GetGenome()
      }
    }
    if len(empty) > 0 {
      t := AllocSimpleTrack("classification", genome, config.BinSize)
      for _, i := range empty {
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

func modhmm_chromatin_state_eval_dep(config ConfigModHmm) []string {
  files := []string{}
  for _, feature := range EnrichmentList {
    files = append(files, config.EnrichmentProb.GetTargetFile(feature).Filename)
  }
  return files
}

func modhmm_chromatin_state_eval(config ConfigModHmm, state string, tracks []Track) []Track {

  if !ChromatinStateList.Contains(strings.ToLower(state)) {
    log.Fatalf("unknown state: %s", state)
  }

  localConfig := config
  localConfig.BinSummaryStatistics = "mean"

  dependencies   := []string{}
  dependencies    = append(dependencies, modhmm_chromatin_state_eval_dep(config)...)
  dependencies    = append(dependencies, modhmm_enrichment_eval_dep(config)...)
  dependencies    = append(dependencies, modhmm_coverage_dep(config)...)

  trackFiles     := modhmm_chromatin_state_eval_dep(config)
  filenameResult := config.ChromatinStateProb.GetTargetFile(state)

  if updateRequired(config, filenameResult, dependencies...) {
    modhmm_enrichment_eval_all(config)
    printStderr(config, 1, "==> Evaluating Multi-Feature Model (%s) <==\n", strings.ToUpper(state))
    classifier := get_chromatin_state_model(config, state)
    tracks = chromatin_state_eval(localConfig, classifier, trackFiles, tracks, filenameResult.Filename)
  }
  return tracks
}

func modhmm_chromatin_state_eval_loop(config ConfigModHmm, states []string) {
  var tracks []Track
  for _, feature := range states {
    tracks = modhmm_chromatin_state_eval(config, feature, tracks)
  }
}

func modhmm_chromatin_state_eval_all(config ConfigModHmm) {
  modhmm_chromatin_state_eval_loop(config, ChromatinStateList)
}

/* -------------------------------------------------------------------------- */

func modhmm_chromatin_state_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-multi-feature", os.Args[0]))
  options.SetParameters("[STATE]...\n")

  optHelp := options.BoolLong("help", 'h', "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_chromatin_state_eval_all(config)
  } else {
    modhmm_chromatin_state_eval_loop(config, options.Args())
  }
}
