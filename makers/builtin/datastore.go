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
	if ds.TypeHolder == nil || len(ds.TypeHolder.Name) == 0 {
		return ""
	}

	return filepath.Join(
		makers.BasePath,
		ds.ID(),
		strings.ToLower(ds.TypeHolder.Name)+".go")
}

// Make generates the result
func (ds *Datastore) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {
	if generatedOutput == nil {
		return nil, errs.ErrNoContent
	}

	// always include type definition into this datastore generated file
	foundFirstFunc := false
	for i, decl := range generatedOutput.Ast.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			// trick to prepend instead of appending. The idea here is to include type
			// definition just before first found function in generated output
			generatedOutput.Ast.Decls = append(generatedOutput.Ast.Decls[:i],
				append([]ast.Decl{ds.TypeHolder.Decl}, generatedOutput.Ast.Decls[i:]...)...)
			foundFirstFunc = true
		}

		if foundFirstFunc {
			break
		}
	}

	if !foundFirstFunc {
		return nil, errs.NewErrNotFound("First function in generated output")
	}

	return generatedOutput, nil
}

func init() {
	makers.Register(&Datastore{})
}
