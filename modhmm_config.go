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
import   "os"

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
  Atac       []string `json:"ATAC"`
  H3k27ac    []string `json:"H3K27ac"`
  H3k27me3   []string `json:"H3K27me3"`
  H3k9me3    []string `json:"H3K9me3"`
  H3k4me1    []string `json:"H3K4me1"`
  H3k4me3    []string `json:"H3K4me3"`
  Rna        []string `json:"RNA"`
  Control    []string `json:"Control"`
}

func (config *ConfigBam) CompletePaths(prefix, suffix string) {
  for i, _ := range config.Atac {
    config.Atac[i]      = completePath(prefix, config.Atac[i],     "")
  }
  for i, _ := range config.H3k27ac {
    config.H3k27ac[i]   = completePath(prefix, config.H3k27ac[i],  "")
  }
  for i, _ := range config.H3k27me3 {
    config.H3k27me3[i]  = completePath(prefix, config.H3k27me3[i], "")
  }
  for i, _ := range config.H3k9me3 {
    config.H3k9me3[i]   = completePath(prefix, config.H3k9me3[i],  "")
  }
  for i, _ := range config.H3k4me1 {
    config.H3k4me1[i]   = completePath(prefix, config.H3k4me1[i],  "")
  }
  for i, _ := range config.H3k4me3 {
    config.H3k4me3[i]   = completePath(prefix, config.H3k4me3[i],  "")
  }
  for i, _ := range config.Rna {
    config.Rna[i]       = completePath(prefix, config.Rna[i],      "")
  }
  for i, _ := range config.Control {
    config.Control[i]   = completePath(prefix, config.Control[i],  "")
  }
}

/* -------------------------------------------------------------------------- */

type ConfigSingleFeaturePaths struct {
  Atac       string `json:"ATAC"`
  H3k27ac    string `json:"H3K27ac"`
  H3k27me3   string `json:"H3K27me3"`
  H3k9me3    string `json:"H3K27me3"`
  H3k4me1    string `json:"H3K4me1"`
  H3k4me3    string `json:"H3K4me3"`
  H3k4me3o1  string `json:"H3K4me3o1"`
  Rna        string `json:"RNA"`
  RnaLow     string `json:"RNA low"`
  Control    string `json:"Control"`
}

func (config *ConfigSingleFeaturePaths) CompletePaths(prefix, suffix string) {
  config.Atac      = completePath(prefix, config.Atac,      fmt.Sprintf("atac%s", suffix))
  config.H3k27ac   = completePath(prefix, config.H3k27ac,   fmt.Sprintf("h3k27ac%s", suffix))
  config.H3k27me3  = completePath(prefix, config.H3k27me3,  fmt.Sprintf("h3k27me3%s", suffix))
  config.H3k9me3   = completePath(prefix, config.H3k9me3,   fmt.Sprintf("h3k9me3%s", suffix))
  config.H3k4me1   = completePath(prefix, config.H3k4me1,   fmt.Sprintf("h3k4me1%s", suffix))
  config.H3k4me3   = completePath(prefix, config.H3k4me3,   fmt.Sprintf("h3k4me3%s", suffix))
  config.H3k4me3o1 = completePath(prefix, config.H3k4me3o1, fmt.Sprintf("h3k4me3o1%s", suffix))
  config.Rna       = completePath(prefix, config.Rna,       fmt.Sprintf("rna%s", suffix))
  config.RnaLow    = completePath(prefix, config.RnaLow,    fmt.Sprintf("rna-low%s", suffix))
  config.Control   = completePath(prefix, config.Control,   fmt.Sprintf("control%s", suffix))
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
  config.EA = completePath(prefix, config.EA, fmt.Sprintf("multifeature-EA%s", suffix))
  config.EP = completePath(prefix, config.EP, fmt.Sprintf("multifeature-EP%s", suffix))
  config.PA = completePath(prefix, config.PA, fmt.Sprintf("multifeature-PA%s", suffix))
  config.PB = completePath(prefix, config.PB, fmt.Sprintf("multifeature-PB%s", suffix))
  config.TR = completePath(prefix, config.TR, fmt.Sprintf("multifeature-TR%s", suffix))
  config.TL = completePath(prefix, config.TL, fmt.Sprintf("multifeature-TL%s", suffix))
  config.R1 = completePath(prefix, config.R1, fmt.Sprintf("multifeature-R1%s", suffix))
  config.R2 = completePath(prefix, config.R2, fmt.Sprintf("multifeature-R2%s", suffix))
  config.CL = completePath(prefix, config.CL, fmt.Sprintf("multifeature-CL%s", suffix))
  config.NS = completePath(prefix, config.NS, fmt.Sprintf("multifeature-NS%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigModHmm struct {
  SessionConfig
  ThreadsCoverage            int                      `json:"Threads Coverage"`
  SingleFeatureBamDir        string                   `json:"Bam Directory"`
  SingleFeatureBam           ConfigBam                `json:"Bam Files"`
  SingleFeatureDataDir       string                   `json:"Coverage Directory"`
  SingleFeatureData          ConfigSingleFeaturePaths
  SingleFeatureJsonDir       string                   `json:"Model Directory"`
  SingleFeatureJson          ConfigSingleFeaturePaths
  SingleFeatureComp          ConfigSingleFeaturePaths
  SingleFeatureFg            ConfigSingleFeaturePaths
  SingleFeatureBg            ConfigSingleFeaturePaths
  MultiFeatureClass          ConfigMultiFeaturePaths
  MultiFeatureClassExp       ConfigMultiFeaturePaths
  Prefix                     string                    `json:"Prefix"`
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
  if config.Prefix == "" {
    config.Prefix = "."
  }
  if config.SingleFeatureBamDir == "" {
    config.SingleFeatureBamDir = config.Prefix
  }
  if config.SingleFeatureDataDir == "" {
    config.SingleFeatureDataDir = config.Prefix
  }
  if config.SingleFeatureJsonDir == "" {
    config.SingleFeatureJsonDir = config.Prefix
  }
  if config.Model == "" {
    config.Model = completePath(config.Prefix, config.Model, "segmentation.json")
  }
  if config.Segmentation == "" {
    config.Segmentation = completePath(config.Prefix, config.Segmentation, "segmentation.bed.gz")
  }
  config.SingleFeatureBam    .CompletePaths(config.SingleFeatureBamDir, "")
  config.SingleFeatureData   .CompletePaths(config.SingleFeatureDataDir, ".bw")
  config.SingleFeatureJson   .CompletePaths(config.SingleFeatureJsonDir, ".json")
  config.SingleFeatureComp   .CompletePaths(config.SingleFeatureJsonDir, ".components.json")
  config.SingleFeatureFg     .CompletePaths(config.Prefix, ".fg.bw")
  config.SingleFeatureBg     .CompletePaths(config.Prefix, ".bg.bw")
  config.MultiFeatureClass   .CompletePaths(config.Prefix, ".bw")
  config.MultiFeatureClassExp.CompletePaths(config.Prefix, ".exp.bw")
}

/* -------------------------------------------------------------------------- */

func fileCheckMark(filename string) string {
  if _, err := os.Stat(filename); err != nil {
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
  fmt.Fprintf(&buffer, "Single-feature foreground classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFg.String())
  fmt.Fprintf(&buffer, "Single-feature background classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBg.String())
  fmt.Fprintf(&buffer, "Multi-feature classifications (log-scale):\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureClass.String())
  fmt.Fprintf(&buffer, "Multi-feature classifications:\n")
  fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureClassExp.String())
  fmt.Fprintf(&buffer, "ModHmm options:\n")
  fmt.Fprintf(&buffer, " ->  ModHMM Model File           : %v %s\n", config.Model, fileCheckMark(config.Model))
  fmt.Fprintf(&buffer, " ->  Genome Segmentation File    : %v %s\n", config.Segmentation, fileCheckMark(config.Segmentation))
  fmt.Fprintf(&buffer, " ->  Description                 : %v\n", config.Description)

  return buffer.String()
}
