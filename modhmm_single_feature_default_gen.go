// +build ignore

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

import   "log"
import   "net/http"

import   "github.com/shurcooL/vfsgen"

/* -------------------------------------------------------------------------- */

var Model http.FileSystem = http.Dir("modhmm_single_feature_default")

/* -------------------------------------------------------------------------- */

func main() {
  if err := vfsgen.Generate(Model, vfsgen.Options{
    Filename    : "modhmm_single_feature_default.go",
    PackageName : "main",
    VariableName: "assets" }); err != nil {
    log.Fatalln(err)
  }
}
