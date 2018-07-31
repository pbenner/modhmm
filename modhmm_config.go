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
  Rna_low     string `json:"RNA low"`
}

func (config *ConfigSingleFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.ConfigCoveragePaths.CompletePaths(dir, prefix, suffix)
  config.Rna_low    = completePath(dir, prefix, config.Rna_low,    fmt.Sprintf("rna-low%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigMultiFeaturePaths struct {
  EA       string
  PA       string
  BI       string
  PR       string
  TR       string
  TL       string
  R1       string
  R2       string
  CL       string
  NS       string
}

func (config *ConfigMultiFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.PA = completePath(dir, prefix, config.PA, fmt.Sprintf("PA%s", suffix))
  config.EA = completePath(dir, prefix, config.EA, fmt.Sprintf("EA%s", suffix))
  config.BI = completePath(dir, prefix, config.BI, fmt.Sprintf("BI%s", suffix))
  config.PR = completePath(dir, prefix, config.PR, fmt.Sprintf("PR%s", suffix))
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
  BamDir                     string                   `json:"Bam Directory"`
  Bam                        ConfigBam                `json:"Bam Files"`
  CoverageBinSize            int                      `json:"Coverage Bin Size`
  CoverageThreads            int                      `json:"Coverage Threads"`
  CoverageDir                string                   `json:"Coverage Directory"`
  Coverage                   ConfigCoveragePaths      `json:"Coverage Files"`
  CoverageCnts               ConfigCoveragePaths      `json:"Coverage Counts Files"`
  SingleFeatureModelDir      string                   `json:"Single-Feature Model Directory"`
  SingleFeatureModel         ConfigSingleFeaturePaths `json:"Single-Feature Model Files"`
  SingleFeatureComp          ConfigSingleFeaturePaths `json:"Single-Feature Model Component Files"`
  SingleFeatureDir           string                   `json:"Single-Feature Directory"`
  SingleFeatureFg            ConfigSingleFeaturePaths `json:"Single-Feature Foreground"`
  SingleFeatureBg            ConfigSingleFeaturePaths `json:"Single-Feature Background"`
  SingleFeatureFgExp         ConfigSingleFeaturePaths `json:"Single-Feature Foreground [exp]"`
  SingleFeatureBgExp         ConfigSingleFeaturePaths `json:"Single-Feature Background [exp]"`
  MultiFeatureDir            string                   `json:"Multi-Feature Directory"`
  MultiFeatureProb           ConfigMultiFeaturePaths  `json:"Multi-Feature Probabilities"`
  MultiFeatureProbExp        ConfigMultiFeaturePaths  `json:"Multi-Feature Probabilities [exp]"`
  MultiFeatureProbNorm       ConfigMultiFeaturePaths  `json:"Normalized Multi-Feature Probabilities"`
  MultiFeatureProbNormExp    ConfigMultiFeaturePaths  `json:"Normalized Multi-Feature Probabilities [exp]"`
  Posterior                  ConfigMultiFeaturePaths  `json:"Posterior Marginals"`
  PosteriorExp               ConfigMultiFeaturePaths  `json:"Posterior Marginals [exp]"`
  PosteriorDir               string                   `json:"Posterior Marginals Directory"`
  ModelType                  string                   `json:"Model Type"`
  ModelUnconstrained         bool                     `json:"Model Unconstrained"`
  Model                      string                   `json:"Model File"`
  ModelDir                   string                   `json:"Model Directory"`
  Segmentation               string                   `json:"Segmentation File"`
  SegmentationDir            string                   `json:"Segmentation Directory"`
  Directory                  string
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
  if err := ImportFile(config, filename); err != nil {
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
  config.CoverageThreads      = 1
  config.CoverageBinSize      = 10
  config.ModelType            = "posterior"
  config.Threads              = 1
  config.Verbose              = 0
  return config
}

func (config *ConfigModHmm) setDefaultDir(target, def string) string {
  if target != "" {
    return target
  }
  if config.Directory != "" {
    return config.Directory
  }
  return def
}

func (config *ConfigModHmm) CompletePaths() {
  config.BamDir                 = config.setDefaultDir(config.BamDir,  "")
  config.CoverageDir            = config.setDefaultDir(config.CoverageDir, config.BamDir)
  config.SingleFeatureModelDir  = config.setDefaultDir(config.SingleFeatureModelDir, config.CoverageDir)
  config.SingleFeatureDir       = config.setDefaultDir(config.SingleFeatureDir, config.SingleFeatureModelDir)
  config.MultiFeatureDir        = config.setDefaultDir(config.MultiFeatureDir, config.SingleFeatureDir)
  config.ModelDir               = config.setDefaultDir(config.ModelDir, config.MultiFeatureDir)
  config.SegmentationDir        = config.setDefaultDir(config.SegmentationDir, config.ModelDir)
  config.PosteriorDir           = config.setDefaultDir(config.PosteriorDir, config.SegmentationDir)
  config.Model                  = completePath(config.ModelDir, "", config.Model, "segmentation.json")
  config.Segmentation           = completePath(config.SegmentationDir, "", config.Segmentation, "segmentation.bed.gz")
  config.Bam                    .CompletePaths(config.BamDir, "", "")
  config.Coverage               .CompletePaths(config.CoverageDir, "coverage-", ".bw")
  config.CoverageCnts           .CompletePaths(config.SingleFeatureModelDir, "", ".counts.json")
  config.SingleFeatureModel     .CompletePaths(config.SingleFeatureModelDir, "", ".json")
  config.SingleFeatureComp      .CompletePaths(config.SingleFeatureModelDir, "", ".components.json")
  config.SingleFeatureFg        .CompletePaths(config.SingleFeatureDir, "single-feature-", ".fg.bw")
  config.SingleFeatureFgExp     .CompletePaths(config.SingleFeatureDir, "single-feature-exp-", ".fg.bw")
  config.SingleFeatureBg        .CompletePaths(config.SingleFeatureDir, "single-feature-", ".bg.bw")
  config.SingleFeatureBgExp     .CompletePaths(config.SingleFeatureDir, "single-feature-exp-", ".bg.bw")
  config.MultiFeatureProb       .CompletePaths(config.MultiFeatureDir, "multi-feature-", ".bw")
  config.MultiFeatureProbExp    .CompletePaths(config.MultiFeatureDir, "multi-feature-exp-", ".bw")
  config.MultiFeatureProbNorm   .CompletePaths(config.MultiFeatureDir, "multi-feature-norm-", ".bw")
  config.MultiFeatureProbNormExp.CompletePaths(config.MultiFeatureDir, "multi-feature-norm-exp-", ".bw")
  config.Posterior              .CompletePaths(config.PosteriorDir, "posterior-marginal-", ".bw")
  config.PosteriorExp           .CompletePaths(config.PosteriorDir, "posterior-marginal-exp-", ".bw")
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
  fmt.Fprintf(&buffer, " -> RNA (low)            : %v %s\n", config.Rna_low,   fileCheckMark(config.Rna_low))
  fmt.Fprintf(&buffer, " -> Control              : %v %s\n", config.Control,   fileCheckMark(config.Control))

  return buffer.String()
}

func (config ConfigMultiFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> PA                   : %v %s\n", config.PA, fileCheckMark(config.PA))
  fmt.Fprintf(&buffer, " -> EA                   : %v %s\n", config.EA, fileCheckMark(config.EA))
  fmt.Fprintf(&buffer, " -> BI                   : %v %s\n", config.BI, fileCheckMark(config.BI))
  fmt.Fprintf(&buffer, " -> PR                   : %v %s\n", config.PR, fileCheckMark(config.PR))
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

  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "%v", config.SessionConfig.String())
    fmt.Fprintf(&buffer, " -> Coverage Bin Size      : %d\n\n", config.CoverageBinSize)
    fmt.Fprintf(&buffer, "Alignment files (BAM):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Bam.String())
    fmt.Fprintf(&buffer, "Coverage files (bigWig):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Coverage.String())
    fmt.Fprintf(&buffer, "Single-feature mixture distributions:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureModel.String())
    fmt.Fprintf(&buffer, "Single-feature count statistics:\n")
    fmt.Fprintf(&buffer, "%v\n", config.CoverageCnts.String())
    fmt.Fprintf(&buffer, "Single-feature foreground mixture components:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureComp.String())
    fmt.Fprintf(&buffer, "Single-feature foreground probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFg.String())
    fmt.Fprintf(&buffer, "Single-feature background probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBg.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Single-feature foreground probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFgExp.String())
    fmt.Fprintf(&buffer, "Single-feature background probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBgExp.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Multi-feature probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProb.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Multi-feature probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProbExp.String())
    fmt.Fprintf(&buffer, "Normalized multi-feature probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProbNorm.String())
    fmt.Fprintf(&buffer, "Normalized multi-feature probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProbNormExp.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Posterior marginals (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Posterior.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Posterior marginals:\n")
    fmt.Fprintf(&buffer, "%v\n", config.PosteriorExp.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "ModHmm options:\n")
    fmt.Fprintf(&buffer, " ->  Description                  : %v\n"   , config.Description)
    fmt.Fprintf(&buffer, " ->  Directory                    : %v\n"   , config.Directory)
    fmt.Fprintf(&buffer, " ->  ModHmm Model File            : %v %s\n", config.Model, fileCheckMark(config.Model))
    fmt.Fprintf(&buffer, " ->  ModHmm Model Directory       : %v\n"   , config.ModelDir)
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation File     : %v %s\n", config.Segmentation, fileCheckMark(config.Segmentation))
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation Directory: %v\n"   , config.SegmentationDir)
  }
  return buffer.String()
}
