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
import   "io"
import   "math"
import   "os"
import   "strconv"
import   "strings"

import   "github.com/pborman/getopt"

import . "github.com/pbenner/gonetics"
import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/ngstat/track"
import . "github.com/pbenner/ngstat/utility"
import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

type Config struct {
  SessionConfig
  Threshold float64
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

func printResult(config Config, writer io.Writer, result [][]int, invertNames map[string]int) {
  names := make([]string, len(invertNames))
  for name, i := range invertNames {
    names[i] = name
  }
  // print header
  fmt.Fprintf(writer, "%7s ", "")
  for j := 0; j < len(names); j++ {
    fmt.Fprintf(writer, "%7s ", names[j])
  }
  fmt.Fprintf(writer, "\n")
  // print result
  for i := 0; i < len(names); i++ {
    fmt.Fprintf(writer, "%7s ", names[i])
    for j := 0; j < len(names); j++ {
      fmt.Fprintf(writer, "%7d ", result[i][j])
    }
    fmt.Fprintf(writer, "\n")
  }
}

/* -------------------------------------------------------------------------- */

func diffMarkovModel(config Config, config1, config2 ConfigModHmm) {
  genome, err := BigWigImportGenome(config1.MultiFeatureProb.GetTargetFile(MultiFeatureList[0]).Filename); if err != nil {
    log.Fatal(err)
  }
  seg1 := AllocSimpleTrack("", genome, config1.BinSize)
  seg2 := AllocSimpleTrack("", genome, config2.BinSize)

  stateNames1, err := (GenericMutableTrack{seg1}).ImportSegmentation(config1.Segmentation.Filename)
  if err != nil {
    log.Fatal(err)
  }
  stateNames2, err := (GenericMutableTrack{seg2}).ImportSegmentation(config2.Segmentation.Filename)
  if err != nil {
    log.Fatal(err)
  }
  for i := 0; i < len(stateNames1); i++ {
    // convert state names, i.e. EA:tr -> EA
    stateNames1[i] = strings.Split(stateNames1[i], ":")[0]
  }
  for i := 0; i < len(stateNames2); i++ {
    stateNames2[i] = strings.Split(stateNames2[i], ":")[0]
  }
  invertNames := make(map[string]int)
  for _, name := range stateNames1 {
    if _, ok := invertNames[name]; !ok {
      invertNames[name] = len(invertNames)
    }
  }
  result := make([][]int, len(invertNames))
  for i := 0; i < len(invertNames); i++ {
    result[i] = make([]int, len(invertNames))
  }

  posterior1 := make([]SimpleTrack, len(invertNames))
  posterior2 := make([]SimpleTrack, len(invertNames))

  for i, state := range stateNames1 {
    if t, err := ImportTrack(config1.SessionConfig, config1.MultiFeatureProb.GetTargetFile(state).Filename); err != nil {
      log.Fatal(err)
    } else {
      posterior1[i] = t
    }
  }
  for i, state := range stateNames2 {
    if t, err := ImportTrack(config2.SessionConfig, config2.MultiFeatureProb.GetTargetFile(state).Filename); err != nil {
      log.Fatal(err)
    } else {
      posterior2[i] = t
    }
  }

  // counter
  l := 0
  // total track length
  L := 0
  for _, length := range genome.Lengths {
    L += length/config1.BinSize
  }
  if config.Verbose > 0 {
    NewProgress(L, L).PrintStderr(l)
  }
  for _, name := range seg1.GetSeqNames() {
    seq1, err := seg1.GetSequence(name); if err != nil {
      log.Fatal(err)
    }
    seq2, err := seg2.GetSequence(name); if err != nil {
      log.Fatal(err)
    }
    for i := 0; i < seq1.NBins(); i++ {
      s1 := int(seq1.AtBin(i))
      s2 := int(seq2.AtBin(i))
      r1 := invertNames[stateNames1[s1]]
      r2 := invertNames[stateNames2[s2]]
      if stateNames1[s1] != stateNames2[s2] {
        t1, err := posterior1[s1].GetSequence(name); if err != nil {
          log.Fatal(err)
        }
        t2, err := posterior2[s2].GetSequence(name); if err != nil {
          log.Fatal(err)
        }
        p1 := t1.AtBin(i)
        p2 := t2.AtBin(i)
        if math.Exp(p1) > config.Threshold && math.Exp(p2) > config.Threshold {
          result[r1][r2]++
        } else {
          result[r1][r1]++
        }
      } else {
        result[r1][r2]++
      }
    }
    l += seq1.NBins()

    if config.Verbose > 0 {
      NewProgress(L, L).PrintStderr(l)
    }
  }
  printResult(config, os.Stdout, result, invertNames)
}

/* -------------------------------------------------------------------------- */

func main() {
  log.SetFlags(0)

  options := getopt.New()
  config  := Config{}

  optThreshold := options. StringLong("threshold",  0 , "0.8", "threshold")
  optHelp      := options.   BoolLong("help",      'h',        "print help")
  optVerbose   := options.CounterLong("verbose",   'v',        "verbose level [-v or -vv]")

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
  if t, err := strconv.ParseFloat(*optThreshold, 64); err != nil {
    log.Fatal(err)
  } else {
    config.Threshold = t
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
