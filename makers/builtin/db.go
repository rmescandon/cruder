// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017 Roberto Mier Escandon <rmescandon@gmail.com>
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
	"fmt"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/parser"
)

// Db maker to include types in datastore interface
type Db struct {
	makers.Base
}

// ID returns 'db'
func (db *Db) ID() string {
	return "db"
}

// OutputFilepath returns the path to generated file
func (db *Db) OutputFilepath() string {
	return filepath.Join(makers.BasePath, fmt.Sprintf("datastore/%v.go", db.ID()))
}

// Make generates the results
func (db *Db) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {
	if currentOutput != nil {
		generatedIface := parser.GetInterface(generatedOutput.Ast, "Datastore")
		currentIface := parser.GetInterface(currentOutput.Ast, "Datastore")

		// Search for generatedIface methods into currentIface and add them if not found
		for _, method := range parser.GetInterfaceMethods(generatedIface) {
			if !parser.HasMethod(currentIface, method.Names[0].Name) {
				parser.AddMethod(currentIface, method)
			}
		}

		return currentOutput, nil
	}

	return generatedOutput, nil
}

func init() {
	makers.Register(&Db{})
}
