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

var EnrichmentModelList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "rna", "control"}

/* -------------------------------------------------------------------------- */

var EnrichmentList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "rna", "rna-low", "control"}

/* -------------------------------------------------------------------------- */

var ChromatinStateList = StringList{
  "pa", "ea", "bi", "pr", "tr", "tl", "r1", "r2", "ns", "cl"}

/* -------------------------------------------------------------------------- */

func EnrichmentIsOptional(feature string) bool {
  switch strings.ToLower(feature) {
  case "h3k27me3": return true
  case "h3k9me3" : return true
  case "control" : return true
  default        : return false
  }
}

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

type EnrichmentFiles struct {
  Feature           string
  Probabilities     TargetFile
  Model             TargetFile
  Components        TargetFile
  Coverage          TargetFile
  CoverageCnts      TargetFile
  // H3K4me3 source coverage and counts files
  SrcCoverage     []TargetFile
  SrcCoverageCnts []TargetFile
}

func (obj EnrichmentFiles) Dependencies() []string {
  filenames := []string{}
  filenames  = append(filenames, obj.Model       .Filename)
  filenames  = append(filenames, obj.Coverage    .Filename)
  filenames  = append(filenames, obj.CoverageCnts.Filename)
  filenames  = append(filenames, obj.Components  .Filename)
  return filenames
}

func (obj EnrichmentFiles) DependenciesModel() []string {
  return []string{obj.Coverage.Filename}
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

func (config *ConfigBam) GetTargetFiles(feature string) []string {
  switch strings.ToLower(feature) {
  case "open"    : return append(config.Atac, config.Dnase...)
  case "atac"    : return config.Atac
  case "dnase"   : return config.Dnase
  case "h3k27ac" : return config.H3k27ac
  case "h3k27me3": return config.H3k27me3
  case "h3k9me3" : return config.H3k9me3
  case "h3k4me1" : return config.H3k4me1
  case "h3k4me3" : return config.H3k4me3
  case "rna"     : return config.Rna
  case "control" : return config.Control
  default:
    panic("internal error")
  }
}

func (config *ConfigBam) GetFilenames() []string {
  filenames := []string{}
  for _, feature := range CoverageList {
    filenames = append(filenames, config.GetTargetFiles(feature)...)
  }
  return filenames
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
  Open     TargetFile `json:"-"`
  Atac     TargetFile `json:"ATAC"`
  Dnase    TargetFile `json:"DNase"`
  H3k27ac  TargetFile `json:"H3K27ac"`
  H3k27me3 TargetFile `json:"H3K27me3"`
  H3k9me3  TargetFile `json:"H3K9me3"`
  H3k4me1  TargetFile `json:"H3K4me1"`
  H3k4me3  TargetFile `json:"H3K4me3"`
  Rna      TargetFile `json:"RNA"`
  Control  TargetFile `json:"Control"`
}

func (config *ConfigCoveragePaths) GetTargetFile(feature string) TargetFile {
  switch strings.ToLower(feature) {
  case "open"    : return config.Open
  case "atac"    : return config.Atac
  case "dnase"   : return config.Dnase
  case "h3k27ac" : return config.H3k27ac
  case "h3k27me3": return config.H3k27me3
  case "h3k9me3" : return config.H3k9me3
  case "h3k4me1" : return config.H3k4me1
  case "h3k4me3" : return config.H3k4me3
  case "rna"     : return config.Rna
  case "control" : return config.Control
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
  config.Atac    .Static = static
  config.Dnase   .Static = static
  config.Open    .Static = static
  config.H3k27ac .Static = static
  config.H3k27me3.Static = static
  config.H3k9me3 .Static = static
  config.H3k4me1 .Static = static
  config.H3k4me3 .Static = static
  config.Rna     .Static = static
  config.Control .Static = static
}

/* -------------------------------------------------------------------------- */

type ConfigEnrichmentPaths struct {
  ConfigCoveragePaths
  Rna_low     TargetFile `json:"RNA low"`
}

func (config *ConfigEnrichmentPaths) GetTargetFile(feature string) TargetFile {
  switch strings.ToLower(feature) {
  case "rna_low"  : return config.Rna_low
  case "rna-low"  : return config.Rna_low
  default         : return config.ConfigCoveragePaths.GetTargetFile(feature)
  }
}

func (config *ConfigEnrichmentPaths) CompletePaths(dir, prefix, suffix string) {
  config.ConfigCoveragePaths.CompletePaths(dir, prefix, suffix)
  config.Rna_low.Filename = completePath(dir, prefix, config.Rna_low  .Filename, fmt.Sprintf("rna-low%s", suffix))
}

func (config *ConfigEnrichmentPaths) GetFilenames() []string {
  filenames := []string{}
  for _, feature := range EnrichmentList {
    filenames = append(filenames, config.GetTargetFile(feature).Filename)
  }
  return filenames
}

func (config *ConfigEnrichmentPaths) SetStatic(static bool) {
  config.ConfigCoveragePaths.SetStatic(static)
  config.Rna_low.Static = static
}

/* -------------------------------------------------------------------------- */

type ConfigEnrichmentParameters struct {
  Open     []float64 `json:"Open"`
  H3k27ac  []float64 `json:"H3K27ac"`
  H3k27me3 []float64 `json:"H3K27me3"`
  H3k9me3  []float64 `json:"H3K9me3"`
  H3k4me1  []float64 `json:"H3K4me1"`
  H3k4me3  []float64 `json:"H3K4me3"`
  Rna      []float64 `json:"RNA"`
  Control  []float64 `json:"Control"`
}

func (config *ConfigEnrichmentParameters) getParameters(feature string) []float64 {
  switch strings.ToLower(feature) {
  case "open"    : return config.Open
  case "atac"    : return config.Open
  case "dnase"   : return config.Open
  case "h3k27ac" : return config.H3k27ac
  case "h3k27me3": return config.H3k27me3
  case "h3k9me3" : return config.H3k9me3
  case "h3k4me1" : return config.H3k4me1
  case "h3k4me3" : return config.H3k4me3
  case "rna"     : return config.Rna
  case "control" : return config.Control
  default:
    panic("internal error")
  }
}

func (config *ConfigEnrichmentParameters) GetParameters(feature string) []float64 {
  parameters := config.getParameters(feature)
  if len(parameters) != 3 {
    log.Fatalf("config file has invalid number of enrichment parameters for feature `%s'", feature)
  }
  return parameters
}

/* -------------------------------------------------------------------------- */

type ConfigChromatinStatePaths struct {
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

func (config *ConfigChromatinStatePaths) GetTargetFile(state string) TargetFile {
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

func (config *ConfigChromatinStatePaths) CompletePaths(dir, prefix, suffix string) {
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
  OpenChromatinAssay      string                     `json:"Open Chromatin Assay"`
  BamDir                  string                     `json:"Bam Directory"`
  Bam                     ConfigBam                  `json:"Bam Files"`
  CoverageBinSize         int                        `json:"Coverage Bin Size`
  CoverageThreads         int                        `json:"Coverage Threads"`
  CoverageDir             string                     `json:"Coverage Directory"`
  Coverage                ConfigCoveragePaths        `json:"Coverage Files"`
  CoverageCnts            ConfigCoveragePaths        `json:"Coverage Counts Files"`
  CoverageFraglen         bool                       `json:"Coverage Fraglen"`
  CoverageMAPQ            int                        `json:"Coverage MAPQ"`
  EnrichmentMethod        string                     `json:"Enrichment Method"`
  EnrichmentModelDir      string                     `json:"Enrichment Model Directory"`
  EnrichmentModel         ConfigEnrichmentPaths      `json:"Enrichment Model Files"`
  EnrichmentComp          ConfigEnrichmentPaths      `json:"Enrichment Model Component Files"`
  EnrichmentModelStatic   bool                       `json:"Enrichment Model Static"`
  EnrichmentDir           string                     `json:"Enrichment Directory"`
  EnrichmentProb          ConfigEnrichmentPaths      `json:"Enrichment Probabilities"`
  EnrichmentPeak          ConfigEnrichmentPaths      `json:"Enrichment Peaks"`
  EnrichmentParameters    ConfigEnrichmentParameters `json:"Enrichment Parameters"`
  ChromatinStateDir       string                     `json:"Chromatin-State Directory"`
  ChromatinStateProb      ConfigChromatinStatePaths  `json:"Chromatin-State Probabilities"`
  ChromatinStatePeak      ConfigChromatinStatePaths  `json:"Chromatin-State Peaks"`
  PosteriorProb           ConfigChromatinStatePaths  `json:"Posterior Marginals"`
  PosteriorPeak           ConfigChromatinStatePaths  `json:"Posterior Marginals Peaks"`
  PosteriorDir            string                     `json:"Posterior Marginals Directory"`
  ModelEstimate           bool                       `json:"Model Estimate"`
  ModelFallback           string                     `json:"Model Fallback"`
  ModelUnconstrained      bool                       `json:"Model Unconstrained"`
  Model                   TargetFile                 `json:"Model File"`
  ModelDir                string                     `json:"Model Directory"`
  Segmentation            TargetFile                 `json:"Segmentation File"`
  SegmentationDir         string                     `json:"Segmentation Directory"`
  Directory               string
  Description             string
  FontSize                float64
  XLim                 [2]float64
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
  config.CoverageFraglen      = false
  config.CoverageBinSize      = 10
  config.CoverageMAPQ         = 30
  config.ModelFallback        = "mm10"
  config.FontSize             = 12
  config.OpenChromatinAssay   = ""
  config.EnrichmentMethod     = "heuristic"
  config.Threads              = 1
  config.Verbose              = 0
  // default parameters for assigning enrichment probabilities
  config.EnrichmentParameters.Open     = []float64{0.95, 0.01, 0.90}
  config.EnrichmentParameters.H3k27ac  = []float64{0.95, 0.01, 0.90}
  config.EnrichmentParameters.H3k27me3 = []float64{0.95, 0.01, 0.90}
  config.EnrichmentParameters.H3k9me3  = []float64{0.95, 0.01, 0.90}
  config.EnrichmentParameters.H3k4me1  = []float64{0.80, 0.01, 0.90}
  config.EnrichmentParameters.H3k4me3  = []float64{0.95, 0.01, 0.40}
  config.EnrichmentParameters.Rna      = []float64{0.60, 0.01, 0.90}
  config.EnrichmentParameters.Control  = []float64{0.95, 0.01, 0.90}
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
  if config.EnrichmentProb.Atac.Filename != "" && config.EnrichmentProb.Dnase.Filename != "" {
    if FileExists(config.EnrichmentProb.Atac.Filename) && FileExists(config.EnrichmentProb.Dnase.Filename) {
      log.Fatal("EnrichmentProb bigWig files exist for both ATAC- and DNase-seq. Please select a single open chromatin assay.")
    }
  }
  if config.EnrichmentProb.Atac.Filename != "" && FileExists(config.EnrichmentProb.Atac.Filename) {
    return "atac"
  }
  if config.EnrichmentProb.Dnase.Filename != "" && FileExists(config.EnrichmentProb.Dnase.Filename) {
    return "dnase"
  }
  // return default assay
  return "atac"
}

func (config *ConfigModHmm) SetOpenChromatinAssay(assay string) {
  switch strings.ToLower(assay) {
  case "atac":
    config.Coverage       .Open = config.Coverage       .Atac
    config.CoverageCnts   .Open = config.CoverageCnts   .Atac
    config.EnrichmentModel.Open = config.EnrichmentModel.Atac
    config.EnrichmentComp .Open = config.EnrichmentComp .Atac
    config.EnrichmentPeak .Open = config.EnrichmentPeak .Atac
    config.EnrichmentProb .Open = config.EnrichmentProb .Atac
  case "dnase":
    config.Coverage       .Open = config.Coverage       .Dnase
    config.CoverageCnts   .Open = config.CoverageCnts   .Dnase
    config.EnrichmentModel.Open = config.EnrichmentModel.Dnase
    config.EnrichmentComp .Open = config.EnrichmentComp .Dnase
    config.EnrichmentPeak .Open = config.EnrichmentPeak .Dnase
    config.EnrichmentProb .Open = config.EnrichmentProb .Dnase
  default:
    log.Fatalf("invalid open chromatin assay `%s'", assay)
  }
  config.OpenChromatinAssay = assay
}

func (config *ConfigModHmm) CompletePaths(prefix string) {
  config.BamDir                 = config.setDefaultDir(prefix, config.BamDir               ,  "")
  config.CoverageDir            = config.setDefaultDir(prefix, config.CoverageDir          , config.BamDir)
  config.EnrichmentModelDir     = config.setDefaultDir(prefix, config.EnrichmentModelDir   , config.CoverageDir)
  config.EnrichmentDir          = config.setDefaultDir(prefix, config.EnrichmentDir        , config.EnrichmentModelDir)
  config.ChromatinStateDir      = config.setDefaultDir(prefix, config.ChromatinStateDir    , config.EnrichmentDir)
  config.ModelDir               = config.setDefaultDir(prefix, config.ModelDir             , config.ChromatinStateDir)
  config.SegmentationDir        = config.setDefaultDir(prefix, config.SegmentationDir      , config.ModelDir)
  config.PosteriorDir           = config.setDefaultDir(prefix, config.PosteriorDir         , config.SegmentationDir)
  config.Model.Filename         = completePath(config.ModelDir, "", config.Model.Filename, "segmentation.json")
  config.Segmentation.Filename  = completePath(config.SegmentationDir, "", config.Segmentation.Filename, "segmentation.bed.gz")
  config.Bam                    .CompletePaths(config.BamDir, "", "")
  config.Coverage               .CompletePaths(config.CoverageDir, "coverage-", ".bw")
  config.CoverageCnts           .CompletePaths(config.EnrichmentModelDir, "", ".counts.json")
  config.EnrichmentModel        .CompletePaths(config.EnrichmentModelDir, "", ".json")
  config.EnrichmentComp         .CompletePaths(config.EnrichmentModelDir, "", ".components.json")
  config.EnrichmentPeak         .CompletePaths(config.EnrichmentDir, "enrichment-peaks-", ".table")
  config.EnrichmentProb         .CompletePaths(config.EnrichmentDir, "enrichment-", ".bw")
  config.ChromatinStateProb     .CompletePaths(config.ChromatinStateDir, "chromatin-state-", ".bw")
  config.ChromatinStatePeak     .CompletePaths(config.ChromatinStateDir, "chromatin-state-", ".table")
  config.PosteriorProb          .CompletePaths(config.PosteriorDir, "posterior-marginal-", ".bw")
  config.PosteriorPeak          .CompletePaths(config.PosteriorDir, "posterior-marginal-peaks-", ".table")
  config.SetOpenChromatinAssay(config.DetectOpenChromatinAssay())
  if config.EnrichmentModelStatic {
    config.CoverageCnts   .SetStatic(true)
    config.EnrichmentModel.SetStatic(true)
    config.EnrichmentComp .SetStatic(true)
  }
}

func (config *ConfigModHmm) EnrichmentFiles(feature string) EnrichmentFiles {

  if !EnrichmentList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
  files := EnrichmentFiles{}
  files.Feature = config.CoerceOpenChromatinAssay(strings.ToLower(feature))

  switch files.Feature {
  case "rna-low":
    files.Probabilities = config.EnrichmentProb .Rna_low
    files.Model         = config.EnrichmentModel.Rna
    files.Components    = config.EnrichmentComp .Rna_low
    files.Coverage      = config.Coverage       .Rna
    files.CoverageCnts  = config.CoverageCnts   .Rna
  default:
    files.Probabilities = config.EnrichmentProb .GetTargetFile(feature)
    files.Model         = config.EnrichmentModel.GetTargetFile(feature)
    files.Components    = config.EnrichmentComp .GetTargetFile(feature)
    files.Coverage      = config.Coverage       .GetTargetFile(feature)
    files.CoverageCnts  = config.CoverageCnts   .GetTargetFile(feature)
  }
  return files
}

func (config ConfigModHmm) ModelFallbackPath() string {
  switch strings.ToLower(config.ModelFallback) {
  case "mm10"  :
    return "mm10-forebrain-embryo-day11.5"
  case "mm10-liver-embryo-day12.5":
    return "mm10-liver-embryo-day12.5"
  case "grch38":
    return "GRCh38-gastrocnemius-medialis"
  default:
    log.Fatalf("invalid single-feature model fallback `%s'", config.ModelFallback)
    panic("internal error")
  }
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

func (config ConfigEnrichmentPaths) String(openChromatinAssay string) string {
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
  fmt.Fprintf(&buffer, " -> RNA (low)            : %v\n", config.Rna_low)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigEnrichmentParameters) String() string {
  var buffer bytes.Buffer

  fmt.Fprintf(&buffer, " -> Open                 : %v\n", config.Open)
  fmt.Fprintf(&buffer, " -> H3K27ac              : %v\n", config.H3k27ac)
  fmt.Fprintf(&buffer, " -> H3K27me3             : %v\n", config.H3k27me3)
  fmt.Fprintf(&buffer, " -> H3K4me1              : %v\n", config.H3k4me1)
  fmt.Fprintf(&buffer, " -> H3K4me3              : %v\n", config.H3k4me3)
  fmt.Fprintf(&buffer, " -> RNA                  : %v\n", config.Rna)
  fmt.Fprintf(&buffer, " -> Control              : %v\n", config.Control)

  return buffer.String()
}

func (config ConfigChromatinStatePaths) String() string {
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
    fmt.Fprintf(&buffer, " -> Open Chromatin Assay   : %s\n"  , config.OpenChromatinAssay)
    fmt.Fprintf(&buffer, " -> Coverage Bin Size      : %d\n\n", config.CoverageBinSize)
    fmt.Fprintf(&buffer, "Alignment files (BAM):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Bam.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Coverage files (bigWig):\n")
    fmt.Fprintf(&buffer, "%v\n", config.Coverage.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 1 && config.EnrichmentMethod == "model" {
    fmt.Fprintf(&buffer, "Enrichment mixture distributions:\n")
    fmt.Fprintf(&buffer, "%v\n", config.EnrichmentModel.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Enrichment count statistics:\n")
    fmt.Fprintf(&buffer, "%v\n", config.CoverageCnts.String(config.OpenChromatinAssay))
    fmt.Fprintf(&buffer, "Enrichment mixture components:\n")
    fmt.Fprintf(&buffer, "%v\n", config.EnrichmentComp.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Enrichment probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.EnrichmentProb.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Enrichment peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.EnrichmentPeak.String(config.OpenChromatinAssay))
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Enrichment parameters:\n")
    fmt.Fprintf(&buffer, "%v\n", config.EnrichmentParameters.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Chromatin state probabilities:\n")
    fmt.Fprintf(&buffer, "%v\n", config.ChromatinStateProb.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Chromatin state peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.ChromatinStatePeak.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "Posterior marginals:\n")
    fmt.Fprintf(&buffer, "%v\n", config.PosteriorProb.String())
  }
  if config.Verbose > 1 {
    fmt.Fprintf(&buffer, "Posterior marginals peaks:\n")
    fmt.Fprintf(&buffer, "%v\n", config.PosteriorPeak.String())
  }
  if config.Verbose > 0 {
    fmt.Fprintf(&buffer, "ModHmm options:\n")
    fmt.Fprintf(&buffer, " ->  Description                  : %v\n", config.Description)
    fmt.Fprintf(&buffer, " ->  Directory                    : %v\n", config.Directory)
    fmt.Fprintf(&buffer, " ->  ModHmm Model File            : %v\n", config.Model)
    fmt.Fprintf(&buffer, " ->  ModHmm Model Directory       : %v\n", config.ModelDir)
    fmt.Fprintf(&buffer, " ->  ModHMM Model Fallback        : %s\n", config.ModelFallback)
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation File     : %v\n", config.Segmentation)
    fmt.Fprintf(&buffer, " ->  ModHmm Segmentation Directory: %v\n", config.SegmentationDir)
  }
  return buffer.String()
}
