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

func single_feature_import_model(config ConfigModHmm, files SingleFeatureFiles, normalize bool) Track {
  if files.Feature == "h3k4me3o1" {
    // check if h3k4me1 or h3k4me3 must be updated first
    files1 := config.SingleFeatureFiles("h3k4me1")
    files2 := config.SingleFeatureFiles("h3k4me3")
    if (FileExists(files1.Model.Filename) && updateRequired(config, files1.Model, files1.DependenciesModel()...)) ||
      ((FileExists(files2.Model.Filename) && updateRequired(config, files2.Model, files2.DependenciesModel()...))) {
      log.Fatalf("ERROR: Please first update single-feature models for `h3k4me1' and `h3k4me3'.\n" +
        "Custom single-feature models are being used. This error occurs because the\n" +
        "time-stamp of coverage files is newer than those of the single-feature model\n" +
        "files. Please make sure that the models are up to date. Use\n" +
        "\t\"Single-Feature Model Static\": true\n" +
        "in the config file prevent this check.")
    }
    config.BinSummaryStatistics = "mean"
    config.BinOverlap = 1
    track1 := single_feature_import_and_normalize(config, files.SrcCoverage[0].Filename, files.SrcCoverageCnts[0].Filename, normalize)
    track2 := single_feature_import_and_normalize(config, files.SrcCoverage[1].Filename, files.SrcCoverageCnts[1].Filename, normalize)
    return single_feature_compute_h3k4me3o1(config, track1, track2)
  } else {
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
    return single_feature_import_and_normalize(config, files.Coverage.Filename, files.CoverageCnts.Filename, normalize)
  }
}

/* -------------------------------------------------------------------------- */

func single_feature_eval_classifier(config ConfigModHmm, files SingleFeatureFiles) {
  mixture := ImportMixtureDistribution(config, files.Model.Filename)
  k, _    := ImportComponents(config, files.Components.Filename, mixture.NComponents())

  scalarClassifier := scalarClassifier.MixturePosterior{mixture, k}
  vectorClassifier := vectorClassifier.ScalarBatchIid{scalarClassifier, 1}

  data := single_feature_import_model(config, files, true)

  result, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier, data); if err != nil {
    log.Fatal(err)
  }
  // rna-low is a special case
  if files.Feature == "rna" {
    single_feature_eval_rna_low(config, result, data)
  }
  if err := (GenericMutableTrack{result}).Map(result, func(seqname string, position int, value float64) float64 {
    return math.Exp(value)
  }); err != nil {
    log.Fatal(err)
  }
  if err := ExportTrack(config.SessionConfig, result, files.Probabilities.Filename); err != nil {
    log.Fatal(err)
  }
}
