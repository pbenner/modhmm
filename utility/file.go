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

package utility

/* -------------------------------------------------------------------------- */

import "io"
import "net/http"
import "os"

/* -------------------------------------------------------------------------- */

func FileExists(filename string) bool {
  if _, err := os.Stat(filename); err != nil {
    return false
  } else {
    return true
  }
}

func FileCheckMark(filename string) string {
  if !FileExists(filename) {
    return "\xE2\x9C\x97"
  } else {
    return "\xE2\x9C\x93"
  }
}

func DownloadFile(filepath string, url string) error {
  resp, err := http.Get(url)
  if err != nil {
    return err
  }
  defer resp.Body.Close()

  out, err := os.Create(filepath)
  if err != nil {
    return err
  }
  defer out.Close()

  _, err = io.Copy(out, resp.Body)
  return err
}
