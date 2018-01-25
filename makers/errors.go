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

package makers

import (
	"errors"
	"fmt"
)

// Errors
var (
	ErrNoMakerRegistered = errors.New("Any maker has been registered")
)

// ErrOutputExists error struct for an existing output file
type ErrOutputExists struct {
	Path string
}

// Error returns the error string
func (e ErrOutputExists) Error() string {
	return fmt.Sprintf("File %v already exists. Skip writting", e.Path)
}

// NewErrOutputExists returns a new ErrOutputExists struct
func NewErrOutputExists(output string) ErrOutputExists {
	return ErrOutputExists{Path: output}
}

// ErrNotFound error struct for a not existing thing
type ErrNotFound struct {
	What string
}

// Error returns the error string
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%v not Found", e.What)
}

// NewErrNotFound returns a new ErrNotFound
func NewErrNotFound(what string) ErrNotFound {
	return ErrNotFound{What: what}
}
