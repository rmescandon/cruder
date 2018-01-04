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

package src

// This file contains the model construction by parsing source files.

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/rmescandon/cruder/io"
)

func ComposeTypeHolders(source *io.GoFile) ([]*TypeHolder, error) {
	var holders []*TypeHolder
	decls := getTypeDecls(source.Ast)
	for _, decl := range decls {
		for _, spec := range decl.Specs {
			fields, err := composeTypeFields(source.Content, spec)
			if err != nil {
				return holders, err
			}

			idField := TypeField{}
			if len(fields) > 0 {
				idField = fields[0]
			}

			holders = append(holders, &TypeHolder{
				Name:    spec.(*ast.TypeSpec).Name.Name,
				Source:  source,
				IDField: idField,
				Fields:  fields,
				Decl:    decl,
			})
		}
	}
	return holders, nil
}

func getTypeDecls(file *ast.File) []*ast.GenDecl {
	var typeDecls []*ast.GenDecl
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			if decl.(*ast.GenDecl).Tok == token.TYPE {
				typeDecls = append(typeDecls, decl.(*ast.GenDecl))
			}
		}
	}
	return typeDecls
}

func composeTypeFields(content []byte, spec ast.Spec) ([]TypeField, error) {
	var fields []TypeField
	for _, field := range spec.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		if len(field.Names) != 1 {
			return fields, fmt.Errorf("Unexpected length of %v for a field", len(field.Names))
		}

		fields = append(fields, TypeField{
			Name: field.Names[0].Name,
			Type: string(content[field.Type.Pos()-1 : field.Type.End()-1]),
		})
	}
	return fields, nil
}
