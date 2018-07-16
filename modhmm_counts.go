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

import   "io"

import . "github.com/pbenner/ngstat/config"
import . "github.com/pbenner/autodiff"

/* -------------------------------------------------------------------------- */

type Counts struct {
  X []float64
  Y []int
}

/* -------------------------------------------------------------------------- */

func (config *Counts) Import(reader io.Reader, args... interface{}) error {
  return JsonImport(reader, config)
}

func (config *Counts) Export(writer io.Writer) error {
  return JsonExport(writer, config)
}

func (config *Counts) ImportFile(filename string) error {
  if err := ImportFile(config, filename, BareRealType); err != nil {
    return err
  }
  return nil
}

func (config *Counts) ExportFile(filename string) error {
  if err := ExportFile(config, filename); err != nil {
    return err
  }
  return nil
}

