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

func multi_feature_classify(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, tracks []Track, result1, result2 string) []Track {
  if len(tracks) != len(trackFiles) {
    tracks := make([]Track, len(trackFiles))
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

func modhmm_multi_feature_classify(config ConfigModHmm, state string, tracks []Track) []Track {

  localConfig := config
  localConfig.BinSummaryStatistics = "mean"

  trackFiles := []string{
    config.SingleFeatureFg.atac,
    config.SingleFeatureBg.atac,
    config.SingleFeatureFg.h3k27ac,
    config.SingleFeatureBg.h3k27ac,
    config.SingleFeatureFg.h3k27me3,
    config.SingleFeatureBg.h3k27me3,
    config.SingleFeatureFg.h3k9me3,
    config.SingleFeatureBg.h3k9me3,
    config.SingleFeatureFg.h3k4me1,
    config.SingleFeatureBg.h3k4me1,
    config.SingleFeatureFg.h3k4me3,
    config.SingleFeatureBg.h3k4me3,
    config.SingleFeatureFg.h3k4me3o1,
    config.SingleFeatureBg.h3k4me3o1,
    config.SingleFeatureFg.rna,
    config.SingleFeatureBg.rna,
    config.SingleFeatureFg.rnaLow,
    config.SingleFeatureBg.rnaLow,
    config.SingleFeatureFg.control,
    config.SingleFeatureBg.control }

  filenameResult1 := ""
  filenameResult2 := ""

  var classifier MatrixBatchClassifier

  switch strings.ToLower(state) {
  case "pa":
    filenameResult1 = config.MultiFeatureClass.PA
    filenameResult2 = config.MultiFeatureClassExp.PA
    classifier = ClassifierPA{}
  case "pb":
    filenameResult1 = config.MultiFeatureClass.PB
    filenameResult2 = config.MultiFeatureClassExp.PB
    classifier = ClassifierPB{}
  case "ea":
    filenameResult1 = config.MultiFeatureClass.EA
    filenameResult2 = config.MultiFeatureClassExp.EA
    classifier = ClassifierEA{}
  case "ep":
    filenameResult1 = config.MultiFeatureClass.EP
    filenameResult2 = config.MultiFeatureClassExp.EP
    classifier = ClassifierEP{}
  case "tr":
    filenameResult1 = config.MultiFeatureClass.TR
    filenameResult2 = config.MultiFeatureClassExp.TR
    classifier = ClassifierTR{}
  case "tl":
    filenameResult1 = config.MultiFeatureClass.TL
    filenameResult2 = config.MultiFeatureClassExp.TL
    classifier = ClassifierTL{}
  case "r1":
    filenameResult1 = config.MultiFeatureClass.R1
    filenameResult2 = config.MultiFeatureClassExp.R1
    classifier = ClassifierR1{}
  case "r2":
    filenameResult1 = config.MultiFeatureClass.R2
    filenameResult2 = config.MultiFeatureClassExp.R2
    classifier = ClassifierR2{}
  case "ns":
    filenameResult1 = config.MultiFeatureClass.NS
    filenameResult2 = config.MultiFeatureClassExp.NS
    classifier = ClassifierNS{}
  case "cl":
    filenameResult1 = config.MultiFeatureClass.CL
    filenameResult2 = config.MultiFeatureClassExp.CL
    classifier = ClassifierCL{}
  default:
    log.Fatalf("unknown state: %s", state)
  }
  if updateRequired(config, filenameResult1, trackFiles...) ||
    (updateRequired(config, filenameResult2, trackFiles...)) {
    modhmm_single_feature_classify_all(config)
    tracks = multi_feature_classify(localConfig, classifier, trackFiles, tracks, filenameResult1, filenameResult2)
  }
  return tracks
}

func modhmm_multi_feature_classify_all(config ConfigModHmm) {
  var tracks []Track
  for _, feature := range []string{"pa", "pb", "ea", "ep", "tr", "tl", "r1", "r2", "ns", "cl"} {
    tracks = modhmm_multi_feature_classify(config, feature, tracks)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_classify_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s classify-multi-feature-mixture", os.Args[0]))
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
    modhmm_multi_feature_classify_all(config)
  } else {
    modhmm_multi_feature_classify(config, options.Args()[0], nil)
  }
}
