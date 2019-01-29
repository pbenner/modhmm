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

import   "github.com/pborman/getopt"

import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

type Config struct {
  SessionConfig
}

/* -------------------------------------------------------------------------- */

func printStderr(config Config, level int, format string, args ...interface{}) {
  if config.Verbose >= level {
    fmt.Fprintf(os.Stderr, format, args...)
  }
}

/* -------------------------------------------------------------------------- */

func loadModConfig(config Config, filename string, modconfig ConfigModHmm) ConfigModHmm {
  printStderr(config, 1, "Importing config file `%s'... ", filename)
  if err := modconfig.ImportFile(filename); err != nil {
    printStderr(config, 1, "failed\n")
    log.Fatalf("reading config file `%s' failed: %v", filename, err)
  }
  printStderr(config, 1, "done\n")
  return modconfig
}

/* -------------------------------------------------------------------------- */

func diffMarkovModel(config Config, config1, config2 ConfigModHmm) {
  posterior1 := make([]SimpleTrack, len(MultiFeatureList))
  posterior2 := make([]SimpleTrack, len(MultiFeatureList))

  for i, state := range MultiFeatureList {
    if t, err := ImportTrack(config.SessionConfig, config1.MultiFeatureProb.GetTargetFile(state).Filename); err != nil {
      log.Fatal(err)
    } else {
      posterior1[i] = t
    }
    if t, err := ImportTrack(config.SessionConfig, config2.MultiFeatureProb.GetTargetFile(state).Filename); err != nil {
      log.Fatal(err)
    } else {
      posterior2[i] = t
    }
  }

}

/* -------------------------------------------------------------------------- */

func main() {
  log.SetFlags(0)

  options := getopt.New()
  config  := Config{}

  optHelp    := options.   BoolLong("help",    'h',     "print help")
  optVerbose := options.CounterLong("verbose", 'v',     "verbose level [-v or -vv]")

  options.SetParameters("<CONFIG1.json> <CONFIG2.json>\n")
  options.Parse(os.Args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  if *optVerbose != 0 {
    config.Verbose = *optVerbose
  }
  // command arguments
  if len(options.Args()) != 2 {
    options.PrintUsage(os.Stderr)
    os.Exit(1)
  }
  config1 := DefaultModHmmConfig()
  config2 := DefaultModHmmConfig()
  config1  = loadModConfig(config, options.Args()[0], config1)
  config2  = loadModConfig(config, options.Args()[1], config2)

  config1.CompletePaths()
  config2.CompletePaths()

  diffMarkovModel(config, config1, config2)
}
