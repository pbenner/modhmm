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
import   "os"

import . "github.com/pbenner/modhmm/config"

/* file utilities
 * -------------------------------------------------------------------------- */

func updateRequired(config ConfigModHmm, target TargetFile, deps ...string) bool {
  if s1, err := os.Stat(target.Filename); err != nil {
    if target.Static {
      log.Fatalf("Target `%s' is marked static but does not exist\n", target)
    }
    printStderr(config, 2, "Target `%s' does not exist...\n", target.Filename)
    return true
  } else {
    if target.Static {
      printStderr(config, 2, "Target `%s' is static and requires no update...\n", target)
      return false
    }
    for _, dep := range deps {
      if s2, err := os.Stat(dep); err == nil {
        if s1.ModTime().Before(s2.ModTime()) {
          printStderr(config, 2, "Target `%s' requires update...\n", target.Filename)
          printStderr(config, 3, " -> `%s' has more recent timestamp\n", dep)
          return true
        }
      }
    }
  }
  printStderr(config, 2, "Target `%s' is up to date...\n", target.Filename)
  return false
}
