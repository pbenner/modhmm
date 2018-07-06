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

func multi_feature_classify(config ConfigModHmm, classifier MatrixBatchClassifier, trackFiles []string, result1, result2 string) {
  result, err := ImportAndBatchClassifyMultiTrack(config.SessionConfig, classifier, trackFiles, false); if err != nil {
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
}

/* -------------------------------------------------------------------------- */

func modhmm_multi_feature_classify(config ConfigModHmm, state string) {

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
    config.SingleFeatureFg.control,
    config.SingleFeatureBg.control }

  filenameResult1 := ""
  filenameResult2 := ""

  var classifier MatrixBatchClassifier

  switch strings.ToLower(state) {
  case "pa":
    filenameResult1 = config.MultiFeatureClass.PA
    filenameResult1 = config.MultiFeatureClassExp.PA
    classifier = ClassifierPA{}
  case "pb":
    filenameResult1 = config.MultiFeatureClass.PB
    filenameResult1 = config.MultiFeatureClassExp.PB
    classifier = ClassifierPB{}
  default:
    log.Fatal("unknown state: %s", state)
  }

  multi_feature_classify(localConfig, classifier, trackFiles, filenameResult1, filenameResult2)
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
  if len(options.Args()) != 1 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }

  modhmm_single_feature_classify(config, options.Args()[0])
}
