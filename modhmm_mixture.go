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

import . "github.com/pbenner/ngstat/config"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/logarithmetic"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"

/* -------------------------------------------------------------------------- */

func ImportMixture(config ConfigModHmm, filenameModel string) *scalarDistribution.Mixture {
  mixture := &scalarDistribution.Mixture{}

  printStderr(config, 1, "Importing mixture model from `%s'... ", filenameModel)
  if err := ImportDistribution(filenameModel, mixture, BareRealType); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  return mixture
}

/* -------------------------------------------------------------------------- */

type Components []int

func (obj Components) Invert() []int {
  m := make(map[int]bool)
  r := []int{}
  for _, j := range obj {
    m[j] = true
  }
  for j := 0; j < len(obj); j++ {
    if _, ok := m[j]; !ok {
      r = append(r, j)
    }
  }
  return r
}

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
    log.Fatalf("ERROR: could not import components from `%s': %v", filename, err)
  }
  printStderr(config, 1, "done\n")
  return []int(k)
}

func ExportComponents(config ConfigModHmm, filename string, k []int) {
  printStderr(config, 1, "Exporting foreground components to `%s'... ", filename)
  if err := ExportFile((*Components)(&k), filename); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("ERROR: could not export components to `%s': %v", filename, err)
  }
  printStderr(config, 1, "done\n")
}

/* -------------------------------------------------------------------------- */

func ImportMixtureWeights(config ConfigModHmm, filenameModel, filenameComp string) (float64, float64) {
  mixture := ImportMixture(config, filenameModel)

  k := ImportComponents(config, filenameComp)
  r := Components(k).Invert()

  p := math.Inf(-1)
  q := math.Inf(-1)

  for _, i := range k {
    p = LogAdd(p, mixture.LogWeights.ValueAt(i))
  }
  for _, i := range r {
    q = LogAdd(q, mixture.LogWeights.ValueAt(i))
  }
  return p, q
}
