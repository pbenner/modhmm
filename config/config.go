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

package config

/* -------------------------------------------------------------------------- */

import   "fmt"
import   "bytes"
import   "log"
import   "io"
import   "path"
import   "path/filepath"
import   "strings"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/modhmm/utility"

/* -------------------------------------------------------------------------- */

var CoverageList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "rna", "control"}

/* -------------------------------------------------------------------------- */

var SingleFeatureModelList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "control"}

/* -------------------------------------------------------------------------- */

var SingleFeatureList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "rna-low", "control"}

/* -------------------------------------------------------------------------- */

var MultiFeatureList = StringList{
  "pa", "ea", "bi", "pr", "tr", "tl", "r1", "r2", "ns", "cl"}

/* -------------------------------------------------------------------------- */

type TargetFile struct {
  Filename string
  Static   bool
}

func (obj TargetFile) String() string {
  if obj.Static {
    return fmt.Sprintf("%s %s [static]", obj.Filename, FileCheckMark(obj.Filename))
  } else {
    return fmt.Sprintf("%s %s", obj.Filename, FileCheckMark(obj.Filename))
  }
}

/* -------------------------------------------------------------------------- */

type SingleFeatureFiles struct {
  Feature      string
  Foreground   TargetFile
  Background   TargetFile
  Model        TargetFile
  Components   TargetFile
  Coverage   []TargetFile
  Counts     []TargetFile
}

func (obj SingleFeatureFiles) Dependencies() []string {
  filenames := []string{}
  filenames  = append(filenames, obj.Model     .Filename)
  filenames  = append(filenames, obj.Components.Filename)
  for _, file := range obj.Coverage {
    filenames = append(filenames, file.Filename)
  }
  for _, file := range obj.Counts {
    filenames = append(filenames, file.Filename)
  }
  return filenames
}

/* -------------------------------------------------------------------------- */

func completePath(dir, prefix, mypath, def string) string {
  if mypath == "" && def != "" {
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
  Dnase      []string `json:"DNase"`
  H3k27ac    []string `json:"H3K27ac"`
  H3k27me3   []string `json:"H3K27me3"`
  H3k9me3    []string `json:"H3K9me3"`
  H3k4me1    []string `json:"H3K4me1"`
  H3k4me3    []string `json:"H3K4me3"`
  Rna        []string `json:"RNA"`
  Control    []string `json:"Control"`
}

func (config *ConfigBam) GetFilenames(feature string) []string {
  switch strings.ToLower(feature) {
  case "atac"     : return config.Atac
  case "dnase"    : return config.Dnase
  case "h3k27ac"  : return config.H3k27ac
  case "h3k27me3" : return config.H3k27me3
  case "h3k9me3"  : return config.H3k9me3
  case "h3k4me1"  : return config.H3k4me1
  case "h3k4me3"  : return config.H3k4me3
  case "rna"      : return config.Rna
  case "control"  : return config.Control
  default:
    panic("internal error")
  }
}

func (config *ConfigBam) CompletePaths(dir, prefix, suffix string) {
  for i, _ := range config.Atac {
    config.Atac[i]      = completePath(dir, prefix, config.Atac[i],     "")
  }
  for i, _ := range config.Dnase {
    config.Dnase[i]     = completePath(dir, prefix, config.Dnase[i],    "")
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
  Open       TargetFile `json:"-"`
  Atac       TargetFile `json:"ATAC"`
  Dnase      TargetFile `json:"DNase"`
  H3k27ac    TargetFile `json:"H3K27ac"`
  H3k27me3   TargetFile `json:"H3K27me3"`
  H3k9me3    TargetFile `json:"H3K27me3"`
  H3k4me1    TargetFile `json:"H3K4me1"`
  H3k4me3    TargetFile `json:"H3K4me3"`
  Rna        TargetFile `json:"RNA"`
  Control    TargetFile `json:"Control"`
}

func (config *ConfigCoveragePaths) GetTargetFile(feature string) TargetFile {
  switch strings.ToLower(feature) {
  case "open"     : return config.Open
  case "atac"     : return config.Atac
  case "dnase"    : return config.Dnase
  case "h3k27ac"  : return config.H3k27ac
  case "h3k27me3" : return config.H3k27me3
  case "h3k9me3"  : return config.H3k9me3
  case "h3k4me1"  : return config.H3k4me1
  case "h3k4me3"  : return config.H3k4me3
  case "rna"      : return config.Rna
  case "control"  : return config.Control
  default:
    panic("internal error")
  }
}

func (config *ConfigCoveragePaths) CompletePaths(dir, prefix, suffix string) {
  atac  := completePath(dir, prefix, config.Atac .Filename, fmt.Sprintf("atac%s", suffix))
  dnase := completePath(dir, prefix, config.Dnase.Filename, fmt.Sprintf("dnase%s", suffix))
  if config.Atac.Filename == "" && config.Dnase.Filename == "" {
    config.Atac .Filename = atac
    config.Dnase.Filename = dnase
  }
  if config.Atac.Filename == "" {
    config.Open .Filename = dnase
    config.Atac .Filename = ""
    config.Dnase.Filename = dnase
  }
  if config.Dnase.Filename == "" {
    config.Open .Filename = atac
    config.Atac .Filename = atac
    config.Dnase.Filename = ""
  }
  config.Atac     .Filename = completePath(dir, prefix, config.Atac     .Filename, fmt.Sprintf("atac%s", suffix))
  config.Dnase    .Filename = completePath(dir, prefix, config.Dnase    .Filename, fmt.Sprintf("dnase%s", suffix))
  config.H3k27ac  .Filename = completePath(dir, prefix, config.H3k27ac  .Filename, fmt.Sprintf("h3k27ac%s", suffix))
  config.H3k27me3 .Filename = completePath(dir, prefix, config.H3k27me3 .Filename, fmt.Sprintf("h3k27me3%s", suffix))
  config.H3k9me3  .Filename = completePath(dir, prefix, config.H3k9me3  .Filename, fmt.Sprintf("h3k9me3%s", suffix))
  config.H3k4me1  .Filename = completePath(dir, prefix, config.H3k4me1  .Filename, fmt.Sprintf("h3k4me1%s", suffix))
  config.H3k4me3  .Filename = completePath(dir, prefix, config.H3k4me3  .Filename, fmt.Sprintf("h3k4me3%s", suffix))
  config.Rna      .Filename = completePath(dir, prefix, config.Rna      .Filename, fmt.Sprintf("rna%s", suffix))
  config.Control  .Filename = completePath(dir, prefix, config.Control  .Filename, fmt.Sprintf("control%s", suffix))
}

func (config *ConfigCoveragePaths) GetFilenames() []string {
  filenames := []string{}
  for _, feature := range CoverageList {
    filenames = append(filenames, config.GetTargetFile(feature).Filename)
  }
  return filenames
}

func (config *ConfigCoveragePaths) SetStatic(static bool) {
  config.Atac     .Static = static
  config.Dnase    .Static = static
  config.Open     .Static = static
  config.H3k27ac  .Static = static
  config.H3k27me3 .Static = static
  config.H3k9me3  .Static = static
  config.H3k4me1  .Static = static
  config.H3k4me3  .Static = static
  config.Rna      .Static = static
  config.Control  .Static = static
}

/* -------------------------------------------------------------------------- */

type ConfigCountsPaths struct {
  ConfigCoveragePaths
  H3k4me3o1   TargetFile `json:"H3K4me3o1"`
  Rna_low     TargetFile `json:"RNA low"`
}

func (config *ConfigCountsPaths) GetTargetFile(feature string) TargetFile {
  switch strings.ToLower(feature) {
  case "h3K4me3o1": return config.H3k4me3o1
  case "rna-low"  : return config.Rna_low
  default         : return config.ConfigCoveragePaths.GetTargetFile(feature)
  }
}

func (config *ConfigCountsPaths) CompletePaths(dir, prefix, suffix string) {
  config.ConfigCoveragePaths.CompletePaths(dir, prefix, suffix)
  config.H3k4me3o1.Filename = completePath(dir, prefix, config.H3k4me3o1.Filename, fmt.Sprintf("h3k4me3o1%s", suffix))
}

func (config *ConfigCountsPaths) GetFilenames() []string {
  filenames := []string{}
  for _, feature := range SingleFeatureList {
    filenames = append(filenames, config.GetTargetFile(feature).Filename)
  }
  return filenames
}

func (config *ConfigCountsPaths) SetStatic(static bool) {
  config.ConfigCoveragePaths.SetStatic(static)
  config.H3k4me3o1.Static = static
}

/* -------------------------------------------------------------------------- */

type ConfigSingleFeaturePaths struct {
  ConfigCoveragePaths
  H3k4me3o1   TargetFile `json:"H3K4me3o1"`
  Rna_low     TargetFile `json:"RNA low"`
}

func (config *ConfigSingleFeaturePaths) GetTargetFile(feature string) TargetFile {
  switch strings.ToLower(feature) {
  case "h3k4me3o1": return config.H3k4me3o1
  case "rna_low"  : return config.Rna_low
  case "rna-low"  : return config.Rna_low
  default         : return config.ConfigCoveragePaths.GetTargetFile(feature)
  }
}

func (config *ConfigSingleFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.ConfigCoveragePaths.CompletePaths(dir, prefix, suffix)
  config.H3k4me3o1.Filename = completePath(dir, prefix, config.H3k4me3o1.Filename, fmt.Sprintf("h3k4me3o1%s", suffix))
  config.Rna_low  .Filename = completePath(dir, prefix, config.Rna_low  .Filename, fmt.Sprintf("rna-low%s", suffix))
}

func (config *ConfigSingleFeaturePaths) GetFilenames() []string {
  filenames := []string{}
  for _, feature := range SingleFeatureList {
    filenames = append(filenames, config.GetTargetFile(feature).Filename)
  }
  return filenames
}

func (config *ConfigSingleFeaturePaths) SetStatic(static bool) {
  config.ConfigCoveragePaths.SetStatic(static)
  config.H3k4me3o1.Static = static
  config.Rna_low  .Static = static
}

/* -------------------------------------------------------------------------- */

type ConfigMultiFeaturePaths struct {
  EA       TargetFile
  PA       TargetFile
  BI       TargetFile
  PR       TargetFile
  TR       TargetFile
  TL       TargetFile
  R1       TargetFile
  R2       TargetFile
  CL       TargetFile
  NS       TargetFile
}

func (config *ConfigMultiFeaturePaths) GetTargetFile(state string) TargetFile {
  switch strings.ToUpper(state) {
  case "EA": return config.EA
  case "PA": return config.PA
  case "BI": return config.BI
  case "PR": return config.PR
  case "TR": return config.TR
  case "TL": return config.TL
  case "R1": return config.R1
  case "R2": return config.R2
  case "CL": return config.CL
  case "NS": return config.NS
  default:
    panic("internal error")
  }
}

func (config *ConfigMultiFeaturePaths) CompletePaths(dir, prefix, suffix string) {
  config.PA.Filename = completePath(dir, prefix, config.PA.Filename, fmt.Sprintf("PA%s", suffix))
  config.EA.Filename = completePath(dir, prefix, config.EA.Filename, fmt.Sprintf("EA%s", suffix))
  config.BI.Filename = completePath(dir, prefix, config.BI.Filename, fmt.Sprintf("BI%s", suffix))
  config.PR.Filename = completePath(dir, prefix, config.PR.Filename, fmt.Sprintf("PR%s", suffix))
  config.TR.Filename = completePath(dir, prefix, config.TR.Filename, fmt.Sprintf("TR%s", suffix))
  config.TL.Filename = completePath(dir, prefix, config.TL.Filename, fmt.Sprintf("TL%s", suffix))
  config.R1.Filename = completePath(dir, prefix, config.R1.Filename, fmt.Sprintf("R1%s", suffix))
  config.R2.Filename = completePath(dir, prefix, config.R2.Filename, fmt.Sprintf("R2%s", suffix))
  config.CL.Filename = completePath(dir, prefix, config.CL.Filename, fmt.Sprintf("CL%s", suffix))
  config.NS.Filename = completePath(dir, prefix, config.NS.Filename, fmt.Sprintf("NS%s", suffix))
}

/* -------------------------------------------------------------------------- */

type ConfigModHmm struct {
  SessionConfig
  OpenChromatinAssay         string                   `json:"Open Chromatin Assay"`
  BamDir                     string                   `json:"Bam Directory"`
  Bam                        ConfigBam                `json:"Bam Files"`
  CoverageBinSize            int                      `json:"Coverage Bin Size`
  CoverageThreads            int                      `json:"Coverage Threads"`
  CoverageDir                string                   `json:"Coverage Directory"`
  Coverage                   ConfigCoveragePaths      `json:"Coverage Files"`
  CoverageCnts               ConfigCountsPaths        `json:"Coverage Counts Files"`
  CoverageMAPQ               int                      `json:"Coverage MAPQ"`
  SingleFeatureModelDir      string                   `json:"Single-Feature Model Directory"`
  SingleFeatureModel         ConfigSingleFeaturePaths `json:"Single-Feature Model Files"`
  SingleFeatureComp          ConfigSingleFeaturePaths `json:"Single-Feature Model Component Files"`
  SingleFeatureDir           string                   `json:"Single-Feature Directory"`
  SingleFeatureFg            ConfigSingleFeaturePaths `json:"Single-Feature Foreground"`
  SingleFeatureBg            ConfigSingleFeaturePaths `json:"Single-Feature Background"`
  SingleFeatureFgExp         ConfigSingleFeaturePaths `json:"Single-Feature Foreground [exp]"`
  SingleFeatureBgExp         ConfigSingleFeaturePaths `json:"Single-Feature Background [exp]"`
  SingleFeaturePeak          ConfigSingleFeaturePaths `json:"Single-Feature Peaks"`
  MultiFeatureDir            string                   `json:"Multi-Feature Directory"`
  MultiFeatureProb           ConfigMultiFeaturePaths  `json:"Multi-Feature Probabilities"`
  MultiFeaturePeak           ConfigMultiFeaturePaths  `json:"Multi-Feature Peaks"`
  MultiFeatureProbExp        ConfigMultiFeaturePaths  `json:"Multi-Feature Probabilities [exp]"`
  MultiFeatureProbNorm       ConfigMultiFeaturePaths  `json:"Normalized Multi-Feature Probabilities"`
  MultiFeatureProbNormExp    ConfigMultiFeaturePaths  `json:"Normalized Multi-Feature Probabilities [exp]"`
  Posterior                  ConfigMultiFeaturePaths  `json:"Posterior Marginals"`
  PosteriorExp               ConfigMultiFeaturePaths  `json:"Posterior Marginals [exp]"`
  PosteriorPeak              ConfigMultiFeaturePaths  `json:"Posterior Marginals Peaks"`
  PosteriorDir               string                   `json:"Posterior Marginals Directory"`
  ModelUnconstrained         bool                     `json:"Model Unconstrained"`
  Model                      TargetFile               `json:"Model File"`
  ModelDir                   string                   `json:"Model Directory"`
  Segmentation               TargetFile               `json:"Segmentation File"`
  SegmentationDir            string                   `json:"Segmentation Directory"`
  Directory                  string
  Description                string
  FontSize                   float64
  XLim                    [2]float64
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
  config.CoverageMAPQ         = 30
  config.FontSize             = 12
  config.OpenChromatinAssay   = ""
  config.Threads              = 1
  config.Verbose              = 0
  return config
}

func (config *ConfigModHmm) setDefaultDir(prefix, target, def string) string {
  if target != "" {
    if path.IsAbs(target) {
      return target
    } else {
      return path.Join(prefix, target)
    }
  }
  if config.Directory != "" {
    if path.IsAbs(config.Directory) {
      return config.Directory
    } else {
      return path.Join(prefix, config.Directory)
    }
  }
  return def
}

func (config *ConfigModHmm) CoerceOpenChromatinAssay(feature string) string {
  switch strings.ToLower(feature) {
  case "atac" :
    if config.OpenChromatinAssay == "dnase" {
      log.Fatalf("unknown feature: %s", feature)
    }
    feature = "open"
  case "dnase":
    if config.OpenChromatinAssay == "atac" {
      log.Fatalf("unknown feature: %s", feature)
    }
    feature = "open"
  }
  return feature
}

func (config *ConfigModHmm) DetectOpenChromatinAssay() string {
  switch strings.ToLower(config.OpenChromatinAssay) {
  case "atac" : return "atac"
  case "dnase": return "dnase"
  case "":
  default:
    log.Fatalf("invalid open chromatin assay `%s'", config.OpenChromatinAssay)
  }
  if len(config.Bam.Atac) != 0 && len(config.Bam.Dnase) != 0 {
    log.Fatal("Config file specifies BAM files for ATAC and DNase-seq! Please select a single open chromatin assay.")
  }
  if len(config.Bam.Atac) != 0 {
    return "atac"
  }
  if len(config.Bam.Dnase) != 0 {
    return "dnase"
  }
  if config.Coverage.Atac.Filename != "" && config.Coverage.Dnase.Filename != "" {
    if FileExists(config.Coverage.Atac.Filename) && FileExists(config.Coverage.Dnase.Filename) {
      log.Fatal("Coverage bigWig files exist for both ATAC- and DNase-seq. Please select a single open chromatin assay.")
    }
  }
  if config.Coverage.Atac.Filename != "" && FileExists(config.Coverage.Atac.Filename) {
    return "atac"
  }
  if config.Coverage.Dnase.Filename != "" && FileExists(config.Coverage.Dnase.Filename) {
    return "dnase"
  }
  if config.SingleFeatureFg.Atac.Filename != "" && config.SingleFeatureFg.Dnase.Filename != "" {
    if FileExists(config.SingleFeatureFg.Atac.Filename) && FileExists(config.SingleFeatureFg.Dnase.Filename) {
      log.Fatal("SingleFeatureFg bigWig files exist for both ATAC- and DNase-seq. Please select a single open chromatin assay.")
    }
  }
  if config.SingleFeatureFg.Atac.Filename != "" && FileExists(config.SingleFeatureFg.Atac.Filename) {
    return "atac"
  }
  if config.SingleFeatureFg.Dnase.Filename != "" && FileExists(config.SingleFeatureFg.Dnase.Filename) {
    return "dnase"
  }
  // return default assay
  return "atac"
}

func (config *ConfigModHmm) SetOpenChromatinAssay(assay string) {
  switch strings.ToLower(assay) {
  case "atac":
    config.Coverage          .Open = config.Coverage          .Atac
    config.CoverageCnts      .Open = config.CoverageCnts      .Atac
    config.SingleFeatureModel.Open = config.SingleFeatureModel.Atac
    config.SingleFeatureComp .Open = config.SingleFeatureComp .Atac
    config.SingleFeaturePeak .Open = config.SingleFeaturePeak .Atac
    config.SingleFeatureFg   .Open = config.SingleFeatureFg   .Atac
    config.SingleFeatureFgExp.Open = config.SingleFeatureFgExp.Atac
    config.SingleFeatureBg   .Open = config.SingleFeatureBg   .Atac
    config.SingleFeatureBgExp.Open = config.SingleFeatureBgExp.Atac
  case "dnase":
    config.Coverage          .Open = config.Coverage          .Dnase
    config.CoverageCnts      .Open = config.CoverageCnts      .Dnase
    config.SingleFeatureModel.Open = config.SingleFeatureModel.Dnase
    config.SingleFeatureComp .Open = config.SingleFeatureComp .Dnase
    config.SingleFeaturePeak .Open = config.SingleFeaturePeak .Dnase
    config.SingleFeatureFg   .Open = config.SingleFeatureFg   .Dnase
    config.SingleFeatureFgExp.Open = config.SingleFeatureFgExp.Dnase
    config.SingleFeatureBg   .Open = config.SingleFeatureBg   .Dnase
    config.SingleFeatureBgExp.Open = config.SingleFeatureBgExp.Dnase
  default:
    log.Fatalf("invalid open chromatin assay `%s'", assay)
  }
  config.OpenChromatinAssay = assay
}

func (config *ConfigModHmm) CompletePaths(prefix string) {
  config.BamDir                 = config.setDefaultDir(prefix, config.BamDir               ,  "")
  config.CoverageDir            = config.setDefaultDir(prefix, config.CoverageDir          , config.BamDir)
  config.SingleFeatureModelDir  = config.setDefaultDir(prefix, config.SingleFeatureModelDir, config.CoverageDir)
  config.SingleFeatureDir       = config.setDefaultDir(prefix, config.SingleFeatureDir     , config.SingleFeatureModelDir)
  config.MultiFeatureDir        = config.setDefaultDir(prefix, config.MultiFeatureDir      , config.SingleFeatureDir)
  config.ModelDir               = config.setDefaultDir(prefix, config.ModelDir             , config.MultiFeatureDir)
  config.SegmentationDir        = config.setDefaultDir(prefix, config.SegmentationDir      , config.ModelDir)
  config.PosteriorDir           = config.setDefaultDir(prefix, config.PosteriorDir         , config.SegmentationDir)
  config.Model.Filename         = completePath(config.ModelDir, "", config.Model.Filename, "segmentation.json")
  config.Segmentation.Filename  = completePath(config.SegmentationDir, "", config.Segmentation.Filename, "segmentation.bed.gz")
  config.Bam                    .CompletePaths(config.BamDir, "", "")
  config.Coverage               .CompletePaths(config.CoverageDir, "coverage-", ".bw")
  config.CoverageCnts           .CompletePaths(config.SingleFeatureModelDir, "", ".counts.json")
  config.SingleFeatureModel     .CompletePaths(config.SingleFeatureModelDir, "", ".json")
  config.SingleFeatureComp      .CompletePaths(config.SingleFeatureModelDir, "", ".components.json")
  config.SingleFeaturePeak      .CompletePaths(config.SingleFeatureDir, "single-feature-peaks-", ".table")
  config.SingleFeatureFg        .CompletePaths(config.SingleFeatureDir, "single-feature-", ".fg.bw")
  config.SingleFeatureFgExp     .CompletePaths(config.SingleFeatureDir, "single-feature-exp-", ".fg.bw")
  config.SingleFeatureBg        .CompletePaths(config.SingleFeatureDir, "single-feature-", ".bg.bw")
  config.SingleFeatureBgExp     .CompletePaths(config.SingleFeatureDir, "single-feature-exp-", ".bg.bw")
  config.MultiFeatureProb       .CompletePaths(config.MultiFeatureDir, "multi-feature-", ".bw")
  config.MultiFeaturePeak       .CompletePaths(config.MultiFeatureDir, "multi-feature-peaks-", ".table")
  config.MultiFeatureProbExp    .CompletePaths(config.MultiFeatureDir, "multi-feature-exp-", ".bw")
  config.MultiFeatureProbNorm   .CompletePaths(config.MultiFeatureDir, "multi-feature-norm-", ".bw")
  config.MultiFeatureProbNormExp.CompletePaths(config.MultiFeatureDir, "multi-feature-norm-exp-", ".bw")
  config.Posterior              .CompletePaths(config.PosteriorDir, "posterior-marginal-", ".bw")
  config.PosteriorExp           .CompletePaths(config.PosteriorDir, "posterior-marginal-exp-", ".bw")
  config.PosteriorPeak          .CompletePaths(config.PosteriorDir, "posterior-marginal-peaks-", ".bw")
  config.SetOpenChromatinAssay(config.DetectOpenChromatinAssay())
}

func (config *ConfigModHmm) SingleFeatureFiles(feature string, logScale bool) SingleFeatureFiles {

  if !SingleFeatureList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
  files := SingleFeatureFiles{}
  files.Feature = config.CoerceOpenChromatinAssay(strings.ToLower(feature))

  switch files.Feature {
  case "rna-low":
    if logScale {
      files.Foreground = config.SingleFeatureFg.Rna_low
      files.Background = config.SingleFeatureBg.Rna_low
    } else {
      files.Foreground = config.SingleFeatureFgExp.Rna_low
      files.Background = config.SingleFeatureBgExp.Rna_low
    }
    files.Model      = config.SingleFeatureModel.Rna
    files.Components = config.SingleFeatureComp .Rna_low
    files.Coverage   = append(files.Coverage, config.Coverage    .Rna)
    files.Counts     = append(files.Counts,   config.CoverageCnts.Rna)
  default:
    if logScale {
      files.Foreground = config.SingleFeatureFg.GetTargetFile(feature)
      files.Background = config.SingleFeatureBg.GetTargetFile(feature)
    } else {
      files.Foreground = config.SingleFeatureFgExp.GetTargetFile(feature)
      files.Background = config.SingleFeatureBgExp.GetTargetFile(feature)
    }
    files.Model      = config.SingleFeatureModel.GetTargetFile(feature)
    files.Components = config.SingleFeatureComp .GetTargetFile(feature)
    if files.Feature == "h3k4me3o1" {
      files.Coverage = append(files.Coverage, config.Coverage    .H3k4me1)
      files.Coverage = append(files.Coverage, config.Coverage    .H3k4me3)
      files.Counts   = append(files.Counts,   config.CoverageCnts.H3k4me1)
      files.Counts   = append(files.Counts,   config.CoverageCnts.H3k4me3)
    } else {
      files.Coverage = append(files.Coverage, config.Coverage    .GetTargetFile(feature))
      files.Counts   = append(files.Counts,   config.CoverageCnts.GetTargetFile(feature))
    }
  }
  return files
}

/* -------------------------------------------------------------------------- */

func (config ConfigBam) String(openChromatinAssay string) string {
  var buffer bytes.Buffer

  switch strings.ToLower(openChromatinAssay) {
  case "atac":
    fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.Atac)
  case "dnase":
    fmt.Fprintf(&buffer, " -> DNase                : %v\n", config.Dnase)
  default:
    panic("internal error")
  }
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.H3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.H3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.H3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.H3k4me3)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.Rna)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigCoveragePaths) String(openChromatinAssay string) string {
  var buffer bytes.Buffer

  switch strings.ToLower(openChromatinAssay) {
  case "atac":
    fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.Atac)
  case "dnase":
    fmt.Fprintf(&buffer, " -> DNase                : %v\n", config.Dnase)
  default:
    panic("internal error")
  }
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.H3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.H3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.H3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.H3k4me3)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.Rna)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigSingleFeaturePaths) String(openChromatinAssay string) string {
  var buffer bytes.Buffer

  switch strings.ToLower(openChromatinAssay) {
  case "atac":
    fmt.Fprintf(&buffer, " -> ATAC                 : %v\n", config.Atac)
  case "dnase":
    fmt.Fprintf(&buffer, " -> DNase                : %v\n", config.Dnase)
  default:
    panic("internal error")
  }
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.H3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.H3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.H3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.H3k4me3)
  fmt.Fprintf(&buffer, " -> H3K4me3o1            : %v\n", config.H3k4me3o1)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.Rna)
  fmt.Fprintf(&buffer, " -> RNA (low)            : %v\n", config.Rna_low)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigMultiFeaturePaths) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> PA                   : %v\n", config.PA)
  fmt.Fprintf(&buffer, " -> EA                   : %v\n", config.EA)
  fmt.Fprintf(&buffer, " -> BI                   : %v\n", config.BI)
  fmt.Fprintf(&buffer, " -> PR                   : %v\n", config.PR)
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

  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "%v", config.SessionConfig.String())
    fmt.Fprintf(&buffer, " -> Open Chromatin Assay   : %s\n", config.OpenChromatinAssay)
    fmt.Fprintf(&buffer, " -> Coverage Bin Size      : %d\n\n", config.CoverageBinSize)
    fmt.Fprintf(&buffer, "Alignment files (BAM):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Bam.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Coverage files (bigWig):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Coverage.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature mixture distributions:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureModel.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature count statistics:\n")
    fmt.Fprintf(&buffer, "%v\n", config.CoverageCnts.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature foreground mixture components:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureComp.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature foreground probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFg.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature background probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBg.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Single-feature peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeaturePeak.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature foreground probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureFgExp.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Single-feature background probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.SingleFeatureBgExp.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Multi-feature probabilities (log-scale):\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeatureProb.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Multi-feature peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.MultiFeaturePeak.String())
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
    fmt.Fprintf(&buffer, "Posterior marginals peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.PosteriorPeak.String())
    fmt.Fprintf(&buffer, "Posterior marginals:\n")
    fmt.Fprintf(&buffer, "%v\n", config.PosteriorExp.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "ModHmm options:\n")
    fmt.Fprintf(&buffer, " ->  Description                  : %v\n", config.Description)
    fmt.Fprintf(&buffer, " ->  Directory                    : %v\n", config.Directory)
    fmt.Fprintf(&buffer, " ->  ModHmm Model File            : %v\n", config.Model)
    fmt.Fprintf(&buffer, " ->  ModHmm Model Directory       : %v\n", config.ModelDir)
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation File     : %v\n", config.Segmentation)
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation Directory: %v\n", config.SegmentationDir)
  }
  return buffer.String()
}
