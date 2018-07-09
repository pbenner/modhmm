/* Copyright (C) 2016-2018 Philipp Benner
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

//import "fmt"
import "os"
import "path"

/* -------------------------------------------------------------------------- */

func updateRequired(config ConfigModHmm, target string, deps ...string) bool {
  if s1, err := os.Stat(target); err != nil {
    printStderr(config, 2, "Target `%s' does not exist...\n", path.Base(target))
    return true
  } else {
    for _, dep := range deps {
      if s2, err := os.Stat(dep); err == nil {
        if s1.ModTime().Before(s2.ModTime()) {
          printStderr(config, 2, "Target `%s' required update...\n", path.Base(target))
          return true
        }
      }
    }
  }
  printStderr(config, 2, "Target `%s' is up to date...\n", path.Base(target))
  return false
}
