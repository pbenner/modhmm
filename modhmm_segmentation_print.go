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
import   "math"
import   "os"

import . "github.com/pbenner/modhmm/config"

import   "github.com/pborman/getopt"

/* -------------------------------------------------------------------------- */

func modhmm_transition_matrix_print(config ConfigModHmm) {
  modhmm := ImportHMM(config)

  tr := modhmm.Tr
  sn := modhmm.StateNames

  n, _ := tr.Dims()

  // header
  fmt.Printf("%5s", "")
  for i := 0; i < n; i++ {
    fmt.Printf(" %8s", sn[i])
  }
  fmt.Println()
  for i := 0; i < n; i++ {
    fmt.Printf("%5s", sn[i])
    for j := 0; j < n; j++ {
      fmt.Printf(" %8.2e", math.Exp(tr.ValueAt(i,j)))
    }
    fmt.Println()
  }
}

/* -------------------------------------------------------------------------- */

func modhmm_transition_matrix_print_main(config ConfigModHmm, args []string) {

  options := getopt.New()
  options.SetProgram(fmt.Sprintf("%s print-transition-matrix", os.Args[0]))
  options.SetParameters("[FEATURE]...\n")

  optHelp        := options.  BoolLong("help",  'h', "print help")

  options.Parse(args)

  // command options
  if *optHelp {
    options.PrintUsage(os.Stdout)
    os.Exit(0)
  }
  modhmm_transition_matrix_print(config)
}
