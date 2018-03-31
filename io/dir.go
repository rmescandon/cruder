// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package io

import (
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir checks for a dir existence and creates it if not exists
func EnsureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// NormalizePath takes a path and resolves ~ to HOME env value and
// returns abs path
func NormalizePath(ptrStr *string) error {
	if strings.Contains(*ptrStr, "~") {
		*ptrStr = strings.Replace(*ptrStr, "~", os.Getenv("HOME"), -1)
	}

	var err error
	*ptrStr, err = filepath.Abs(*ptrStr)
	return err
}
