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
import   "image/color"
import   "io/ioutil"
import   "log"
import   "math"
import   "os"
import   "os/exec"
import   "strconv"
import   "strings"

import . "github.com/pbenner/autodiff"
import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

import   "gonum.org/v1/plot"
import   "gonum.org/v1/plot/plotter"
import   "gonum.org/v1/plot/plotutil"
import   "gonum.org/v1/plot/vg"

/* -------------------------------------------------------------------------- */

func eval_counts(counts Counts, xlim [2]float64) (plotter.XYs, float64) {
  xy  := make(plotter.XYs, 0)
  sum := 0
  for i := 0; i < len(counts.X); i++ {
    sum += counts.Y[i]
  }
  y_min := math.Inf(1)
  for i := 0; i < len(counts.X); i++ {
    if !math.IsNaN(xlim[0]) && counts.X[i] < xlim[0] {
      continue
    }
    if !math.IsNaN(xlim[1]) && counts.X[i] > xlim[1] {
      continue
    }
    y := float64(counts.Y[i])/float64(sum)
    xy = append(xy, plotter.XY{counts.X[i], y})
    if y < y_min {
      y_min = y
    }
  }
  return xy, y_min
}

func eval_component(mixture *scalarDistribution.Mixture, k_ []int, counts Counts, xlim [2]float64, y_min float64) plotter.XYs {
  r  := NullBareReal()
  xy := make(plotter.XYs, 0)
  for i := 0; i < len(counts.X); i++ {
    if !math.IsNaN(xlim[0]) && counts.X[i] < xlim[0] {
      continue
    }
    if !math.IsNaN(xlim[1]) && counts.X[i] > xlim[1] {
      continue
    }
    y := 0.0
    for _, k := range k_ {
      if err := mixture.Edist[k].LogPdf(r, ConstReal(counts.X[i])); err != nil {
        log.Fatal("evaluating mixture component failed:", err)
      } else {
        y += math.Exp(mixture.LogWeights.ValueAt(k) + r.GetValue())
      }
    }
    if math.IsInf(y, 0) || y == 0.0 || y < y_min {
      continue
    }
    xy = append(xy, plotter.XY{counts.X[i], y})
  }
  return xy
}

func eval_delta_component(mixture *scalarDistribution.Mixture, k int, xlim [2]float64, y_min float64) plotter.XYs {
  x := mixture.Edist[k].(*scalarDistribution.DeltaDistribution).X.GetValue()
  y := math.Exp(mixture.LogWeights.ValueAt(k))
  if !math.IsNaN(xlim[0]) && x < xlim[0] {
    return nil
  }
  if !math.IsNaN(xlim[1]) && x > xlim[1] {
    return nil
  }
  if math.IsInf(y, 0) || y == 0.0 || y < y_min {
    return nil
  }
  xy := plotter.XY{}
  xy.X = x
  xy.Y = y
  return plotter.XYs([]plotter.XY{xy})
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot_isolated(config ConfigModHmm, p *plot.Plot, mixture *scalarDistribution.Mixture, counts Counts, xlim [2]float64) {
  counts_xy, y_min := eval_counts(counts, xlim)
  plotutil.DefaultColors = []color.Color{color.RGBA{0, 0, 0, 255}}
  if err := plotutil.AddLines(p, "Raw", counts_xy); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
  var list_points []interface{}
  var list_lines  []interface{}
  for k := 0; k < mixture.NComponents(); k ++ {
    switch mixture.Edist[k].(type) {
    case *scalarDistribution.DeltaDistribution:
      xys := eval_delta_component(mixture, k, xlim, y_min)
      list_points = append(list_points, fmt.Sprintf("Component %d", k+1))
      list_points = append(list_points, xys)
    default:
      xys := eval_component(mixture, []int{k}, counts, xlim, y_min)
      list_lines = append(list_lines, fmt.Sprintf("Component %d", k+1))
      list_lines = append(list_lines, xys)
    }
  }
  plotutil.DefaultColors = plotutil.SoftColors
  if err := plotutil.AddScatters(p, list_points...); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
  if err := plotutil.AddLines(p, list_lines...); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
}

func modhmm_single_feature_plot_joined(config ConfigModHmm, p *plot.Plot, mixture *scalarDistribution.Mixture, counts Counts, k_fg, k_bg []int, xlim [2]float64) {
  counts_xy, y_min := eval_counts(counts, xlim)
  plotutil.DefaultColors = []color.Color{color.RGBA{0, 0, 0, 255}}
  if err := plotutil.AddLines(p, "Raw", counts_xy); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
  xys_fg := eval_component(mixture, k_fg, counts, xlim, y_min)
  xys_bg := eval_component(mixture, k_bg, counts, xlim, y_min)
  plotutil.DefaultColors = plotutil.SoftColors
  if err := plotutil.AddLines(p, "Foreground", xys_fg, "Background", xys_bg); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot(config ConfigModHmm, feature string, xlim [2]float64, save string) {
  mixture := &scalarDistribution.Mixture{}
  counts  := Counts{}

  filenameModel, filenameComp, _, filenameCnts, _, _ := single_feature_files(config, feature, false)

  if !FileExists(filenameModel.Filename) {
    log.Fatalf("%s: file does not exist", filenameModel.Filename)
  }
  if !FileExists(filenameCnts.Filename) {
    log.Fatalf("%s: file does not exist", filenameCnts.Filename)
  }
  printStderr(config, 1, "Importing mixture model from `%s'... ", filenameModel.Filename)
  if err := ImportDistribution(filenameModel.Filename, mixture, BareRealType); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  printStderr(config, 1, "Importing reference counts from `%s'... ", filenameCnts.Filename)
  if err := counts.ImportFile(filenameCnts.Filename); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatal(err)
  }
  printStderr(config, 1, "done\n")

  p, err := plot.New()
  if err != nil {
    log.Fatal(err)
  }
  p.Title.Text   = feature
  p.X.Label.Text = "coverage value"
  p.Y.Label.Text = "probability"
  p.X.Scale = plot.LinearScale{}
  p.Y.Scale = plot.LogScale{}
  p.Y.Tick.Marker = plot.LogTicks{}
  p.Legend.Top = true

  if !FileExists(filenameComp.Filename) {
    modhmm_single_feature_plot_isolated(config, p, mixture, counts, xlim)
  } else {
    k_fg := ImportComponents(config, filenameComp.Filename, mixture.NComponents())
    k_bg := Components(k_fg).Invert(mixture.NComponents())

    modhmm_single_feature_plot_joined(config, p, mixture, counts, k_fg, k_bg, xlim)
  }
  size_x := 10*vg.Inch
  size_y :=  6*vg.Inch
  if save == "" {
    tmpfile, err := ioutil.TempFile("", "*.png")
    if err != nil {
      log.Fatal(err)
    }
    tmpfile.Close()

    if err := p.Save(size_x, size_y, tmpfile.Name()); err != nil {
      os.Remove(tmpfile.Name())
      log.Fatal(err)
    }
    cmd := exec.Command("display", tmpfile.Name())
    if err := cmd.Run(); err != nil {
      os.Remove(tmpfile.Name())
      log.Fatalf("%v: opening image viewer failed - try using `--save`", err)
    }
    os.Remove(tmpfile.Name())
  } else {
    if err := p.Save(size_x, size_y, save); err != nil {
      log.Fatal(err)
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s plot-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]\n")

  optSave := options.StringLong("save",  0 , "", "save plot to file")
  optXlim := options.StringLong("xlim",  0 , "", "range of the x-axis (e.g. 0-100)")
  optHelp := options.  BoolLong("help", 'h',     "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) != 1 {
    options.PrintUsage(os.Stdout)
    os.Exit(1)
  }
  xlim := [2]float64{math.NaN(), math.NaN()}
  if *optXlim != "" {
    r := strings.Split(*optXlim, "-")
    if len(r) != 2 {
      options.PrintUsage(os.Stdout)
      os.Exit(1)
    }
    if v, err := strconv.ParseFloat(r[0], 64); err != nil {
      log.Fatal(err)
    } else {
      xlim[0] = v
    }
    if v, err := strconv.ParseFloat(r[1], 64); err != nil {
      log.Fatal(err)
    } else {
      xlim[1] = v
    }
  }
  modhmm_single_feature_plot(config, options.Args()[0], xlim, *optSave)
}
