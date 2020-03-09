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

import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/modhmm/config"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func modhmm_call_enrichment_peaks(config ConfigModHmm, feature string, threshold float64) {
  printStderr(config, 1, "==> Calling Single-Feature Peaks (%s) <==\n", feature)
  filenameIn  := config.EnrichmentProb.GetTargetFile(feature).Filename
  filenameOut := config.EnrichmentPeak.GetTargetFile(feature)

  if !updateRequired(config, filenameOut, filenameIn) {
    return
  }
  modhmm_enrichment_eval(config, feature)

  if track, err := ImportTrack(config.SessionConfig, filenameIn); err != nil {
    log.Fatal(err)
  } else {
    if peaks, err := getPeaks(track, threshold); err != nil {
      log.Fatal(err)
    } else {
      printStderr(config, 1, "Writing table `%s'... ", filenameOut.Filename)
      if err := peaks.ExportTable(filenameOut.Filename, true, false, false, OptionPrintScientific{true}); err != nil {
        printStderr(config, 1, "failed\n")
        log.Fatal(err)
      } else {
        printStderr(config, 1, "done\n")
      }
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_call_enrichment_peaks_loop(config ConfigModHmm, features []string, threshold float64) {
  for _, feature := range features {
    feature = config.CoerceOpenChromatinAssay(feature)
    modhmm_call_enrichment_peaks(config, feature, threshold)
  }
}

func modhmm_call_enrichment_peaks_all(config ConfigModHmm, threshold float64) {
  modhmm_call_enrichment_peaks_loop(config, EnrichmentList, threshold)
}

func modhmm_call_enrichment_peaks_main(config ConfigModHmm, args []string) {

  var threshold float64

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s call-single-feature-peaks", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

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
    modhmm_call_enrichment_peaks_all(config, threshold)
  } else {
    modhmm_call_enrichment_peaks_loop(config, options.Args(), threshold)
  }
}
