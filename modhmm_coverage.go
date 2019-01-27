/* Copyright (C) 2016-2018 Philipp Benner
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
import   "bufio"
import   "log"
import   "math"
import   "path/filepath"
import   "strconv"
import   "strings"
import   "os"

import . "github.com/pbenner/ngstat/track"

import   "github.com/pborman/getopt"
import . "github.com/pbenner/gonetics"
import   "github.com/pbenner/threadpool"
import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "gonum.org/v1/plot"
import   "gonum.org/v1/plot/plotter"
import   "gonum.org/v1/plot/plotutil"
import   "gonum.org/v1/plot/vg"

/* utility
 * -------------------------------------------------------------------------- */

func parseFilename(filename string) (string, int) {
  if tmp := strings.Split(filename, ":"); len(tmp) == 2 {
    t, err := strconv.ParseInt(tmp[1], 10, 64)
    if err != nil {
      log.Fatal(err)
    }
    return tmp[0], int(t)
  } else
  if len(tmp) >= 2 {
    log.Fatalf("invalid input file description `%s'", filename)
  }
  return filename, 0
}


/* fragment length estimation
 * -------------------------------------------------------------------------- */

func saveFraglen(config ConfigModHmm, feature, filename string, fraglen int) {
  basename := strings.TrimRight(filename, filepath.Ext(filename))
  filename  = fmt.Sprintf("%s.fraglen.txt", basename)

  f, err := os.Create(filename)
  if err != nil {
    log.Fatalf("[%s] ERROR: opening `%s' failed: %v", feature, filename, err)
  }
  defer f.Close()

  fmt.Fprintf(f, "%d\n", fraglen)

  printStderr(config, 1, "[%s] Wrote fragment length estimate to `%s'\n", feature, filename)
}

func saveCrossCorr(config ConfigModHmm, feature, filename string, x []int, y []float64) {
  basename := strings.TrimRight(filename, filepath.Ext(filename))
  filename  = fmt.Sprintf("%s.fraglen.table", basename)

  f, err := os.Create(filename)
  if err != nil {
    log.Fatalf("[%s] ERROR: opening `%s' failed: %v", feature, filename, err)
  }
  defer f.Close()

  for i := 0; i < len(x); i++ {
    fmt.Fprintf(f, "%d %f\n", x[i], y[i])
  }
  printStderr(config, 1, "[%s] Wrote crosscorrelation to `%s'\n", feature, filename)
}

func saveCrossCorrPlot(config ConfigModHmm, feature, filename string, fraglen int, x []int, y []float64) {
  basename := strings.TrimRight(filename, filepath.Ext(filename))
  filename  = fmt.Sprintf("%s.fraglen.pdf", basename)

  // draw cross-correlation
  xy := make(plotter.XYs, len(x))
  for i := 0; i < len(x); i++ {
    xy[i].X = float64(x[i])+1
    xy[i].Y = y[i]
  }
  p, err := plot.New()
  if err != nil {
    log.Fatal(err)
  }
  p.Title.Text = ""
  p.X.Label.Text = "shift"
  p.Y.Label.Text = "cross-correlation"

  err = plotutil.AddLines(p, xy)
  if err != nil {
    log.Fatal(err)
  }

  if fraglen != -1 {
    // determine cross-correlation maximum
    max_value := 0.0
    for i := 0; i < len(x); i++ {
      if y[i] > max_value {
        max_value = y[i]
      }
    }
    // draw vertical line at fraglen estimate
    fr := make(plotter.XYs, 2)
    fr[0].X = float64(fraglen)
    fr[0].Y = 0.0
    fr[1].X = float64(fraglen)
    fr[1].Y = max_value

    err = plotutil.AddLines(p, fr)
    if err != nil {
      log.Fatal(err)
    }
  }
  if err := p.Save(8*vg.Inch, 4*vg.Inch, filename); err != nil {
    log.Fatal(err)
  }
  printStderr(config, 1, "[%s] Wrote cross-correlation plot to `%s'\n", feature, filename)
}

func importFraglen(config ConfigModHmm, feature, filename string) int {
  // try reading the fragment length from file
  basename := strings.TrimRight(filename, filepath.Ext(filename))
  filename  = fmt.Sprintf("%s.fraglen.txt", basename)
  if f, err := os.Open(filename); err != nil {
    return 0
  } else {
    defer f.Close()
    printStderr(config, 1, "[%s] Reading fragment length from `%s'\n", feature, filename)
    scanner := bufio.NewScanner(f)
    if scanner.Scan() {
      if fraglen, err := strconv.ParseInt(scanner.Text(), 10, 64); err == nil {
        return int(fraglen)
      }
    }
    return 0
  }
}

/* -------------------------------------------------------------------------- */

func coverage_h3k4me3o1(config ConfigModHmm) error {
  config.BinSummaryStatistics = "discrete mean"
  config.BinOverlap = 1
  track1, err := ImportTrack(config.SessionConfig, config.Coverage.H3k4me1.Filename); if err != nil {
    return err
  }
  track2, err := ImportTrack(config.SessionConfig, config.Coverage.H3k4me3.Filename); if err != nil {
    return err
  }
  if err := (GenericMutableTrack{track1}).MapList([]Track{track1, track2}, func(seqname string, position int, values ...float64) float64 {
    x1 := values[0]
    x2 := values[1]
    return math.Round((x2+1.0)/(x1+1.0)*10)
  }); err != nil {
    return err
  }
  if err := ExportTrack(config.SessionConfig, track1, config.Coverage.H3k4me3o1.Filename); err != nil {
    return err
  }
  return err
}

/* -------------------------------------------------------------------------- */

func coverage(config ConfigModHmm, feature string, filenameBam []string, filenameData string, optionsList []interface{}) error {
  fraglen := make([]int, len(filenameBam))

  // split filename:fraglen
  for i, filename := range filenameBam {
    filenameBam[i], fraglen[i] = parseFilename(filename)
  }
  // import fragment length
  for i, filename := range filenameBam {
    if fraglen[i] == 0 {
      fraglen[i] = importFraglen(config, feature, filename)
    }
  }

  //////////////////////////////////////////////////////////////////////////////
  result, fraglenEstimate, _, err := BamCoverage(filenameData, filenameBam, nil, fraglen, nil, optionsList...)

  // save fraglen estimates
  //////////////////////////////////////////////////////////////////////////////
  for i, estimate := range fraglenEstimate {
    filename := filenameBam[i]
    if estimate.Error == nil && estimate.Fraglen != 0 {
      saveFraglen(config, feature, filename, estimate.Fraglen)
    }
    if estimate.X != nil && estimate.Y != nil {
      saveCrossCorr(config, feature, filename, estimate.X, estimate.Y)
    }
    if estimate.X != nil && estimate.Y != nil {
      saveCrossCorrPlot(config, feature, filename, estimate.Fraglen, estimate.X, estimate.Y)
    }
  }

  // process result
  //////////////////////////////////////////////////////////////////////////////
  if err != nil {
    return err
  } else {
    configLocal := config
    configLocal.Verbose = 0
    printStderr(config, 1, "Attempting to write track `%s'\n", filenameData)
    if err := ExportTrack(configLocal.SessionConfig, result, filenameData); err != nil {
      return fmt.Errorf("Writing track `%s' failed: %v", filenameData, err)
    }
    printStderr(config, 1, "Wrote track `%s'\n", filenameData)
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func modhmm_coverage(config ConfigModHmm, feature string) error {

  if !CoverageList.Contains(strings.ToLower(feature)) {
    return fmt.Errorf("unknown feature: %s", feature)
  }

  filenameBam  := []string{}
  filenameData := TargetFile{}
  optionsList  := []interface{}{}
  logPrefix    := feature

  switch strings.ToLower(feature) {
  case "open":
    switch strings.ToLower(config.OpenChromatinAssay) {
    case "atac":
      filenameBam  = config.Bam.Atac
      filenameData = config.Coverage.Atac
      // ATAC-seq typically uses paired-end sequencing;
      // for obtain only open-chromatin information we drop paired-end information
      optionsList  = append(optionsList, OptionPairedAsSingleEnd{true})
      optionsList  = append(optionsList, OptionFilterChroms{[]string{"chrM","M"}})
      logPrefix    = "atac"
    case "dnase":
      filenameBam  = config.Bam.Dnase
      filenameData = config.Coverage.Dnase
      // DNase-seq can be single- or paired-end;
      // for single-end sequencing no fragment-length estimation is performed
      optionsList  = append(optionsList, OptionFilterChroms{[]string{"chrM","M"}})
      logPrefix    = "dnase"
    default:
      panic("internal error")
    }
  case "rna":
    filenameBam  = config.Bam.Rna
    filenameData = config.Coverage.Rna
  case "h3k4me3o1":
    if updateRequired(config, config.Coverage.H3k4me3o1, config.Coverage.H3k4me1.Filename, config.Coverage.H3k4me3.Filename) {
      return coverage_h3k4me3o1(config)
    } else {
      return nil
    }
  default:
    filenameBam  = config.Bam     .GetFilenames (feature)
    filenameData = config.Coverage.GetTargetFile(feature)
    optionsList = append(optionsList, OptionEstimateFraglen{true})
    optionsList = append(optionsList, OptionFraglenRange{[2]int{100,300}})
    optionsList = append(optionsList, OptionFraglenBinSize{10})
  }
  if config.Verbose > 0 {
    optionsList = append(optionsList, OptionLogger{log.New(os.Stderr, fmt.Sprintf("[%s] ", logPrefix), 0)})
  }
  optionsList = append(optionsList, OptionBinningMethod{"mean overlap"})
  optionsList = append(optionsList, OptionBinSize{config.CoverageBinSize})
  optionsList = append(optionsList, OptionFilterMapQ{config.CoverageMAPQ})
  optionsList = append(optionsList, OptionFilterDuplicates{true})

  if updateRequired(config, filenameData, filenameBam...) {
    if len(filenameBam) == 0 {
      return fmt.Errorf("ERROR: no bam files specified for feature `%s'", logPrefix)
    }
    return coverage(config, feature, filenameBam, filenameData.Filename, optionsList)
  }
  return nil
}

func modhmm_coverage_loop(config ConfigModHmm, features []string) {
  pool := threadpool.New(config.CoverageThreads, 10)
  for _, feature := range features {
    if strings.ToLower(feature) == "h3k4me3o1" {
      continue
    }
    f := config.CoerceOpenChromatinAssay(feature)
    pool.AddJob(0, func(pool threadpool.ThreadPool, erf func() error) error {
      return modhmm_coverage(config, f)
    })
  }
  if err := pool.Wait(0); err != nil {
    log.Fatal(err)
  }
  if InsensitiveStringList(features).Contains("h3k4me3o1") {
    if err := modhmm_coverage(config, "h3k4me3o1"); err != nil {
      log.Fatal(err)
    }
  }
}

func modhmm_coverage_all(config ConfigModHmm) {
  modhmm_coverage_loop(config, CoverageList)
}

/* -------------------------------------------------------------------------- */

func modhmm_coverage_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s coverage", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optHelp := options.   BoolLong("help",     'h',     "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_coverage_all(config)
  } else {
    modhmm_coverage_loop(config, options.Args())
  }
}
