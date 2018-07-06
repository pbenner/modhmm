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
import   "path"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

func completePath(prefix, mypath, def string) string {
  if mypath == "" {
    mypath = def
  }
  if !path.IsAbs(mypath) {
    mypath = fmt.Sprintf("%s/%s", prefix, mypath)
  }
  return mypath
}

/* -------------------------------------------------------------------------- */

type ConfigFeaturePaths struct {
  atac       string `json:"ATAC"`
  h3k27ac    string `json:"H3K27ac"`
  h3k27me3   string `json:"H3K27me3"`
  h3k9me3    string `json:"H3K27me3"`
  h3k4me1    string `json:"H3K4me1"`
  h3k4me3    string `json:"H3K4me3"`
  h3k4me3o1  string `json:"H3K4me3o1"`
}

func (config *ConfigFeaturePaths) CompletePaths(prefix, suffix string) {
  config.atac      = completePath(prefix, config.atac,      fmt.Sprintf("atac%s", suffix))
  config.h3k27ac   = completePath(prefix, config.h3k27ac,   fmt.Sprintf("h3k27ac%s", suffix))
  config.h3k27me3  = completePath(prefix, config.h3k27me3,  fmt.Sprintf("h3k27me3%s", suffix))
  config.h3k9me3   = completePath(prefix, config.h3k9me3,   fmt.Sprintf("h3k9me3%s", suffix))
  config.h3k4me1   = completePath(prefix, config.h3k4me1,   fmt.Sprintf("h3k4me1%s", suffix))
  config.h3k4me3   = completePath(prefix, config.h3k4me3,   fmt.Sprintf("h3k4me3%s", suffix))
  config.h3k4me3o1 = completePath(prefix, config.h3k4me3o1, fmt.Sprintf("h3k4me3o1%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigClassifierPaths struct {
  EA       string
  EP       string
  PA       string
  PB       string
}

/* -------------------------------------------------------------------------- */

type ConfigModHmm struct {
  SessionConfig
  SingleFeatureData   ConfigFeaturePaths
  SingleFeatureJson   ConfigFeaturePaths
  SingleFeaturePrefix string             `json:"Single Feature Prefix"`
}

/* -------------------------------------------------------------------------- */

func (config *ConfigModHmm) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, config)
}

func (config *ConfigModHmm) Export(writer io.Writer) error {
  return JsonExport(writer, config)
}

func (config *ConfigModHmm) ImportFile(filename string) error {
  if err := ImportFile(config, filename, BareRealType); err != nil {
    return err
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func DefaultModHmmConfig() ConfigModHmm {
  config := ConfigModHmm{}
  // set default values
  config.BinSize              = 200
  config.BinSummaryStatistics = "mean"
  config.Threads              = 1
  config.Verbose              = 0
  return config
}

func (config *ConfigModHmm) CompletePaths() {
  if config.SingleFeaturePrefix == "" {
    config.SingleFeaturePrefix = "."
  }
  config.SingleFeatureData.CompletePaths(config.SingleFeaturePrefix, ".bw")
  config.SingleFeatureJson.CompletePaths(config.SingleFeaturePrefix, ".json")
}

/* -------------------------------------------------------------------------- */

func (config ConfigFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.atac)
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.h3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.h3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.h3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.h3k4me3)
  fmt.Fprintf(&buffer, " -> H3K4me3o1            : %v\n", config.h3k4me3o1)

  return buffer.String()
}

func (config ConfigModHmm) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, "%v\n", config.SessionConfig.String())
  fmt.Fprintf(&buffer, "Input data bigWig files:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureData.String())
  fmt.Fprintf(&buffer, "Single-feature mixture distributions:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureJson.String())
  fmt.Fprintf(&buffer, "ModHmm options:\n")
  fmt.Fprintf(&buffer, " -> Single Feature Prefix: %v\n", config.SingleFeaturePrefix)

  return buffer.String()
}
