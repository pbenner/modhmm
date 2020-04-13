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

import . "github.com/pbenner/modhmm/config"

/* -------------------------------------------------------------------------- */

var iPA int // promoter active
var iEA int // enhancer active
var iBI int // bivalent state
var iPR int // primed state
var iTR int // transcribed
var iR1 int // repressed h3k27me3
var iR2 int // repressed h3k9me3
var iCL int // control
var iNS int // no signal

/* -------------------------------------------------------------------------- */

var jOpen      int
var jH3k27ac   int
var jH3k27me3  int
var jH3k9me3   int
var jH3k4me1   int
var jH3k4me3   int
var jRna       int
var jControl   int

/* -------------------------------------------------------------------------- */

func init() {
  // track indices for multi-feature classifiers
  jOpen     = EnrichmentList.Index("open")
  jH3k27ac  = EnrichmentList.Index("h3k27ac")
  jH3k27me3 = EnrichmentList.Index("h3k27me3")
  jH3k9me3  = EnrichmentList.Index("h3k9me3")
  jH3k4me1  = EnrichmentList.Index("h3k4me1")
  jH3k4me3  = EnrichmentList.Index("h3k4me3")
  jRna      = EnrichmentList.Index("rna")
  jControl  = EnrichmentList.Index("control")
  // track indices for modhmm
  iPA = ChromatinStateList.Index("pa")
  iEA = ChromatinStateList.Index("ea")
  iBI = ChromatinStateList.Index("bi")
  iPR = ChromatinStateList.Index("pr")
  iTR = ChromatinStateList.Index("tr")
  iR1 = ChromatinStateList.Index("r1")
  iR2 = ChromatinStateList.Index("r2")
  iCL = ChromatinStateList.Index("cl")
  iNS = ChromatinStateList.Index("ns")
}
