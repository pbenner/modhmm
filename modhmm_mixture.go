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
import   "bytes"
import   "io"
import   "io/ioutil"
import   "log"
import   "math"
import   "path"

import . "github.com/pbenner/ngstat/config"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/logarithmetic"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"

import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

func ImportDefaultFile(config ConfigModHmm, object Serializable, filename string, args... interface{}) error {
  // remove directory from filename
  _, filename = path.Split(filename)
  filename = path.Join(config.SingleFeatureModelFallbackPath(), filename)

  f, err := assets.Open(filename)
  if err != nil {
    return err
  }
  str, err := ioutil.ReadAll(f)
  if err != nil {
    return err
  }
  return object.Import(bytes.NewReader(str), args...)
}

func ImportDefaultDistribution(config ConfigModHmm, filename string, distribution BasicDistribution, t ScalarType) error {
  // remove directory from filename
  _, filename = path.Split(filename)
  filename = path.Join(config.SingleFeatureModelFallbackPath(), filename)

  cfg := ConfigDistribution{}

  f, err := assets.Open(filename)
  if err != nil {
    return err
  }
  if err := cfg.ReadJson(f); err != nil {
    return err
  }
  if err := distribution.ImportConfig(cfg, t); err != nil {
    return err
  }
  return nil
}

func ImportMixtureDistribution(config ConfigModHmm, filename string) *scalarDistribution.Mixture {
  mixture := &scalarDistribution.Mixture{}
  printStderr(config, 1, "Importing mixture model from `%s'... ", filename)
  if err := ImportDistribution(filename, mixture, BareRealType); err != nil {
    printStderr(config, 1, "failed\n")
    printStderr(config, 1, "Importing `%s' fallback mixture model... ", config.SingleFeatureModelFallback)
    if err := ImportDefaultDistribution(config, filename, mixture, BareRealType); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  } else {
    printStderr(config, 1, "done\n")
  }
  return mixture
}

/* -------------------------------------------------------------------------- */

type Components []int

func (obj Components) Invert(n int) []int {
  m := make(map[int]bool)
  r := []int{}
  for _, j := range obj {
    m[j] = true
  }
  for j := 0; j < n; j++ {
    if _, ok := m[j]; !ok {
      r = append(r, j)
    }
  }
  return r
}

func (obj Components) Check(n int) error {
  m := make(map[int]struct{})
  for _, c := range obj {
    if c < 0 {
      return fmt.Errorf("invalid negative component `%d'", c)
    }
    if c >= n {
      return fmt.Errorf("component `%d' is too large; mixture distribution has only `%d' components", c, n)
    }
    if _, ok := m[c]; ok {
      return fmt.Errorf("component `%d' appears multiple times", c)
    }
    m[c] = struct{}{}
  }
  return nil
}

func (obj *Components) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, obj)
}

func (obj *Components) Export(writer io.Writer) error {
  return JsonExport(writer, obj)
}

func ImportComponents(config ConfigModHmm, filename string, n int) ([]int, []int) {
  var k Components
  printStderr(config, 1, "Importing foreground components from `%s'... ", filename)
  if err := ImportFile(&k, filename); err != nil {
    printStderr(config, 1, "failed\n")
    printStderr(config, 1, "Importing foreground components from `%s' fallback model... ", config.SingleFeatureModelFallback)
    if err := ImportDefaultFile(config, &k, filename); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
  }
  if err := k.Check(n); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("ERROR: invalid components file `%s': %v", filename, err)
  } else {
    printStderr(config, 1, "done\n")
  }
  r := Components(k).Invert(n)
  return []int(k), []int(r)
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
  mixture := ImportMixtureDistribution(config, filenameModel)

  k, r := ImportComponents(config, filenameComp, mixture.NComponents())

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
