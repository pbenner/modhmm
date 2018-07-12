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
import   "path/filepath"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

func completePath(dir, prefix, mypath, def string) string {
  if mypath == "" {
    mypath = fmt.Sprintf("%s%s", prefix, def)
  }
  if d, f := path.Split(mypath); d == "" {
    mypath = filepath.Join(dir, f)
  }
  return mypath
}

/* -------------------------------------------------------------------------- */

type ConfigBam struct {
  Atac       []string `json:"ATAC"`
  H3k27ac    []string `json:"H3K27ac"`
  H3k27me3   []string `json:"H3K27me3"`
  H3k9me3    []string `json:"H3K9me3"`
  H3k4me1    []string `json:"H3K4me1"`
  H3k4me3    []string `json:"H3K4me3"`
  Rna        []string `json:"RNA"`
  Control    []string `json:"Control"`
}

func (config *ConfigBam) CompletePaths(dir, prefix, suffix string) {
  for i, _ := range config.Atac {
    config.Atac[i]      = completePath(dir, prefix, config.Atac[i],     "")
  }
  for i, _ := range config.H3k27ac {
    config.H3k27ac[i]   = completePath(dir, prefix, config.H3k27ac[i],  "")
  }
  for i, _ := range config.H3k27me3 {
    config.H3k27me3[i]  = completePath(dir, prefix, config.H3k27me3[i], "")
  }
  for i, _ := range config.H3k9me3 {
    config.H3k9me3[i]   = completePath(dir, prefix, config.H3k9me3[i],  "")
  }
  for i, _ := range config.H3k4me1 {
    config.H3k4me1[i]   = completePath(dir, prefix, config.H3k4me1[i],  "")
  }
  for i, _ := range config.H3k4me3 {
    config.H3k4me3[i]   = completePath(dir, prefix, config.H3k4me3[i],  "")
  }
  for i, _ := range config.Rna {
    config.Rna[i]       = completePath(dir, prefix, config.Rna[i],      "")
  }
  for i, _ := range config.Control {
    config.Control[i]   = completePath(dir, prefix, config.Control[i],  "")
  }
}

/* -------------------------------------------------------------------------- */

type ConfigCoveragePaths struct {
  Atac       string `json:"ATAC"`
  H3k27ac    string `json:"H3K27ac"`
  H3k27me3   string `json:"H3K27me3"`
  H3k9me3    string `json:"H3K27me3"`
  H3k4me1    string `json:"H3K4me1"`
  H3k4me3    string `json:"H3K4me3"`
  H3k4me3o1  string `json:"H3K4me3o1"`
  Rna        string `json:"RNA"`
  Control    string `json:"Control"`
}

func (config *ConfigCoveragePaths) CompletePaths(dir, prefix, suffix string) {
  config.Atac      = completePath(dir, prefix, config.Atac,      fmt.Sprintf("atac%s", suffix))
  config.H3k27ac   = completePath(dir, prefix, config.H3k27ac,   fmt.Sprintf("h3k27ac%s", suffix))
  config.H3k27me3  = completePath(dir, prefix, config.H3k27me3,  fmt.Sprintf("h3k27me3%s", suffix))
  config.H3k9me3   = completePath(dir, prefix, config.H3k9me3,   fmt.Sprintf("h3k9me3%s", suffix))
  config.H3k4me1   = completePath(dir, prefix, config.H3k4me1,   fmt.Sprintf("h3k4me1%s", suffix))
  config.H3k4me3   = completePath(dir, prefix, config.H3k4me3,   fmt.Sprintf("h3k4me3%s", suffix))
  config.H3k4me3o1 = completePath(dir, prefix, config.H3k4me3o1, fmt.Sprintf("h3k4me3o1%s", suffix))
  config.Rna       = completePath(dir, prefix, config.Rna,       fmt.Sprintf("rna%s", suffix))
  config.Control   = completePath(dir, prefix, config.Control,   fmt.Sprintf("control%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigSingleFeaturePaths struct {
  ConfigCoveragePaths
  RnaLow     string `json:"RNA low"`
}

func (config *ConfigSingleFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.ConfigCoveragePaths.CompletePaths(dir, prefix, suffix)
  config.RnaLow    = completePath(dir, prefix, config.RnaLow,    fmt.Sprintf("rna-low%s", suffix))
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

func (config *ConfigMultiFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.EA = completePath(dir, prefix, config.EA, fmt.Sprintf("EA%s", suffix))
  config.EP = completePath(dir, prefix, config.EP, fmt.Sprintf("EP%s", suffix))
  config.PA = completePath(dir, prefix, config.PA, fmt.Sprintf("PA%s", suffix))
  config.PB = completePath(dir, prefix, config.PB, fmt.Sprintf("PB%s", suffix))
  config.TR = completePath(dir, prefix, config.TR, fmt.Sprintf("TR%s", suffix))
  config.TL = completePath(dir, prefix, config.TL, fmt.Sprintf("TL%s", suffix))
  config.R1 = completePath(dir, prefix, config.R1, fmt.Sprintf("R1%s", suffix))
  config.R2 = completePath(dir, prefix, config.R2, fmt.Sprintf("R2%s", suffix))
  config.CL = completePath(dir, prefix, config.CL, fmt.Sprintf("CL%s", suffix))
  config.NS = completePath(dir, prefix, config.NS, fmt.Sprintf("NS%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigModHmm struct {
  SessionConfig
  ThreadsCoverage            int                      `json:"Threads Coverage"`
  SingleFeatureBamDir        string                   `json:"Bam Directory"`
  SingleFeatureBam           ConfigBam                `json:"Bam Files"`
  SingleFeatureDataDir       string                   `json:"Coverage Directory"`
  SingleFeatureData          ConfigCoveragePaths      `json:"Coverage Files"`
  SingleFeatureJsonDir       string                   `json:"Model Directory"`
  SingleFeatureJson          ConfigSingleFeaturePaths `json:"Model Files"`
  SingleFeatureComp          ConfigSingleFeaturePaths `json:"Model Component Files"`
  SingleFeatureFg            ConfigSingleFeaturePaths
  SingleFeatureBg            ConfigSingleFeaturePaths
  MultiFeatureProb           ConfigMultiFeaturePaths
  MultiFeatureProbExp        ConfigMultiFeaturePaths
  Type                       string                    `json:"Type"`
  Directory                  string                    `json:"Directory"`
  Model                      string                    `json:"ModHmm Model File"`
  Segmentation               string                    `json:"Genome Segmentation File"`
  Description                string
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
  config.ThreadsCoverage      = 1
  config.Verbose              = 0
  return config
}

func (config *ConfigModHmm) CompletePaths() {
  if config.Directory == "" {
    config.Directory = "."
  }
  if config.SingleFeatureBamDir == "" {
    config.SingleFeatureBamDir = config.Directory
  }
  if config.SingleFeatureDataDir == "" {
    config.SingleFeatureDataDir = config.Directory
  }
  if config.SingleFeatureJsonDir == "" {
    config.SingleFeatureJsonDir = config.Directory
  }
  if config.Model == "" {
    config.Model = completePath(config.Directory, "", config.Model, "segmentation.json")
  }
  if config.Segmentation == "" {
    config.Segmentation = completePath(config.Directory, "", config.Segmentation, "segmentation.bed.gz")
  }
  config.SingleFeatureBam    .CompletePaths(config.SingleFeatureBamDir, "", "")
  config.SingleFeatureData   .CompletePaths(config.SingleFeatureDataDir, "coverage-", ".bw")
  config.SingleFeatureJson   .CompletePaths(config.SingleFeatureJsonDir, "", ".json")
  config.SingleFeatureComp   .CompletePaths(config.SingleFeatureJsonDir, "", ".components.json")
  config.SingleFeatureFg     .CompletePaths(config.Directory, "single-feature-", ".fg.bw")
  config.SingleFeatureBg     .CompletePaths(config.Directory, "single-feature-", ".bg.bw")
  config.MultiFeatureProb   .CompletePaths(config.Directory, "multi-feature-", ".bw")
  config.MultiFeatureProbExp.CompletePaths(config.Directory, "multi-feature-", ".exp.bw")
}

/* -------------------------------------------------------------------------- */

func fileCheckMark(filename string) string {
  if !fileExists(filename) {
    return "\xE2\x9C\x97"
  } else {
    return "\xE2\x9C\x93"
  }
}

/* -------------------------------------------------------------------------- */

func (config ConfigBam) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.Atac)
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.H3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.H3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.H3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.H3k4me3)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.Rna)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigCoveragePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v %s\n", config.Atac,      fileCheckMark(config.Atac))
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v %s\n", config.H3k27ac,   fileCheckMark(config.H3k27ac))
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v %s\n", config.H3k27me3,  fileCheckMark(config.H3k27me3))
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v %s\n", config.H3k4me1,   fileCheckMark(config.H3k4me1))
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v %s\n", config.H3k4me3,   fileCheckMark(config.H3k4me3))
  fmt.Fprintf(&buffer, " -> H3K4me3o1            : %v %s\n", config.H3k4me3o1, fileCheckMark(config.H3k4me3o1))
  fmt.Fprintf(&buffer, " -> RNA                  : %v %s\n", config.Rna,       fileCheckMark(config.Rna))
  fmt.Fprintf(&buffer, " -> Control              : %v %s\n", config.Control,   fileCheckMark(config.Control))

  return buffer.String()
}

func (config ConfigSingleFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> ATAC                 : %v %s\n", config.Atac,      fileCheckMark(config.Atac))
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v %s\n", config.H3k27ac,   fileCheckMark(config.H3k27ac))
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v %s\n", config.H3k27me3,  fileCheckMark(config.H3k27me3))
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v %s\n", config.H3k4me1,   fileCheckMark(config.H3k4me1))
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v %s\n", config.H3k4me3,   fileCheckMark(config.H3k4me3))
  fmt.Fprintf(&buffer, " -> H3K4me3o1            : %v %s\n", config.H3k4me3o1, fileCheckMark(config.H3k4me3o1))
  fmt.Fprintf(&buffer, " -> RNA                  : %v %s\n", config.Rna,       fileCheckMark(config.Rna))
  fmt.Fprintf(&buffer, " -> RNA (low)            : %v %s\n", config.RnaLow,    fileCheckMark(config.RnaLow))
  fmt.Fprintf(&buffer, " -> Control              : %v %s\n", config.Control,   fileCheckMark(config.Control))

  return buffer.String()
}

func (config ConfigMultiFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> PA                   : %v %s\n", config.PA, fileCheckMark(config.PA))
  fmt.Fprintf(&buffer, " -> PB                   : %v %s\n", config.PB, fileCheckMark(config.PB))
  fmt.Fprintf(&buffer, " -> EA                   : %v %s\n", config.EA, fileCheckMark(config.EA))
  fmt.Fprintf(&buffer, " -> EP                   : %v %s\n", config.EP, fileCheckMark(config.EP))
  fmt.Fprintf(&buffer, " -> TR                   : %v %s\n", config.TR, fileCheckMark(config.TR))
  fmt.Fprintf(&buffer, " -> TL                   : %v %s\n", config.TL, fileCheckMark(config.TL))
  fmt.Fprintf(&buffer, " -> R1                   : %v %s\n", config.R1, fileCheckMark(config.R1))
  fmt.Fprintf(&buffer, " -> R2                   : %v %s\n", config.R2, fileCheckMark(config.R2))
  fmt.Fprintf(&buffer, " -> CL                   : %v %s\n", config.CL, fileCheckMark(config.CL))
  fmt.Fprintf(&buffer, " -> NS                   : %v %s\n", config.NS, fileCheckMark(config.NS))

  return buffer.String()
}

func (config ConfigModHmm) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, "%v\n", config.SessionConfig.String())
  fmt.Fprintf(&buffer, "Alignment files (BAM):\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBam.String())
  fmt.Fprintf(&buffer, "Coverage files (bigWig):\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureData.String())
  fmt.Fprintf(&buffer, "Single-feature mixture distributions:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureJson.String())
  fmt.Fprintf(&buffer, "Single-feature foreground mixture components:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureComp.String())
  fmt.Fprintf(&buffer, "Single-feature foreground probabilities:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFg.String())
  fmt.Fprintf(&buffer, "Single-feature background probabilities:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBg.String())
  fmt.Fprintf(&buffer, "Multi-feature probabilities (log-scale):\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProb.String())
  fmt.Fprintf(&buffer, "Multi-feature probabilities:\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProbExp.String())
  fmt.Fprintf(&buffer, "ModHmm options:\n")
  fmt.Fprintf(&buffer, " ->  ModHMM Model File           : %v %s\n", config.Model, fileCheckMark(config.Model))
  fmt.Fprintf(&buffer, " ->  Genome Segmentation File    : %v %s\n", config.Segmentation, fileCheckMark(config.Segmentation))
  fmt.Fprintf(&buffer, " ->  Description                 : %v\n", config.Description)

  return buffer.String()
}
