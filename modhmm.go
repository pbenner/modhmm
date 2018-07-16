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
    "     coverage                 - compute single-feature coverages from bam files\n" +
    "     estimate-single-feature  - estimate mixture distribution for single-\n" +
    "                                feature enrichment analysis\n" +
    "     eval-single-feature      - call enriched regions of single feature coverages\n" +
    "     eval-multi-feature       - evaluate multi-feature models\n" +
    "     eval-multi-feature-norm  - evaluate multi-feature models where the results are\n" +
    "                                normalized for visual inspection in the genome browser\n" +
    "     posterior-marginals      - compute posterior marginals for the hidden states\n" +
    "     segmentation             - compute genome segmentation\n")
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
  printStderr(config, 1, "%v\n", config)

  switch command {
  case "single-feature-coverage":
    modhmm_coverage_main(config, options.Args())
  case "estimate-single-feature":
    modhmm_single_feature_estimate_main(config, options.Args())
  case "eval-single-feature":
    modhmm_single_feature_eval_main(config, options.Args())
  case "eval-multi-feature":
    modhmm_multi_feature_eval_main(config, options.Args())
  case "eval-multi-feature-norm":
    modhmm_multi_feature_eval_norm_main(config, options.Args())
  case "segmentation":
    modhmm_segmentation_main(config, options.Args())
  case "posterior-marginals":
    modhmm_posterior_main(config, options.Args())
  default:
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
}
