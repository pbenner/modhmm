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
import   "log"
import   "math/rand"
import   "os"
import   "sort"
import   "strconv"
import   "strings"

import . "github.com/pbenner/ngstat/estimation"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"
import   "github.com/pbenner/autodiff/statistics/scalarEstimator"
import   "github.com/pbenner/autodiff/statistics/vectorDistribution"
import   "github.com/pbenner/autodiff/statistics/vectorEstimator"

import . "github.com/pbenner/gonetics"

import . "github.com/pbenner/modhmm/config"
import . "github.com/pbenner/modhmm/utility"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

type SortableMixture struct {
  *scalarDistribution.Mixture
}

func (obj SortableMixture) Len() int {
  return obj.NComponents()
}

func (obj SortableMixture) Less(i, j int) bool {
  xi, yi := distToValue(obj.Edist[i])
  xj, yj := distToValue(obj.Edist[j])
  if xi == xj {
    return yi < yj
  } else {
    return xi < xj
  }
}

func (obj SortableMixture) Swap(i, j int) {
  obj.Edist[i], obj.Edist[j] = obj.Edist[j], obj.Edist[i]
  obj.LogWeights.Swap(i, j)
}

func distToValue(dist ScalarPdf) (int, float64) {
  switch a := dist.(type) {
  case *scalarDistribution.DeltaDistribution:
    return 0, a.GetParameters().Float64At(0)
  case *scalarDistribution.PoissonDistribution:
    return 1, a.GetParameters().Float64At(0)
  case *scalarDistribution.PdfTranslation:
    return distToValue(a.ScalarPdf)
  case *scalarDistribution.GeometricDistribution:
    return 2, -a.GetParameters().Float64At(0)
  default:
    panic("internal error")
  }
}

/* -------------------------------------------------------------------------- */

func newEstimator(config ConfigModHmm, n_delta, n_poisson, n_geometric int) VectorEstimator {
  components := []ScalarEstimator{}
  for i := 0; i < n_delta; i++ {
    if delta, err := scalarEstimator.NewDeltaEstimator(float64(i)); err != nil {
      log.Fatal(err)
    } else {
      components = append(components, delta)
    }
  }
  for i := 0; i < n_poisson; i++ {
    if poisson, err := scalarEstimator.NewPoissonEstimator(rand.Float64()); err != nil {
      log.Fatal(err)
    } else {
      if t, err := scalarEstimator.NewTranslationEstimator(poisson, -float64(n_delta)); err != nil {
        log.Fatal(err)
      } else {
        components = append(components, t)
      }
    }
  }
  for i := 0; i < n_geometric; i++ {
    if geometric, err := scalarEstimator.NewGeometricEstimator(0.01*rand.Float64()); err != nil {
      log.Fatal(err)
    } else {
      components = append(components, geometric)
    }
  }
  if mixture, err := scalarEstimator.NewDiscreteMixtureEstimator(nil, components, 1e-8, -1); err != nil {
    log.Fatal(err)
  } else {
    if estimator, err := vectorEstimator.NewScalarIid(mixture, -1); err != nil {
      log.Fatal(err)
    } else {
      return estimator
    }
  }
  return nil
}

/* -------------------------------------------------------------------------- */

func enrichment_estimate(config ConfigModHmm, track Track, estimator VectorEstimator, files EnrichmentFiles) {
  if err := EstimateOnSingleTrack(config.SessionConfig, estimator, track); err != nil {
    log.Fatal(err)
  }
  if d, err := estimator.GetEstimate(); err != nil {
    log.Fatal(err)
  } else {
    result := d.(*vectorDistribution.ScalarIid).Distribution.(*scalarDistribution.Mixture)

    sort.Sort(SortableMixture{result})

    printStderr(config, 1, "Exporting distribution to `%s'... ", files.Model.Filename)
    if err := ExportDistribution(files.Model.Filename, result); err != nil {
      printStderr(config, 1, "failed\n")
      log.Fatal(err)
    }
    printStderr(config, 1, "done\n")
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_estimate(config ConfigModHmm, feature string, n []int, force bool) {
  var estimator VectorEstimator

  if !EnrichmentModelList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }
  files := config.EnrichmentFiles(feature)
  var track Track
  // update model
  if force || updateRequired(config, files.Model, files.DependenciesModel()...) {
    if track == nil {
      track = enrichment_import_model(config, files, false)
    }
    estimator = newEstimator(config, n[0], n[1], n[2])

    enrichment_estimate(config, track, estimator, files)
  }
  // update counts
  if force || updateRequired(config, files.CoverageCnts, files.DependenciesModel()...) {
    if track == nil {
      track = enrichment_import_model(config, files, false)
    }
    modhmm_compute_counts(config, track, files.CoverageCnts.Filename)
  }
}

func modhmm_enrichment_estimate_default(config ConfigModHmm, feature string, force bool, defcomp string) {
  var n, components []int
  switch strings.ToLower(defcomp) {
  case "mm10":
    switch strings.ToLower(feature) {
    case "open"    : fallthrough
    case "atac"    : fallthrough
    case "dnase"   : n = []int{1, 1, 3}; components = []int{3, 4}
    case "h3k27ac" : n = []int{1, 2, 2}; components = []int{4}
    case "h3k27me3": n = []int{4, 4, 1}; components = []int{8}
    case "h3k4me1" : n = []int{1, 8, 0}; components = []int{5, 6, 7, 8}
    case "h3k4me3" : n = []int{1, 1, 3}; components = []int{3, 4}
    case "h3k9me3" : n = []int{2, 4, 1}; components = []int{5, 6}
    case "rna"     : n = []int{1, 0, 4}; components = []int{2, 3, 4}
    case "control" : n = []int{7, 2, 1}; components = []int{9}
    }
  case "hg19":
    switch strings.ToLower(feature) {
    case "open"    : fallthrough
    case "atac"    : fallthrough
    case "dnase"   : n = []int{1, 1, 3}; components = []int{3, 4}
    case "h3k27ac" : n = []int{1, 2, 2}; components = []int{4}
    case "h3k27me3": n = []int{4, 4, 1}; components = []int{8}
    case "h3k4me1" : n = []int{1, 8, 0}; components = []int{4, 5, 6, 7, 8}
    case "h3k4me3" : n = []int{1, 2, 1}; components = []int{3}
    case "h3k9me3" : n = []int{2, 4, 1}; components = []int{5, 6}
    case "rna"     : n = []int{1, 0, 4}; components = []int{2, 3, 4}
    case "control" : n = []int{7, 2, 1}; components = []int{9}
    }
  default:
    log.Fatalf("unknown default comonents specifier: %s", defcomp)
  }
  // estimate mixture
  if EnrichmentModelList.Contains(strings.ToLower(feature)) {
    modhmm_enrichment_estimate(config, feature, n, force)
  }
  // export foreground mixture components
  files := config.EnrichmentFiles(feature)
  if force || updateRequired(config, files.Components, files.Model.Filename) {
    ExportComponents(config, files.Components.Filename, components)
  }
}

func modhmm_enrichment_estimate_default_loop(config ConfigModHmm, features []string, force bool, defcomp string) {
  // compute coverages here to make use of multi-threading
  modhmm_coverage_loop(config, InsensitiveStringList(features).Intersection(CoverageList))
  // eval single features
  for _, feature := range features {
    modhmm_enrichment_estimate_default(config, feature, force, defcomp)
  }
}

func modhmm_enrichment_estimate_default_all(config ConfigModHmm, force bool, defcomp string) {
  modhmm_enrichment_estimate_default_loop(config, EnrichmentList, force, defcomp)
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_estimate_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s estimate-single-feature", os.Args[0]))
  options.SetParameters("[<FEATURE> [<N_DELTA> <N_POISSON> <N_GEOMETRIC>]]\n")

  optDefComp := options. StringLong("default-components",  0 , "mm10", "default number of components [mm10, hg19]")
  optForce   := options.   BoolLong("force",               0 ,         "always overwrite existing files")
  optHelp    := options.   BoolLong("help",               'h',         "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  // command arguments
  if len(options.Args()) != 0 && len(options.Args()) != 1 && len(options.Args()) != 4 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
  var feature string
  if len(options.Args()) > 0 {
    feature = config.CoerceOpenChromatinAssay(options.Args()[0])
  }

  switch len(options.Args()) {
  case 4:
    n := []int{}
    if m, err := strconv.ParseInt(options.Args()[1], 10, 64); err != nil {
      log.Fatal(err)
    } else {
      n = append(n, int(m))
    }
    if m, err := strconv.ParseInt(options.Args()[2], 10, 64); err != nil {
      log.Fatal(err)
    } else {
      n = append(n, int(m))
    }
    if m, err := strconv.ParseInt(options.Args()[3], 10, 64); err != nil {
      log.Fatal(err)
    } else {
      n = append(n, int(m))
    }
    modhmm_enrichment_estimate(config, feature, n, *optForce)
  case 1:
    modhmm_enrichment_estimate_default(config, feature, *optForce, *optDefComp)
  case 0:
    modhmm_enrichment_estimate_default_all(config, *optForce, *optDefComp)
  }
}
