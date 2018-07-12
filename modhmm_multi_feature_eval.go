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

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func multi_feature_eval_mixture_weights(config ConfigModHmm) []float64 {
  checkModelFiles(config.SingleFeatureJson)
  checkModelFiles(config.SingleFeatureComp)
  pi := []float64{}
  for _, feature := range []string{"atac", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "rnaLow", "control"} {
    filenameModel := getFieldString(config.SingleFeatureJson, feature)
    filenameComp  := getFieldString(config.SingleFeatureComp, feature)
    p, q := ImportMixtureWeights(config, filenameModel, filenameComp)
    pi = append(pi, p, q)
  }
  return pi
}

/* -------------------------------------------------------------------------- */

func multi_feature_eval(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, tracks []Track, result1, result2 string) []Track {
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
  if err := ExportTrack(config.SessionConfig, result, result1); err != nil {
    log.Fatal(err)
  }
  (GenericMutableTrack{result}).Map(result, func(name string, position int, x float64) float64 {
    return math.Exp(x)
  })
  if err := ExportTrack(config.SessionConfig, result, result2); err != nil {
    log.Fatal(err)
  }
  return tracks
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_dep(config ConfigModHmm) []string {
  return []string{
    config.SingleFeatureFg.Atac,
    config.SingleFeatureBg.Atac,
    config.SingleFeatureFg.H3k27ac,
    config.SingleFeatureBg.H3k27ac,
    config.SingleFeatureFg.H3k27me3,
    config.SingleFeatureBg.H3k27me3,
    config.SingleFeatureFg.H3k9me3,
    config.SingleFeatureBg.H3k9me3,
    config.SingleFeatureFg.H3k4me1,
    config.SingleFeatureBg.H3k4me1,
    config.SingleFeatureFg.H3k4me3,
    config.SingleFeatureBg.H3k4me3,
    config.SingleFeatureFg.H3k4me3o1,
    config.SingleFeatureBg.H3k4me3o1,
    config.SingleFeatureFg.Rna,
    config.SingleFeatureBg.Rna,
    config.SingleFeatureFg.RnaLow,
    config.SingleFeatureBg.RnaLow,
    config.SingleFeatureFg.Control,
    config.SingleFeatureBg.Control }
}

func modhmm_multi_feature_eval(config ConfigModHmm, state string, tracks []Track) []Track {

  localConfig := config
  localConfig.BinSummaryStatistics = "mean"

  dependencies := []string{}
  dependencies  = append(dependencies, modhmm_single_feature_eval_dep(config)...)
  dependencies  = append(dependencies, modhmm_multi_feature_eval_dep(config)...)
  trackFiles   := modhmm_multi_feature_eval_dep(config)

  filenameResult1 := ""
  filenameResult2 := ""

  pi := multi_feature_eval_mixture_weights(config)

  var classifier MatrixBatchClassifier

  switch strings.ToLower(state) {
  case "pa":
    filenameResult1 = config.MultiFeatureClass.PA
    filenameResult2 = config.MultiFeatureClassExp.PA
    classifier = ModelPA{BasicMultiFeatureModel{pi}}
  case "pb":
    filenameResult1 = config.MultiFeatureClass.PB
    filenameResult2 = config.MultiFeatureClassExp.PB
    classifier = ModelPB{BasicMultiFeatureModel{pi}}
  case "ea":
    filenameResult1 = config.MultiFeatureClass.EA
    filenameResult2 = config.MultiFeatureClassExp.EA
    classifier = ModelEA{BasicMultiFeatureModel{pi}}
  case "ep":
    filenameResult1 = config.MultiFeatureClass.EP
    filenameResult2 = config.MultiFeatureClassExp.EP
    classifier = ModelEP{BasicMultiFeatureModel{pi}}
  case "tr":
    filenameResult1 = config.MultiFeatureClass.TR
    filenameResult2 = config.MultiFeatureClassExp.TR
    classifier = ModelTR{BasicMultiFeatureModel{pi}}
  case "tl":
    filenameResult1 = config.MultiFeatureClass.TL
    filenameResult2 = config.MultiFeatureClassExp.TL
    classifier = ModelTL{BasicMultiFeatureModel{pi}}
  case "r1":
    filenameResult1 = config.MultiFeatureClass.R1
    filenameResult2 = config.MultiFeatureClassExp.R1
    classifier = ModelR1{BasicMultiFeatureModel{pi}}
  case "r2":
    filenameResult1 = config.MultiFeatureClass.R2
    filenameResult2 = config.MultiFeatureClassExp.R2
    classifier = ModelR2{BasicMultiFeatureModel{pi}}
  case "ns":
    filenameResult1 = config.MultiFeatureClass.NS
    filenameResult2 = config.MultiFeatureClassExp.NS
    classifier = ModelNS{BasicMultiFeatureModel{pi}}
  case "cl":
    filenameResult1 = config.MultiFeatureClass.CL
    filenameResult2 = config.MultiFeatureClassExp.CL
    classifier = ModelCL{BasicMultiFeatureModel{pi}}
  default:
    log.Fatalf("unknown state: %s", state)
  }
  if updateRequired(config, filenameResult1, dependencies...) ||
    (updateRequired(config, filenameResult2, dependencies...)) {
    modhmm_single_feature_eval_all(config)
    printStderr(config, 1, "==> Computing Multi-Feature Classification (%s) <==\n", strings.ToUpper(state))
    tracks = multi_feature_eval(localConfig, classifier, trackFiles, tracks, filenameResult1, filenameResult2)
  }
  return tracks
}

func modhmm_multi_feature_eval_all(config ConfigModHmm) {
  var tracks []Track
  for _, feature := range []string{"pa", "pb", "ea", "ep", "tr", "tl", "r1", "r2", "ns", "cl"} {
    tracks = modhmm_multi_feature_eval(config, feature, tracks)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-multi-feature", os.Args[0]))
  options.SetParameters("<STATE>\n")

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
