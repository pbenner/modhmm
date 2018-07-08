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

type ConfigBam struct {
  atac       []string `json:"ATAC"`
  h3k27ac    []string `json:"H3K27ac"`
  h3k27me3   []string `json:"H3K27me3"`
  h3k9me3    []string `json:"H3K27me3"`
  h3k4me1    []string `json:"H3K4me1"`
  h3k4me3    []string `json:"H3K4me3"`
  rna        []string `json:"RNA"`
  control    []string `json:"Control"`
}

/* -------------------------------------------------------------------------- */

type ConfigSingleFeaturePaths struct {
  atac       string `json:"ATAC"`
  h3k27ac    string `json:"H3K27ac"`
  h3k27me3   string `json:"H3K27me3"`
  h3k9me3    string `json:"H3K27me3"`
  h3k4me1    string `json:"H3K4me1"`
  h3k4me3    string `json:"H3K4me3"`
  h3k4me3o1  string `json:"H3K4me3o1"`
  rna        string `json:"RNA"`
  rnaLow     string `json:"RNA low"`
  control    string `json:"Control"`
}

func (config *ConfigSingleFeaturePaths) CompletePaths(prefix, suffix string) {
  config.atac      = completePath(prefix, config.atac,      fmt.Sprintf("atac%s", suffix))
  config.h3k27ac   = completePath(prefix, config.h3k27ac,   fmt.Sprintf("h3k27ac%s", suffix))
  config.h3k27me3  = completePath(prefix, config.h3k27me3,  fmt.Sprintf("h3k27me3%s", suffix))
  config.h3k9me3   = completePath(prefix, config.h3k9me3,   fmt.Sprintf("h3k9me3%s", suffix))
  config.h3k4me1   = completePath(prefix, config.h3k4me1,   fmt.Sprintf("h3k4me1%s", suffix))
  config.h3k4me3   = completePath(prefix, config.h3k4me3,   fmt.Sprintf("h3k4me3%s", suffix))
  config.h3k4me3o1 = completePath(prefix, config.h3k4me3o1, fmt.Sprintf("h3k4me3o1%s", suffix))
  config.rna       = completePath(prefix, config.rna,       fmt.Sprintf("rna%s", suffix))
  config.rnaLow    = completePath(prefix, config.rnaLow,    fmt.Sprintf("rna-low%s", suffix))
  config.control   = completePath(prefix, config.control,   fmt.Sprintf("control%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigMultiFeaturePaths struct {
  EA       string
  EP       string
  PA       string
  PB       string
  TR       string
  TL       string
  R1       string
  R2       string
  CL       string
  NS       string
}

func (config *ConfigMultiFeaturePaths) CompletePaths(prefix, suffix string) {
  config.EA = completePath(prefix, config.EA, fmt.Sprintf("classification-EA%s", suffix))
  config.EP = completePath(prefix, config.EP, fmt.Sprintf("classification-EP%s", suffix))
  config.PA = completePath(prefix, config.PA, fmt.Sprintf("classification-PA%s", suffix))
  config.PB = completePath(prefix, config.PB, fmt.Sprintf("classification-PB%s", suffix))
  config.TR = completePath(prefix, config.TR, fmt.Sprintf("classification-TR%s", suffix))
  config.TL = completePath(prefix, config.TL, fmt.Sprintf("classification-TL%s", suffix))
  config.R1 = completePath(prefix, config.R1, fmt.Sprintf("classification-R1%s", suffix))
  config.R2 = completePath(prefix, config.R2, fmt.Sprintf("classification-R2%s", suffix))
  config.CL = completePath(prefix, config.CL, fmt.Sprintf("classification-CL%s", suffix))
  config.NS = completePath(prefix, config.NS, fmt.Sprintf("classification-NS%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigModHmm struct {
  SessionConfig
  SingleFeatureBam     ConfigBam
  SingleFeatureData    ConfigSingleFeaturePaths
  SingleFeatureJson    ConfigSingleFeaturePaths
  SingleFeatureComp    ConfigSingleFeaturePaths
  SingleFeatureFg      ConfigSingleFeaturePaths
  SingleFeatureBg      ConfigSingleFeaturePaths
  MultiFeatureClass    ConfigMultiFeaturePaths
  MultiFeatureClassExp ConfigMultiFeaturePaths
  Prefix string                          `json:"Prefix"`
  SingleFeaturePrefix string             `json:"Single Feature Prefix"`
  SingleFeatureMixturePrefix string      `json:"Single Feature Mixture Prefix"`
   MultiFeaturePrefix string             `json:"Multi Feature Prefix"`
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
  if config.Prefix == "" {
    config.Prefix = "."
  }
  if config.SingleFeaturePrefix == "" {
    config.SingleFeaturePrefix = config.Prefix
  }
  if config.SingleFeatureMixturePrefix == "" {
    config.SingleFeatureMixturePrefix = config.Prefix
  }
  if config.MultiFeaturePrefix == "" {
    config.MultiFeaturePrefix = config.Prefix
  }
  config.SingleFeatureData   .CompletePaths(config.SingleFeaturePrefix, ".bw")
  config.SingleFeatureJson   .CompletePaths(config.SingleFeatureMixturePrefix, ".json")
  config.SingleFeatureComp   .CompletePaths(config.SingleFeatureMixturePrefix, ".components.json")
  config.SingleFeatureFg     .CompletePaths(config.SingleFeaturePrefix, ".fg.bw")
  config.SingleFeatureBg     .CompletePaths(config.SingleFeaturePrefix, ".bg.bw")
  config.MultiFeatureClass   .CompletePaths(config.MultiFeaturePrefix, ".bw")
  config.MultiFeatureClassExp.CompletePaths(config.MultiFeaturePrefix, ".exp.bw")
}

/* -------------------------------------------------------------------------- */

func (config ConfigBam) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.atac)
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.h3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.h3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.h3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.h3k4me3)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.rna)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.control)

  return buffer.String()
}

func (config ConfigSingleFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.atac)
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.h3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.h3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.h3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.h3k4me3)
  fmt.Fprintf(&buffer, " -> H3K4me3o1            : %v\n", config.h3k4me3o1)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.rna)
  fmt.Fprintf(&buffer, " -> RNA (low)            : %v\n", config.rnaLow)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.control)

  return buffer.String()
}

func (config ConfigMultiFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> PA                   : %v\n", config.PA)
  fmt.Fprintf(&buffer, " -> PB                   : %v\n", config.PB)
  fmt.Fprintf(&buffer, " -> EA                   : %v\n", config.EA)
  fmt.Fprintf(&buffer, " -> EP                   : %v\n", config.EP)
  fmt.Fprintf(&buffer, " -> TR                   : %v\n", config.TR)
  fmt.Fprintf(&buffer, " -> TL                   : %v\n", config.TL)
  fmt.Fprintf(&buffer, " -> R1                   : %v\n", config.R1)
  fmt.Fprintf(&buffer, " -> R2                   : %v\n", config.R2)
  fmt.Fprintf(&buffer, " -> CL                   : %v\n", config.CL)
  fmt.Fprintf(&buffer, " -> NS                   : %v\n", config.NS)

  return buffer.String()
}

func (config ConfigModHmm) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, "%v\n", config.SessionConfig.String())
  fmt.Fprintf(&buffer, "Alignment files (BAM):\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBam.String())
  fmt.Fprintf(&buffer, "Coverage files (bigWig:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureData.String())
  fmt.Fprintf(&buffer, "Single-feature mixture distributions:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureJson.String())
  fmt.Fprintf(&buffer, "Single-feature foreground mixture components:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureComp.String())
  fmt.Fprintf(&buffer, "Single-feature foreground classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFg.String())
  fmt.Fprintf(&buffer, "Single-feature background classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBg.String())
  fmt.Fprintf(&buffer, "Multi-feature classifications (log-scale):\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureClass.String())
  fmt.Fprintf(&buffer, "Multi-feature classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureClassExp.String())
  fmt.Fprintf(&buffer, "ModHmm options:\n")
  fmt.Fprintf(&buffer, " -> Single Feature Prefix        : %v\n", config.Prefix)
  fmt.Fprintf(&buffer, " -> Single Feature Prefix        : %v\n", config.SingleFeaturePrefix)
  fmt.Fprintf(&buffer, " -> Single Feature Mixture Prefix: %v\n", config.SingleFeaturePrefix)
  fmt.Fprintf(&buffer, " ->  Multi Feature Prefix        : %v\n", config. MultiFeaturePrefix)

  return buffer.String()
}
