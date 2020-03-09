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
import   "io"
import   "os"
import   "path"

import   "github.com/pborman/getopt"

import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

var Version   string
var BuildTime string
var GitHash   string

func printVersion(writer io.Writer) {
  fmt.Fprintf(writer, "ModHMM (https://github.com/pbenner/modhmm)\n")
  fmt.Fprintf(writer, " - Version   : %s\n", Version)
  fmt.Fprintf(writer, " - Build time: %s\n", BuildTime)
  fmt.Fprintf(writer, " - Git Hash  : %s\n", GitHash)
}

/* default model
 * -------------------------------------------------------------------------- */

//go:generate go run modhmm_default_gen.go

/* -------------------------------------------------------------------------- */

func main() {
  log.SetFlags(0)

  options := getopt.New()

  optConfig  := options. StringLong("config",  'c', "", "configuration file")
  optThreads := options.    IntLong("threads", 't',  1, "number of threads")
  optHelp    := options.   BoolLong("help",    'h',     "print help")
  optGenConf := options.   BoolLong("genconf",  0 ,     "print default config file")
  optVerbose := options.CounterLong("verbose", 'v',     "verbose level [-v or -vv]")
  optVersion := options.   BoolLong("version",  0 ,     "print ModHMM version")

  options.SetParameters("<COMMAND>\n\n" +
    " Default usage:\n" +
    "     modhmm -c CONFIG.json segmentation\n\n" +
    " Commands:\n" +
    "     coverage                 [stage 1]   - compute single-feature coverages from bam files\n" +
    "     estimate-single-feature  [stage 1.1] - estimate mixture distribution for single-feature\n" +
    "                                            enrichment analysis\n" +
    "     eval-single-feature      [stage 2]   - evaluate single-feature models\n" +
    "     eval-multi-feature       [stage 3]   - evaluate  multi-feature models\n" +
    "     segmentation             [stage 4]   - compute genome segmentation\n" +
    "     eval-posterior-marginals [stage 5]   - compute posterior marginals of hidden states\n\n" +
    " ModHMM commands are structured in stages. Executing a command also executes all commands with\n" +
    " lower stage number. An exception is stage 1.1 that must be executed manually. By default,\n" +
    " stage 1.1 is bypassed and a default single-feature model is used.\n\n" +
    " Plotting and printing commands:\n" +
    "     plot-single-feature                  - plot fitted single-feature mixture distributions used\n" +
    "                                            for enrichment analysis\n" +
    "     print-single-feature                 - print components and parameters of single-feature\n" +
    "                                            mixture distributions\n" +
    "     print-transition-matrix              - print estimated transition rates\n" +
    " Peak calling commands:\n" +
    "     call-single-feature-peaks            - call peaks of single-feature enrichment analysis\n" +
    "     call-multi-feature-peaks             - call peaks of multi-feature classifications\n" +
    "     call-posterior-marginal-peaks        - call peaks of HMM marginal posterior tracks\n")
  options.Parse(os.Args)

  config := DefaultModHmmConfig()

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if *optVersion {
    printVersion(os.Stdout)
    os.Exit(0)
  }
  if *optVerbose != 0 {
    config.Verbose = *optVerbose
  }
  if *optGenConf  {
    if err := config.Export(os.Stdout); err != nil {
      log.Fatal(err)
    }
    os.Exit(0)
  }
  if *optConfig != "" {
    current_config := config
    printStderr(current_config, 1, "Importing config file `%s'... ", *optConfig)
    if err := config.ImportFile(*optConfig); err != nil {
      printStderr(current_config, 1, "failed\n")
      log.Fatalf("reading config file `%s' failed: %v", *optConfig, err)
    }
    printStderr(current_config, 1, "done\n")
  }
  if *optThreads < 1 {
    log.Fatalf("invalid number of threads `%d'", *optThreads)
  }
  if options.Lookup('t').Seen() {
    config.Threads = *optThreads
  }
  // command arguments
  if len(options.Args()) == 0 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
  command := options.Args()[0]

  // print config
  config.CompletePaths(path.Dir(*optConfig))
  if str := config.String(); str != "" {
    printStderr(config, 0, "%s\n", str)
  }
  switch command {
  case "coverage":
    modhmm_coverage_main(config, options.Args())
  case "estimate-single-feature":
    modhmm_enrichment_estimate_main(config, options.Args())
  case "plot-single-feature":
    modhmm_enrichment_plot_main(config, options.Args())
  case "print-single-feature":
    modhmm_enrichment_print_main(config, options.Args())
  case "print-transition-matrix":
    modhmm_transition_matrix_print_main(config, options.Args())
  case "eval-single-feature":
    modhmm_enrichment_eval_main(config, options.Args())
  case "eval-multi-feature":
    modhmm_chromatin_state_eval_main(config, options.Args())
  case "eval-posterior-marginals":
    modhmm_posterior_main(config, options.Args())
  case "segmentation":
    modhmm_segmentation_main(config, options.Args())
  case "call-single-feature-peaks":
    modhmm_call_enrichment_peaks_main(config, options.Args())
  case "call-multi-feature-peaks":
    modhmm_call_chromatin_state_peaks_main(config, options.Args())
  case "call-posterior-marginal-peaks":
    modhmm_call_posterior_peaks_main(config, options.Args())
  default:
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
}
