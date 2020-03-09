/* Copyright (C) 2019 Philipp Benner
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
import   "os"
import   "strings"

import . "github.com/pbenner/autodiff/statistics"
import   "github.com/pbenner/autodiff/statistics/scalarDistribution"

import . "github.com/pbenner/modhmm/config"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_print_component(k int, pdf ScalarPdf) {
  switch a := pdf.(type) {
  case *scalarDistribution.DeltaDistribution:
    fmt.Printf(": %2d Delta     %e", k+1, a.X.GetValue())
  case *scalarDistribution.PoissonDistribution:
    fmt.Printf(": %2d Poisson   %e", k+1, a.GetParameters().ValueAt(0))
  case *scalarDistribution.GeometricDistribution:
    fmt.Printf(": %2d Geometric %e", k+1, a.GetParameters().ValueAt(0))
  case *scalarDistribution.PdfTranslation:
    modhmm_enrichment_print_component(k, a.ScalarPdf)
  default:
    log.Fatal("unknown distribution")
  }
}

func modhmm_enrichment_print_components(config ConfigModHmm, mixture *scalarDistribution.Mixture, k_fg []int) {
  fg := make([]bool, mixture.NComponents())
  for _, k := range k_fg {
    fg[k] = true
  }
  fmt.Println(":  # Type      Parameter")
  for k := 0; k < mixture.NComponents(); k ++ {
    modhmm_enrichment_print_component(k, mixture.Edist[k])
    if fg[k] {
      fmt.Printf(" [foreground]\n")
    } else {
      fmt.Printf(" [background]\n")
    }
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_print(config ConfigModHmm, feature string) {
  feature = config.CoerceOpenChromatinAssay(feature)

  if !EnrichmentList.Contains(strings.ToLower(feature)) {
    log.Fatalf("unknown feature: %s", feature)
  }

  files   := config.EnrichmentFiles(feature)
  mixture := ImportMixtureDistribution(config, files.Model.Filename)
  k, _    := ImportComponents(config, files.Components.Filename, mixture.NComponents())

  fmt.Printf("Mixture components for feature `%s'\n", feature)
  modhmm_enrichment_print_components(config, mixture, k)
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_print_loop(config ConfigModHmm, features []string) {
  for _, feature := range features {
    modhmm_enrichment_print(config, feature)
  }
}

func modhmm_enrichment_print_all(config ConfigModHmm) {
  modhmm_enrichment_print_loop(config, EnrichmentList)
}

/* -------------------------------------------------------------------------- */

func modhmm_enrichment_print_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s print-single-feature", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optHelp        := options.  BoolLong("help",  'h', "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if len(options.Args()) == 0 {
    modhmm_enrichment_print_all(config)
  } else {
    modhmm_enrichment_print_loop(config, options.Args())
  }
}
