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

import "log"
import "os"
import "strings"

/* file utilities
 * -------------------------------------------------------------------------- */

func fileExists(filename string) bool {
  if _, err := os.Stat(filename); err != nil {
    return false
  } else {
    return true
  }
}

func updateRequired(config ConfigModHmm, target TargetFile, deps ...string) bool {
  if s1, err := os.Stat(target.Filename); err != nil {
    if target.Static {
      log.Fatalf("Target `%s' is marked static but does not exist\n", target)
    }
    printStderr(config, 2, "Target `%s' does not exist...\n", target)
    return true
  } else {
    if target.Static {
      printStderr(config, 2, "Target `%s' is static and requires no update...\n", target)
      return false
    }
    for _, dep := range deps {
      if s2, err := os.Stat(dep); err == nil {
        if s1.ModTime().Before(s2.ModTime()) {
          printStderr(config, 2, "Target `%s' requires update...\n", target)
          return true
        }
      }
    }
  }
  printStderr(config, 2, "Target `%s' is up to date...\n", target)
  return false
}

/* -------------------------------------------------------------------------- */

// Divide a by b, the result is rounded down.
func divIntDown(a, b int) int {
  return a/b
}

// Divide a by b, the result is rounded up.
func divIntUp(a, b int) int {
  return (a+b-1)/b
}

/* -------------------------------------------------------------------------- */

type StringList []string

func (s StringList) Index(item string) int {
  for i, x := range s {
    if item == x {
      return i
    }
  }
  return -1
}

func (s StringList) Contains(item string) bool {
  for _, x := range s {
    if item == x {
      return true
    }
  }
  return false
}

func (s StringList) Intersection(x []string) []string {
  r := []string{}
  for _, elem := range x {
    if s.Contains(elem) {
      r = append(r, elem)
    }
  }
  return r
}

/* -------------------------------------------------------------------------- */

type InsensitiveStringList []string

func (s InsensitiveStringList) Index(item string) int {
  item = strings.ToLower(item)
  for i, x := range s {
    if item == x {
      return i
    }
  }
  return -1
}

func (s InsensitiveStringList) Contains(item string) bool {
  item = strings.ToLower(item)
  for _, x := range s {
    if item == x {
      return true
    }
  }
  return false
}

func (s InsensitiveStringList) Intersection(x []string) []string {
  r := []string{}
  for _, elem := range x {
    elem = strings.ToLower(elem)
    if s.Contains(elem) {
      r = append(r, elem)
    }
  }
  return r
}
