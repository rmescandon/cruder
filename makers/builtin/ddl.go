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

package builtin

import (
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
)

// DDL makes the controller
type DDL struct {
	makers.Base
}

// ID returns the identifier 'handler' for this maker
func (d *DDL) ID() string {
	return "ddl"
}

// OutputFilepath returns the path to the generated file
func (d *DDL) OutputFilepath() string {
	return filepath.Join(makers.BasePath, "datastore/ddl.go")
}

// Make generates the results
func (d *DDL) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {

	// TODO TRACE
	generatedOutput.Trace()

	// TODO IMPLEMENT
	return nil, nil
}

func init() {
	makers.Register(&DDL{})
}
