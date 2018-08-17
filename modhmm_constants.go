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

var coverageList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "control"}

/* -------------------------------------------------------------------------- */

var singleFeatureList = StringList{
  "open", "h3k27ac", "h3k27me3", "h3k9me3", "h3k4me1", "h3k4me3", "h3k4me3o1", "rna", "rna-low", "control"}

/* -------------------------------------------------------------------------- */

var multiFeatureList = StringList{
  "pa", "ea", "bi", "pr", "tr", "tl", "r1", "r2", "ns", "cl"}

/* -------------------------------------------------------------------------- */

var iPA int // promoter active
var iEA int // enhancer active
var iBI int // bivalent state
var iPR int // primed state
var iTR int // transcribed
var iTL int // transcribed (low)
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
var jH3k4me3o1 int
var jRna       int
var jRnaLow    int
var jControl   int

/* -------------------------------------------------------------------------- */

func init() {
  // track indices for multi-feature classifiers
  jOpen      = 2*singleFeatureList.Index("open")
  jH3k27ac   = 2*singleFeatureList.Index("h3k27ac")
  jH3k27me3  = 2*singleFeatureList.Index("h3k27me3")
  jH3k9me3   = 2*singleFeatureList.Index("h3k9me3")
  jH3k4me1   = 2*singleFeatureList.Index("h3k4me1")
  jH3k4me3   = 2*singleFeatureList.Index("h3k4me3")
  jH3k4me3o1 = 2*singleFeatureList.Index("h3k4me3o1")
  jRna       = 2*singleFeatureList.Index("rna")
  jRnaLow    = 2*singleFeatureList.Index("rna-low")
  jControl   = 2*singleFeatureList.Index("control")
  // track indices for modhmm
  iPA = multiFeatureList.Index("pa")
  iEA = multiFeatureList.Index("ea")
  iBI = multiFeatureList.Index("bi")
  iPR = multiFeatureList.Index("pr")
  iTR = multiFeatureList.Index("tr")
  iTL = multiFeatureList.Index("tl")
  iR1 = multiFeatureList.Index("r1")
  iR2 = multiFeatureList.Index("r2")
  iCL = multiFeatureList.Index("cl")
  iNS = multiFeatureList.Index("ns")
}
