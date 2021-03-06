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
import   "io"
import   "log"
import   "math"
import   "path"
import   "sort"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

func ImportCounts(config ConfigModHmm, filename string) Counts {
  counts := Counts{}
  printStderr(config, 1, "Importing reference counts from `%s'... ", filename)
  if err := counts.ImportFile(filename); err != nil {
    printStderr(config, 1, "failed\n")
    // remove directory from filename
    _, filename = path.Split(filename)
    filename = path.Join(config.ModelFallbackPath(), filename)
    printStderr(config, 1, "Importing counts from `%s' fallback model... ", config.ModelFallback)
    if err := ImportDefaultFile(config, &counts, filename); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  } else {
    printStderr(config, 1, "done\n")
  }
  return counts
}

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

func (obj Counts) Quantile(r float64) float64 {
  n := 0
  s := 0.0
  for i, _ := range obj.X {
    n += obj.Y[i]
  }
  for i, _ := range obj.X {
    if s >= r {
      return obj.X[i]
    }
    s += float64(obj.Y[i])/float64(n)
  }
  return obj.X[len(obj.X)-1]
}

func (obj Counts) ThresholdedMean(t float64) float64 {
  n := 0
  r := 0.0
  for i, _ := range obj.X {
    if obj.X[i] > t {
      n += obj.Y[i]
    }
  }
  for i, _ := range obj.X {
    if obj.X[i] > t {
      r += float64(obj.Y[i])/float64(n)*obj.X[i]
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

func (config *Counts) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, config)
}

func (config *Counts) Export(writer io.Writer) error {
  return JsonExport(writer, config)
}

func (config *Counts) ImportFile(filename string) error {
  if err := ImportFile(config, filename, Float64Type); err != nil {
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

func compute_counts(config ConfigModHmm, track Track) Counts {
  config.BinSummaryStatistics = "discrete mean"
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
  return c
}

/* -------------------------------------------------------------------------- */

func modhmm_compute_counts(config ConfigModHmm, track Track, filenameOut string) {
  c := compute_counts(config, track)
  printStderr(config, 1, "Exporting counts to `%s'... ", filenameOut)
  if err := c.ExportFile(filenameOut); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")
}
