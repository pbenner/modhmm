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

/* -------------------------------------------------------------------------- */

func single_feature_eval_classifier(config ConfigModHmm, files SingleFeatureFiles, logScale bool) {
  mixture := ImportMixtureDistribution(config, files.Model.Filename)
  k, r    := ImportComponents(config, files.Components.Filename, mixture.NComponents())

  scalarClassifier1 := scalarClassifier.MixturePosterior{mixture, k}
  scalarClassifier2 := scalarClassifier.MixturePosterior{mixture, r}
  vectorClassifier1 := vectorClassifier.ScalarBatchIid{scalarClassifier1, 1}
  vectorClassifier2 := vectorClassifier.ScalarBatchIid{scalarClassifier2, 1}

  data := single_feature_import(config, files, true)

  // foreground
  result1, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier1, data); if err != nil {
    log.Fatal(err)
  }
  if !logScale {
    if err := (GenericMutableTrack{result1}).Map(result1, func(seqname string, position int, value float64) float64 {
      return math.Exp(value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result1, files.Foreground.Filename); err != nil {
    log.Fatal(err)
  }
  // background
  result2, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier2, data); if err != nil {
    log.Fatal(err)
  }
  if !logScale {
    if err := (GenericMutableTrack{result2}).Map(result2, func(seqname string, position int, value float64) float64 {
      return math.Exp(value)
    }); err != nil {
      log.Fatal(err)
    }
  }
  if err := ExportTrack(config.SessionConfig, result2, files.Background.Filename); err != nil {
    log.Fatal(err)
  }
}
