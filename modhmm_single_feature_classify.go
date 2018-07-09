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
import   "io"
import   "os"
import   "strings"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/ngstat/classification"
import . "github.com/pbenner/ngstat/track"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarClassifier"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"
import   "github.com/pbenner/autodiff/statistics/vectorClassifier"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func invertComponents(k []int, n int) []int {
  m := make(map[int]bool)
  r := []int{}
  for _, j := range k {
    m[j] = true
  }
  for j := 0; j < n; j++ {
    if _, ok := m[j]; !ok {
      r = append(r, j)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

type Components []int

func (obj *Components) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, obj)
}

func (obj *Components) Export(writer io.Writer) error {
  return JsonExport(writer, obj)
}

func ImportComponents(config ConfigModHmm, filename string) []int {
  var k Components
  printStderr(config, 1, "Importing foreground components from `%s'... ", filename)
  if err := ImportFile(&k, filename); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("could not import components from `%s': %v", filename, err)
  }
  printStderr(config, 1, "done\n")
  return []int(k)
}

func ExportComponents(config ConfigModHmm, filename string, k []int) {
  printStderr(config, 1, "Exporting foreground components to `%s'... ", filename)
  if err := ExportFile((*Components)(&k), filename); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("could not export components to `%s': %v", filename, err)
  }
  printStderr(config, 1, "done\n")
}

/* -------------------------------------------------------------------------- */

func single_feature_classify(config ConfigModHmm, filenameModel, filenameComp, filenameData, filenameResult1, filenameResult2 string) {
  mixture := &scalarDistribution.Mixture{}

  if err := ImportDistribution(filenameModel, mixture, BareRealType); err != nil {
    panic(err)
  }

  k := ImportComponents(config, filenameComp)
  r := invertComponents(k, mixture.NComponents())

  scalarClassifier1 := scalarClassifier.MixturePosterior{mixture, k}
  vectorClassifier1 := vectorClassifier.ScalarBatchIid{scalarClassifier1, 1}

  scalarClassifier2 := scalarClassifier.MixturePosterior{mixture, r}
  vectorClassifier2 := vectorClassifier.ScalarBatchIid{scalarClassifier2, 1}

  if data, err := ImportTrack(config.SessionConfig, filenameData); err != nil {
    log.Fatal(err)
  } else {
    result1, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier1, data); if err != nil {
      log.Fatal(err)
    }
    if err := ExportTrack(config.SessionConfig, result1, filenameResult1); err != nil {
      log.Fatal(err)
    }

    result2, err := BatchClassifySingleTrack(config.SessionConfig, vectorClassifier2, data); if err != nil {
      log.Fatal(err)
    }
    if err := ExportTrack(config.SessionConfig, result2, filenameResult2); err != nil {
      log.Fatal(err)
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_classify(config ConfigModHmm, feature string) {

  localConfig := config
  localConfig.BinSummaryStatistics = "discrete mean"
  filenameModel   := ""
  filenameComp    := ""
  filenameData    := ""
  filenameResult1 := ""
  filenameResult2 := ""

  switch strings.ToLower(feature) {
  case "atac":
    filenameData    = config.SingleFeatureData.atac
    filenameModel   = config.SingleFeatureJson.atac
    filenameComp    = config.SingleFeatureComp.atac
    filenameResult1 = config.SingleFeatureFg.atac
    filenameResult2 = config.SingleFeatureBg.atac
  case "h3k27ac":
    filenameData    = config.SingleFeatureData.h3k27ac
    filenameModel   = config.SingleFeatureJson.h3k27ac
    filenameComp    = config.SingleFeatureComp.h3k27ac
    filenameResult1 = config.SingleFeatureFg.h3k27ac
    filenameResult2 = config.SingleFeatureBg.h3k27ac
  case "h3k27me3":
    filenameData    = config.SingleFeatureData.h3k27me3
    filenameModel   = config.SingleFeatureJson.h3k27me3
    filenameComp    = config.SingleFeatureComp.h3k27me3
    filenameResult1 = config.SingleFeatureFg.h3k27me3
    filenameResult2 = config.SingleFeatureBg.h3k27me3
  case "h3k9me3":
    filenameData    = config.SingleFeatureData.h3k9me3
    filenameModel   = config.SingleFeatureJson.h3k9me3
    filenameComp    = config.SingleFeatureComp.h3k9me3
    filenameResult1 = config.SingleFeatureFg.h3k9me3
    filenameResult2 = config.SingleFeatureBg.h3k9me3
  case "h3k4me1":
    filenameData    = config.SingleFeatureData.h3k4me1
    filenameModel   = config.SingleFeatureJson.h3k4me1
    filenameComp    = config.SingleFeatureComp.h3k4me1
    filenameResult1 = config.SingleFeatureFg.h3k4me1
    filenameResult2 = config.SingleFeatureBg.h3k4me1
  case "h3k4me3":
    filenameData    = config.SingleFeatureData.h3k4me3
    filenameModel   = config.SingleFeatureJson.h3k4me3
    filenameComp    = config.SingleFeatureComp.h3k4me3
    filenameResult1 = config.SingleFeatureFg.h3k4me3
    filenameResult2 = config.SingleFeatureBg.h3k4me3
  case "h3k4me3o1":
    filenameData    = config.SingleFeatureData.h3k4me3o1
    filenameModel   = config.SingleFeatureJson.h3k4me3o1
    filenameComp    = config.SingleFeatureComp.h3k4me3o1
    filenameResult1 = config.SingleFeatureFg.h3k4me3o1
    filenameResult2 = config.SingleFeatureBg.h3k4me3o1
    localConfig.BinSummaryStatistics = "mean"
  case "rna":
    filenameData    = config.SingleFeatureData.rna
    filenameModel   = config.SingleFeatureJson.rna
    filenameComp    = config.SingleFeatureComp.rna
    filenameResult1 = config.SingleFeatureFg.rna
    filenameResult2 = config.SingleFeatureBg.rna
  case "rnaLow":
    filenameData    = config.SingleFeatureData.rnaLow
    filenameModel   = config.SingleFeatureJson.rnaLow
    filenameComp    = config.SingleFeatureComp.rnaLow
    filenameResult1 = config.SingleFeatureFg.rnaLow
    filenameResult2 = config.SingleFeatureBg.rnaLow
  case "control":
    filenameData    = config.SingleFeatureData.control
    filenameModel   = config.SingleFeatureJson.control
    filenameComp    = config.SingleFeatureComp.control
    filenameResult1 = config.SingleFeatureFg.control
    filenameResult2 = config.SingleFeatureBg.control
  default:
    log.Fatalf("unknown feature: %s", feature)
  }
  if updateRequired(config, filenameResult1, filenameModel, filenameComp, filenameData) ||
    (updateRequired(config, filenameResult2, filenameModel, filenameComp, filenameData)) {
    modhmm_single_feature_coverage_all(config)
    single_feature_classify(localConfig, filenameModel, filenameComp, filenameData, filenameResult1, filenameResult2)
  }
}

func modhmm_single_feature_classify_all(config ConfigModHmm) {
  for _, feature := range []string{"atac", "h3k27ac", "h3k27me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "rnaLow", "control"} {
    modhmm_single_feature_classify(config, feature)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_classify_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s classify-single-feature", os.Args[0]))
  options.SetParameters("<FEATURE>\n")

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
    modhmm_single_feature_classify_all(config)
  } else {
    modhmm_single_feature_classify(config, options.Args()[0])
  }
}
