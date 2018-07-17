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
import   "io"
import   "log"
import   "math"
import   "os"
import   "sort"
import   "strings"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/autodiff"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

type Counts struct {
  X []float64
  Y []int
}

/* -------------------------------------------------------------------------- */

func (obj Counts) Len() int {
  return len(obj.X)
}

func (obj Counts) Less(i, j int) bool {
  return obj.X[i] < obj.X[j]
}

func (obj Counts) Swap(i, j int) {
  obj.X[i], obj.X[j] = obj.X[j], obj.X[i]
  obj.Y[i], obj.Y[j] = obj.Y[j], obj.Y[i]
}

/* -------------------------------------------------------------------------- */

func (config *Counts) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, config)
}

func (config *Counts) Export(writer io.Writer) error {
  return JsonExport(writer, config)
}

func (config *Counts) ImportFile(filename string) error {
  if err := ImportFile(config, filename, BareRealType); err != nil {
    return err
  }
  return nil
}

func (config *Counts) ExportFile(filename string) error {
  if err := ExportFile(config, filename); err != nil {
    return err
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func compute_counts(config ConfigModHmm, filenameIn, filenameOut string) {
  if track, err := ImportTrack(config.SessionConfig, filenameIn); err != nil {
    log.Fatal(err)
  } else {
    m := make(map[float64]int)
    if err := (GenericMutableTrack{}).Map(track, func(seqname string, position int, value float64) float64 {
      if !math.IsNaN(value) {
        m[value] += 1
      }
      return 0.0
    }); err != nil {
      log.Fatal(err)
    }
    i  := 0
    c  := Counts{}
    c.X = make([]float64, len(m))
    c.Y = make([]int,     len(m))
    for k, v := range m {
      c.X[i] = k
      c.Y[i] = v
      i++
    }
    sort.Sort(c)

    printStderr(config, 1, "Exporting counts to `%s'... ", filenameOut)
    if err := c.ExportFile(filenameOut); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_compute_counts(config ConfigModHmm, feature string) {
  if !coverageList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }

  filenameIn  := getFieldAsString(config.SingleFeatureData, strings.ToLower(feature))
  filenameOut := getFieldAsString(config.SingleFeatureCnts, strings.ToLower(feature))

  if strings.ToLower(feature) != "h3k4me3o1" {
    config.BinSummaryStatistics = "discrete mean"
  }
  if updateRequired(config, filenameOut, filenameOut) {
    modhmm_coverage_all(config)
    compute_counts(config, filenameIn, filenameOut)
  }
}

func modhmm_compute_counts_all(config ConfigModHmm) {
  for _, feature := range coverageList {
    modhmm_compute_counts(config, feature)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_compute_counts_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s compute-counts", os.Args[0]))

  optHelp := options.   BoolLong("help",        'h',     "print help")

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
    modhmm_compute_counts_all(config)
  } else {
    modhmm_compute_counts(config, options.Args()[0])
  }
}
