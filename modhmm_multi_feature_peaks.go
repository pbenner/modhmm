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
import   "strconv"
import   "strings"

import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/gonetics"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func modhmm_call_multi_feature_peaks(config ConfigModHmm, state string, threshold float64) {
  printStderr(config, 1, "==> Calling Multi-Feature Peaks (%s) <==\n", strings.ToUpper(state))
  filenameIn  := getFieldAsString(config.MultiFeatureProb, strings.ToUpper(state))
  filenameOut := getFieldAsString(config.MultiFeaturePeak, strings.ToUpper(state))

  if track, err := ImportTrack(config.SessionConfig, filenameIn); err != nil {
    log.Fatal(err)
  } else {
    if peaks, err := getPeaks(track, threshold); err != nil {
      log.Fatal(err)
    } else {
      printStderr(config, 1, "Writing table `%s'... ", filenameOut)
      if err := peaks.ExportTable(filenameOut, true, false, false, OptionPrintScientific{true}); err != nil {
        printStderr(config, 1, "failed\n")
        log.Fatal(err)
      } else {
        printStderr(config, 1, "done\n")
      }
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_call_multi_feature_peaks_loop(config ConfigModHmm, states []string, threshold float64) {
  for _, state := range states {
    modhmm_call_multi_feature_peaks(config, state, threshold)
  }
}

func modhmm_call_multi_feature_peaks_all(config ConfigModHmm, threshold float64) {
  modhmm_call_multi_feature_peaks_loop(config, multiFeatureList, threshold)
}

func modhmm_call_multi_feature_peaks_main(config ConfigModHmm, args []string) {

  var threshold float64

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s call-multi-feature-peaks", os.Args[0]))
  options.SetParameters("[STATE]...\n")

  optThreshold := options.StringLong("threshold",  0 ,  "0.9", "threshold value [default 0.9]")
  optHelp      := options.BoolLong  ("help",      'h',         "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if t, err := strconv.ParseFloat(*optThreshold, 64); err != nil {
    log.Fatal(err)
  } else {
    threshold = t
  }

  if len(options.Args()) == 0 {
    modhmm_call_multi_feature_peaks_all(config, threshold)
  } else {
    modhmm_call_multi_feature_peaks_loop(config, options.Args(), threshold)
  }
}
