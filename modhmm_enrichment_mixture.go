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

//import   "fmt"
import   "log"
import   "math"

import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/track"

import   "github.com/pbenner/autodiff/statistics/scalarClassifier"
import   "github.com/pbenner/autodiff/statistics/vectorClassifier"

import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

/* -------------------------------------------------------------------------- */

func enrichment_import_model(config ConfigModHmm, files EnrichmentFiles, normalize bool) Track {
  // check if single feature model must be updated
  if normalize && FileExists(files.Model.Filename) && updateRequired(config, files.Model, files.DependenciesModel()...) {
    log.Fatalf("ERROR: Please first update single-feature model for `%s'.\n" +
      "Custom single-feature models are being used. This error occurs because the\n" +
      "time-stamp of coverage files is newer than those of the single-feature model\n" +
      "files. Please make sure that the models are up to date. Use\n" +
      "\t\"Single-Feature Model Static\": true\n" +
      "in the config file prevent this check.", files.Feature)
  }
  config.BinSummaryStatistics = "discrete mean"
  return enrichment_import_and_normalize(config, files.Coverage.Filename, files.CoverageCnts.Filename, normalize)
}

/* -------------------------------------------------------------------------- */

func enrichment_eval_classifier(config ConfigModHmm, files EnrichmentFiles) {
  mixture := ImportMixtureDistribution(config, files.Model.Filename)
  k, _    := ImportComponents(config, files.Components.Filename, mixture.NComponents())

  scalarClassifier := scalarClassifier.MixturePosterior{mixture, k}
  vectorClassifier := vectorClassifier.ScalarBatchIid{scalarClassifier, 1}

  data := enrichment_import_model(config, files, true)

  result, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier, data); if err != nil {
    log.Fatal(err)
  }
  if files.Feature == "rna" {
    counts := compute_counts(config, data)
    q := config.EnrichmentParameters.GetParameters(files.Feature)[0]
    t := counts.Quantile(q)
    enrichment_eval_rna(config, result, data, t)
  } else {
    if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
      return math.Exp(value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result, files.Probabilities.Filename); err != nil {
    log.Fatal(err)
  }
}
