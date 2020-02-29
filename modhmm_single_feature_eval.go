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

import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func single_feature_import_and_normalize(config ConfigModHmm, filenameData, filenameCnts string, normalize bool) MutableTrack {
  if track, err := ImportTrack(config.SessionConfig, filenameData); err != nil {
    log.Fatal(err)
    return nil
  } else {
    if normalize {
      counts := ImportCounts(config, filenameCnts)
      printStderr(config, 1, "Quantile normalizing track to reference distribution... ")
      if err := (GenericMutableTrack{track}).QuantileNormalizeToCounts(counts.X, counts.Y); err != nil {
        printStderr(config, 1, "failed\n")
        log.Fatal(err)
      }
      printStderr(config, 1, "done\n")
    }
    return track
  }
}

func single_feature_compute_h3k4me3o1(config ConfigModHmm, track1, track2 MutableTrack) MutableTrack {
  n1 := int64(0)
  n2 := int64(0)
  if err := (GenericMutableTrack{}).MapList([]Track{track1, track2}, func(seqname string, position int, values ...float64) float64 {
    n1 += int64(values[0])
    n2 += int64(values[1])
    return 0.0
  }); err != nil {
    log.Fatal(err)
  }
  z := float64(n1)/float64(n2)
  if err := (GenericMutableTrack{track1}).MapList([]Track{track1, track2}, func(seqname string, position int, values ...float64) float64 {
    x1 := values[0]
    x2 := values[1]
    // do not add a pseudocount to x2 so that if x1 and x2
    // are both zero, also the result is zero
    // (otherwise strange peaks appear in the distribution)
    return math.Round(z*(x2+0.0)/(x1+1.0)*10)
  }); err != nil {
    log.Fatal(err)
  }
  return track1
}

/* -------------------------------------------------------------------------- */

func single_feature_import(config ConfigModHmm, files SingleFeatureFiles, normalize bool) Track {
  switch strings.ToLower(config.SingleFeatureMethod) {
  case "model"    : return single_feature_import_model    (config, files, normalize)
  case "heuristic": return single_feature_import_heuristic(config, files)
  default:
    log.Fatal("invalid single-feature method")
    panic("internal error")
  }
}

func single_feature_eval(config ConfigModHmm, files SingleFeatureFiles, logScale bool) {
  switch strings.ToLower(config.SingleFeatureMethod) {
  case "model"    : single_feature_eval_classifier(config, files, logScale)
  case "heuristic": single_feature_eval_heuristic (config, files, logScale)
  default:
    log.Fatal("invalid single-feature method")
    panic("internal error")
  }
}

func single_feature_filter_update(config ConfigModHmm, features []string, logScale bool) []string {
  r := []string{}
  for _, feature := range features {
    feature = config.CoerceOpenChromatinAssay(feature)

    files := config.SingleFeatureFiles(feature, logScale)

    dependencies := []string{}
    dependencies  = append(dependencies, files.Dependencies()...)

    switch files.Feature {
    case "h3k4me3o1":
      dependencies  = append(dependencies, modhmm_coverage_dep(config, "h3k4me1", "h3k4me3")...)
    case "rna-low":
      dependencies  = append(dependencies, modhmm_coverage_dep(config, "rna")...)
    default:
      dependencies  = append(dependencies, modhmm_coverage_dep(config, files.Feature)...)
    }
    if updateRequired(config, files.Foreground, dependencies...) ||
      (updateRequired(config, files.Background, dependencies...)) {
      if files.Feature == "rna-low" {
        r = append(r, "rna")
      } else {
        r = append(r, files.Feature)
      }
    }
  }
  return uniqueStrings(r)
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_dep(config ConfigModHmm) []string {
  r := []string{}
  r  = append(r, config.Coverage          .GetFilenames()...)
  r  = append(r, config.SingleFeatureModel.GetFilenames()...)
  r  = append(r, config.SingleFeatureComp .GetFilenames()...)
  r  = append(r, config.CoverageCnts      .GetFilenames()...)
  return r
}

func modhmm_single_feature_eval(config ConfigModHmm, feature string, logScale bool) {

  files := config.SingleFeatureFiles(feature, logScale)

  if updateRequired(config, files.Foreground, files.Dependencies()...) ||
    (updateRequired(config, files.Background, files.Dependencies()...)) {

    printStderr(config, 1, "==> Evaluating Single-Feature Model (%s) <==\n", feature)
    single_feature_eval(config, files, logScale)
  }
}

func modhmm_single_feature_eval_loop(config ConfigModHmm, features []string, logScale bool) {
  // reduce list of features to those that require an update
  features = single_feature_filter_update(config, features, logScale)
  // compute coverages here to make use of multi-threading
  modhmm_coverage_loop(config, InsensitiveStringList(features).Intersection(CoverageList))
  // eval single features
  for _, feature := range features {
    modhmm_single_feature_eval(config, feature, logScale)
  }
}

func modhmm_single_feature_eval_all(config ConfigModHmm, logScale bool) {
  modhmm_single_feature_eval_loop(config, SingleFeatureList, logScale)
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optStdScale := options.BoolLong("std-scale",  0 ,  "single-feature output on standard scale")
  optHelp     := options.BoolLong("help",      'h',  "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_single_feature_eval_all(config, !*optStdScale)
  } else {
    modhmm_single_feature_eval_loop(config, options.Args(), !*optStdScale)
  }
}
