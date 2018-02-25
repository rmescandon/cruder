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

package main

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
)

// Datastore generates datastore/<type>.go output go file
type Datastore struct {
	makers.Base
}

// ID returns 'datastore' as this maker identifier
func (ds *Datastore) ID() string {
	return "datastore"
}

// OutputFilepath returns the path to generated file
func (ds *Datastore) OutputFilepath() string {
	return filepath.Join(
		makers.BasePath,
		ds.ID(),
		strings.ToLower(ds.TypeHolder.Name)+".go")
}

// Make generates the result
func (ds *Datastore) Make(generatedOutput *io.Content, currentOutput *io.Content) (string, error) {
	if currentOutput != nil {
		// in case if does exist, it should match the types file. Otherwise it's an error
		if ds.OutputFilepath() != ds.TypeHolder.Source.Path {
			return errs.NewErrOutputExists(ds.OutputFilepath())
		}

		// if output file is the same as types one, add the type to the generated output
		// - prepend GenType AST to it
		// - write out AST to output, overwriting
		// insert GenType AST just before first function
		foundFirstFunc := false
		for i, decl := range currentOutput.Ast.Decls {
			switch decl.(type) {
			case *ast.FuncDecl:
				// trick to prepend instead of appending
				currentOutput.Ast.Decls = append(currentOutput.Ast.Decls[:i],
					append([]ast.Decl{ds.TypeHolder.Decl}, currentOutput.Ast.Decls[i:]...)...)
				foundFirstFunc = true
			}

			if foundFirstFunc {
				break
			}
		}

		return io.ASTToString(currentOutput.Ast)
	}

	return string(generatedOutput.Bytes)
}

func init() {
	makers.Register(&Datastore{})
}
