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

//import   "fmt"
import   "log"
import   "os"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func main() {
  log.SetFlags(0)

  options := getopt.New()

  optConfig  := options. StringLong("config",  'c', "", "configuration file")
  optThreads := options.    IntLong("threads", 't',  1, "number of threads")
  optHelp    := options.   BoolLong("help",    'h',     "print help")
  optVerbose := options.CounterLong("verbose", 'v',     "verbose level [-v or -vv]")

  options.SetParameters("<COMMAND>\n\n" +
    " Commands:\n" +
    "     coverage                 [stage 1]   - compute single-feature coverages from bam files\n" +
    "     compute-counts           [stage 1.1] - compute coverage counts used for quantile normalization\n" +
    "     estimate-single-feature  [stage 1.2] - estimate mixture distribution for single-feature\n" +
    "                                            enrichment analysis\n" +
    "     eval-single-feature      [stage 2]   - call enriched regions of single feature coverages\n" +
    "     eval-multi-feature       [stage 3]   - evaluate multi-feature models\n" +
    "     eval-posterior-marginals [stage 5]   - compute posterior marginals of hidden states\n\n" +
    "     segmentation             [stage 4]   - compute genome segmentation\n\n" +
    " ModHMM commands are structured in stages. Executing a command also executes all commands with\n" +
    " lower stage number. An exception are stages 1.1 and 1.2 that must be executed manually or\n" +
    " bypassed by providing single-feature models.\n\n")
  options.Parse(os.Args)

  config := DefaultModHmmConfig()

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if *optVerbose != 0 {
    config.Verbose = *optVerbose
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
  config.CompletePaths()
  printStderr(config, 0, "%v\n", config)

  switch command {
  case "coverage":
    modhmm_coverage_main(config, options.Args())
  case "compute-counts":
    modhmm_compute_counts_main(config, options.Args())
  case "estimate-single-feature":
    modhmm_single_feature_estimate_main(config, options.Args())
  case "eval-single-feature":
    modhmm_single_feature_eval_main(config, options.Args())
  case "eval-multi-feature":
    modhmm_multi_feature_eval_main(config, options.Args())
  case "segmentation":
    modhmm_segmentation_main(config, options.Args())
  case "eval-posterior-marginals":
    modhmm_posterior_main(config, options.Args())
  default:
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
}
