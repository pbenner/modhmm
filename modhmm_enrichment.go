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

import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func enrichment_import_and_normalize(config ConfigModHmm, filenameData, filenameCnts string, normalize bool) MutableTrack {
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

/* -------------------------------------------------------------------------- */

func enrichment_eval_rna(config ConfigModHmm, result MutableTrack, data Track, t float64) {
  if err := (GenericMutableTrack{result}).Map(data, func(seqname string, position int, value float64) float64 {
    if value > t {
      return 1.0 - 1e-8
    } else {
      return 0.01
    }
  }); err != nil {
    log.Fatal(err)
  }
}

/* -------------------------------------------------------------------------- */

func enrichment_import(config ConfigModHmm, files EnrichmentFiles, normalize bool) Track {
  switch strings.ToLower(config.EnrichmentMethod) {
  case "model"    : return enrichment_import_model    (config, files, normalize)
  case "heuristic": return enrichment_import_heuristic(config, files)
  default:
    log.Fatal("invalid single-feature method")
    panic("internal error")
  }
}

func enrichment_eval(config ConfigModHmm, files EnrichmentFiles) {
  switch strings.ToLower(config.EnrichmentMethod) {
  case "model"    : enrichment_eval_classifier(config, files)
  case "heuristic": enrichment_eval_heuristic (config, files)
  default:
    log.Fatal("invalid single-feature method")
    panic("internal error")
  }
}

func enrichment_filter_update(config ConfigModHmm, features []string) []string {
  r := []string{}
  for _, feature := range features {
    feature = config.CoerceOpenChromatinAssay(feature)

    files := config.EnrichmentFiles(feature)

    dependencies := []string{}
    dependencies  = append(dependencies, files.Dependencies()...)
    dependencies  = append(dependencies, modhmm_coverage_dep(config, files.Feature)...)
    if updateRequired(config, files.Probabilities, dependencies...) {
      r = append(r, files.Feature)
    }
  }
  return uniqueStrings(r)
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_eval_dep(config ConfigModHmm) []string {
  r := []string{}
  r  = append(r, config.Coverage          .GetFilenames()...)
  r  = append(r, config.EnrichmentModel.GetFilenames()...)
  r  = append(r, config.EnrichmentComp .GetFilenames()...)
  r  = append(r, config.CoverageCnts      .GetFilenames()...)
  return r
}

func modhmm_enrichment_eval(config ConfigModHmm, feature string) {

  files := config.EnrichmentFiles(feature)

  if updateRequired(config, files.Probabilities, files.Dependencies()...) {

    if EnrichmentIsOptional(files.Feature) && !FileExists(files.Coverage.Filename) {
      return
    }
    printStderr(config, 1, "==> Computing Enrichment Probabilities (%s) <==\n", feature)
    enrichment_eval(config, files)
  }
}

func modhmm_enrichment_eval_loop(config ConfigModHmm, features []string) {
  // reduce list of features to those that require an update
  features = enrichment_filter_update(config, features)
  // compute coverages here to make use of multi-threading
  modhmm_coverage_loop(config, InsensitiveStringList(features).Intersection(CoverageList))
  // eval single features
  for _, feature := range features {
    modhmm_enrichment_eval(config, feature)
  }
}

func modhmm_enrichment_eval_all(config ConfigModHmm) {
  modhmm_enrichment_eval_loop(config, EnrichmentList)
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_eval_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s eval-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optHelp := options.BoolLong("help", 'h', "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_enrichment_eval_all(config)
  } else {
    modhmm_enrichment_eval_loop(config, options.Args())
  }
}
