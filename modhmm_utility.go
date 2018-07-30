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

import "fmt"
import "os"
import "reflect"
import "strings"

/* struct operations
 * -------------------------------------------------------------------------- */

func collectStrings(config interface{}) []string {
  r := []string{}
  v := reflect.ValueOf(config)

  for i := 0; i < v.NumField(); i++ {
    switch v.Field(i).Kind() {
    case reflect.Struct:
      r = append(r, collectStrings(v.Field(i).Interface())...)
    case reflect.String:
      r = append(r, v.Field(i).String())
    default:
      panic("internal error")
    }
  }
  return r
}

func getField(config interface{}, field string) reflect.Value {
  field = strings.Replace(field, "-", "_", -1)
  v := reflect.ValueOf(config)
  switch v.Kind() {
  case reflect.Struct:
    if r := reflect.Indirect(v).FieldByName(field); r.IsValid() {
      return r
    } else {
      if s := reflect.Indirect(v).FieldByName(strings.Title(field)); s.IsValid() {
        return s
      }
    }
  }
  panic(fmt.Sprintf("internal error: `%s' not found in struct", field))
}

func getFieldAsString(config interface{}, field string) string {
  return getField(config, field).String()
}

func getFieldAsStringSlice(config interface{}, field string) []string {
  v := getField(config, field)
  switch v.Kind() {
  case reflect.Slice:
    r := make([]string, v.Len())
    for i := 0; i < v.Len(); i++ {
      r[i] = v.Index(i).String()
    }
    return r
  }
  panic("internal error")
}

/* file utilities
 * -------------------------------------------------------------------------- */

func fileExists(filename string) bool {
  if _, err := os.Stat(filename); err != nil {
    return false
  } else {
    return true
  }
}

func updateRequired(config ConfigModHmm, target string, deps ...string) bool {
  if s1, err := os.Stat(target); err != nil {
    printStderr(config, 2, "Target `%s' does not exist...\n", target)
    return true
  } else {
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
