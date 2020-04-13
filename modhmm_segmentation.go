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
import   "path"
import   "os"
import   "math"

import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/estimation"
import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/modhmm/config"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/matrixClassifier"
import   "github.com/pbenner/autodiff/statistics/matrixDistribution"
import   "github.com/pbenner/autodiff/statistics/matrixEstimator"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func ImportHMM(config ConfigModHmm) ModHmm {
  modhmm   := ModHmm{}
  filename := config.Model.Filename
  printStderr(config, 2, "Importing HMM model from `%s'... ", config.Model.Filename)
  if err := ImportDistribution(filename, &modhmm, BareRealType); err != nil {
    printStderr(config, 2, "failed\n")
    // remove directory from filename
    _, filename = path.Split(filename)
    filename = fmt.Sprintf("%s.json", config.ModelFallbackPath())
    printStderr(config, 2, "Importing HMM fallback model (%s)... ", config.ModelFallback)
    if err := ImportDefaultDistribution(config, filename, &modhmm, BareRealType); err != nil {
      printStderr(config, 2, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 2, "done\n")
  } else {
    printStderr(config, 2, "done\n")
  }
  return modhmm
}

/* -------------------------------------------------------------------------- */

func import_chromatin_state_tracks(config ConfigModHmm, tracks []Track, trackFiles []string) []Track {
  if len(tracks) == 0 {
    tracks = make([]Track, len(trackFiles))
  }
  // import track files (do not use ImportAndEstimateOnMultiTrack, which uses lazy imports)
  for i := 0; i < len(trackFiles); i++ {
    if tracks[i] == nil {
      track, err := ImportTrack(config.SessionConfig, trackFiles[i]); if err != nil {
        log.Fatal(err)
      }
      tracks[i] = track
    }
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

type ChromatinStateFilterZeros struct {
}

// filter strange probability assignments at chromosome boundaries
func (ChromatinStateFilterZeros) Eval(x Matrix) Matrix {
  n, m := x.Dims()
  for i := 0; i < n; i++ {
    allZero := true
    for j := 0; j < m; j++ {
      if math.IsNaN(x.ValueAt(i, j)) {
        panic("interal error")
      }
      if x.ValueAt(i, j) != 0.0 {
        allZero = false; break
      }
    }
    if allZero {
      x.At(i, iNS).SetValue(1.0)
    }
  }
  return x
}

/* -------------------------------------------------------------------------- */

func estimate(config ConfigModHmm, tracks []Track, trackFiles []string, model string) {
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
  tracks = import_chromatin_state_tracks(config, nil, trackFiles)

  if err := EstimateOnMultiTrack(config.SessionConfig, estimator, tracks, true, ChromatinStateFilterZeros{}); err != nil {
    log.Fatalf("ERROR: %s", err)
  }
  modhmm := ModHmm{}
  if d, err := estimator.GetEstimate(); err != nil {
    log.Fatalf("ERROR: %s", err)
  } else {
    modhmm.Hmm = *d.(*matrixDistribution.Hmm)
  }
  modhmm.StateNames = stateNames

  printStderr(config, 1, "Exporting model to `%s'... ", config.Model.Filename)
  if err := ExportDistribution(config.Model.Filename, &modhmm); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("ERROR: %s", err)
  }
  printStderr(config, 1, "done\n")
}

/* -------------------------------------------------------------------------- */

func segment(config ConfigModHmm, tracks []Track, trackFiles []string) {
  modhmm := ImportHMM(config)
  tracks  = import_chromatin_state_tracks(config, tracks, trackFiles)

  // compute segmentation
  if result, err := ClassifyMultiTrack(config.SessionConfig, matrixClassifier.HmmClassifier{&modhmm.Hmm}, tracks, true, ChromatinStateFilterZeros{}); err != nil {
    log.Fatal(err)
  } else {
    var name, desc string
    if config.Description == "" {
      name = "ModHMM"
      desc = "Segmentation ModHMM"
    } else {
      name = fmt.Sprintf("ModHMM [%s]", config.Description)
      desc = fmt.Sprintf("Segmentation ModHMM:%s [%s]", Version, config.Description)
    }
    tracksEquivalent := make([]Track, modhmm.NStates())
    for i, state := range ChromatinStateList {
      for _, j := range getStateIndices(modhmm, state) {
        tracksEquivalent[j] = tracks[i]
      }
    }
    printStderr(config, 1, "Writing genome segmentation to `%s'... ", config.Segmentation.Filename)
    if err := ExportTrackSegmentation(config.SessionConfig, result, config.Segmentation.Filename, name, desc, true, modhmm.StateNames, getRGBMap(), tracksEquivalent); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_segmentation_dep(config ConfigModHmm) []string {
  files := make([]string, len(ChromatinStateList))
  for i, state := range ChromatinStateList {
    files[i] = config.ChromatinStateProb.GetTargetFile(state).Filename
  }
  return files
}

func modhmm_segmentation(config ConfigModHmm, model string) {

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_segmentation_dep(config)...)
  dependencies  = append(dependencies, modhmm_chromatin_state_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_enrichment_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_coverage_dep(config)...)

  trackFiles := modhmm_segmentation_dep(config)
  tracks     := make([]Track, len(trackFiles))

  filenameModel        := config.Model
  filenameSegmentation := config.Segmentation

  if config.ModelEstimate && updateRequired(config, filenameModel, dependencies...) {
    modhmm_chromatin_state_eval_all(config)

    printStderr(config, 1, "==> Estimating ModHmm transition parameters <==\n")
    estimate(config, tracks, trackFiles, model)
  }
  if config.ModelEstimate {
    dependencies = append(dependencies, filenameModel.Filename)
  }
  if updateRequired(config, filenameSegmentation, dependencies...) {
    modhmm_chromatin_state_eval_all(config)

    printStderr(config, 1, "==> Computing Segmentation <==\n")
    segment(config, tracks, trackFiles)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_segmentation_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s segmentation", os.Args[0]))

  optHelp  := options.   BoolLong("help",  'h',            "print help")
  optModel := options. StringLong("model",  0 , "default", "default, dense")

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
