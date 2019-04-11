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
import   "io"
import   "io/ioutil"
import   "log"
import   "math"
import   "os"
import   "os/exec"
import   "path"
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
import   "gonum.org/v1/plot/vg/draw"
import   "gonum.org/v1/plot/vg/vgimg"
import   "gonum.org/v1/plot/vg/vgpdf"

/* -------------------------------------------------------------------------- */

type scientificLogTicks struct{}

func (scientificLogTicks) Ticks(min, max float64) []plot.Tick {
  tks := plot.LogTicks{}.Ticks(min, max)
  for i, t := range tks {
    if t.Label == "" { // Skip minor ticks, they are fine.
      continue
    }
    if t, err := strconv.ParseFloat(t.Label, 64); err != nil {
      log.Fatal(err)
    } else {
      tks[i].Label = fmt.Sprintf("%.1e", t)
    }
  }
  return tks
}

/* -------------------------------------------------------------------------- */

func nrc(n int) (int, int) {
  r    := math.Sqrt(float64(n))
  x, y := 0.0, 0.0
  if r - math.Floor(r) >= 0.5 {
    x, y = math.Ceil (r), math.Ceil(r)
  } else {
    x, y = math.Floor(r), math.Ceil(r)
  }
  return int(x), int(y)
}

/* -------------------------------------------------------------------------- */

func plot_result(plots [][]*plot.Plot, save string) error {
  n1 := len(plots)
  n2 := len(plots[0])
  s1 := vg.Points(float64(n2)*300)
  s2 := vg.Points(float64(n1)*200)
  t  := draw.Tiles{
    Rows:      n1,
    Cols:      n2,
    PadX:      vg.Millimeter,
    PadY:      vg.Millimeter,
    PadTop:    vg.Points(2),
    PadBottom: vg.Points(2),
    PadLeft:   vg.Points(2),
    PadRight:  vg.Points(2),
  }
  var img          vg.CanvasSizer
  var canvases [][]draw.Canvas
  var writer       io.Writer
  var filename     string

  switch strings.ToLower(path.Ext(save)) {
  default    : fallthrough
  case ".png":
    img      = vgimg.New(s1, s2)
    dc      := draw .New(img)
    canvases = plot .Align(plots, t, dc)
  case ".pdf":
    img      = vgpdf.New(s1, s2)
    dc      := draw .New(img)
    canvases = plot .Align(plots, t, dc)
  }
  for i := 0; i < n1; i++ {
    for j := 0; j < n2; j++ {
        if plots[i][j] != nil {
            plots[i][j].Draw(canvases[i][j])
        }
    }
  }
  if save == "" {
    if tmpfile, err := ioutil.TempFile("", "*.png"); err != nil {
      return err
    } else {
      defer os.Remove(tmpfile.Name())
      defer tmpfile.Close()
      writer   = tmpfile
      filename = tmpfile.Name()
    }
  } else {
    if f, err := os.Create(save); err != nil {
      return err
    } else {
      defer f.Close()
      writer = f
    }
  }
  switch a := img.(type) {
  case *vgimg.Canvas:
    png := vgimg.PngCanvas{Canvas: a}
    if _, err := png.WriteTo(writer); err != nil {
      return err
    }
  case *vgpdf.Canvas:
    if _, err := a.WriteTo(writer); err != nil {
      return err
    }
  }
  if filename != "" {
    cmd := exec.Command("display", filename)
    if err := cmd.Run(); err != nil {
      return fmt.Errorf("%v: opening image viewer failed - try using `--save`", err)
    }
  }
  return nil
}

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

func modhmm_single_feature_plot_isolated(config ConfigModHmm, p *plot.Plot, mixture *scalarDistribution.Mixture, counts Counts) {
  counts_xy, y_min := eval_counts(counts, config.XLim)
  plotutil.DefaultColors = []color.Color{color.RGBA{0, 0, 0, 255}}
  if err := plotutil.AddLines(p, "raw", counts_xy); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
  var list_points []interface{}
  var list_lines  []interface{}
  for k := 0; k < mixture.NComponents(); k ++ {
    switch mixture.Edist[k].(type) {
    case *scalarDistribution.DeltaDistribution:
      xys := eval_delta_component(mixture, k, config.XLim, y_min)
      list_points = append(list_points, fmt.Sprintf("component %d", k+1))
      list_points = append(list_points, xys)
    default:
      xys := eval_component(mixture, []int{k}, counts, config.XLim, y_min)
      list_lines = append(list_lines, fmt.Sprintf("component %d", k+1))
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

func modhmm_single_feature_plot_joined(config ConfigModHmm, p *plot.Plot, mixture *scalarDistribution.Mixture, counts Counts, k_fg, k_bg []int) {
  counts_xy, y_min := eval_counts(counts, config.XLim)
  plotutil.DefaultColors = []color.Color{color.RGBA{0, 0, 0, 255}}
  if err := plotutil.AddLines(p, "raw", counts_xy); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
  xys_fg := eval_component(mixture, k_fg, counts, config.XLim, y_min)
  xys_bg := eval_component(mixture, k_bg, counts, config.XLim, y_min)
  plotutil.DefaultColors = plotutil.SoftColors
  if err := plotutil.AddLines(p, "foreground", xys_fg, "background", xys_bg); err != nil {
    log.Fatal("plotting mixture distribution failed: ", err)
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot(config ConfigModHmm, feature string) *plot.Plot {
  if !SingleFeatureList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
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
  p.Y.Tick.Marker = scientificLogTicks{}
  // set font size
  p.Title .Font.Size       = vg.Length(config.FontSize)
  p.Legend.Font.Size       = vg.Length(config.FontSize)
  p.X.Label.Font.Size      = vg.Length(config.FontSize)
  p.Y.Label.Font.Size      = vg.Length(config.FontSize)
  p.X.Tick.Label.Font.Size = vg.Length(config.FontSize)
  p.Y.Tick.Label.Font.Size = vg.Length(config.FontSize)

  p.Legend.Top = true

  if !FileExists(filenameComp.Filename) {
    modhmm_single_feature_plot_isolated(config, p, mixture, counts)
  } else {
    k_fg := ImportComponents(config, filenameComp.Filename, mixture.NComponents())
    k_bg := Components(k_fg).Invert(mixture.NComponents())

    modhmm_single_feature_plot_joined(config, p, mixture, counts, k_fg, k_bg)
  }
  return p
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot_loop(config ConfigModHmm, save string, features []string) {
  n1, n2 := nrc(len(features))
  plots := make([][]*plot.Plot, n1)
  for i := 0; i < n1; i++ {
    plots[i] = make([]*plot.Plot, n2)
    for j := 0; j < n2; j++ {
      if i*n2+j >= len(features) {
        break
      }
      plots[i][j] = modhmm_single_feature_plot(config, features[i*n2+j])
    }
  }
  if err := plot_result(plots, save); err != nil {
    log.Fatal(err)
  }
}

func modhmm_single_feature_plot_all(config ConfigModHmm, save string) {
  modhmm_single_feature_plot_loop(config, save, SingleFeatureList)
}

/* -------------------------------------------------------------------------- */

func modhmm_single_feature_plot_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s plot-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optSave     := options.StringLong("save",       0 , "", "save plot to file")
  optXlim     := options.StringLong("xlim",       0 , "", "range of the x-axis (e.g. 0-100)")
  optFontSize := options.StringLong("font-size",  0 , "", "size of the font")
  optHelp     := options.  BoolLong("help",      'h',     "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if *optXlim != "" {
    r := strings.Split(*optXlim, "-")
    if len(r) != 2 {
      options.PrintUsage(os.Stdout)
      os.Exit(1)
    }
    if v, err := strconv.ParseFloat(r[0], 64); err != nil {
      log.Fatal(err)
    } else {
      config.XLim[0] = v
    }
    if v, err := strconv.ParseFloat(r[1], 64); err != nil {
      log.Fatal(err)
    } else {
      config.XLim[1] = v
    }
  }
  if *optFontSize != "" {
    if v, err := strconv.ParseFloat(*optFontSize, 64); err != nil {
      log.Fatal(err)
    } else {
      config.FontSize = v
    }
  }
  if len(options.Args()) == 0 {
    modhmm_single_feature_plot_all(config, *optSave)
  } else {
    modhmm_single_feature_plot_loop(config, *optSave, options.Args())
  }
}
