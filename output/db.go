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

package output

import (
	"go/ast"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
)

// Db maker to include types in datastore interface
type Db struct {
	TypeHolders []*src.TypeHolder
	File        *io.GoFile
	Template    string
}

// OutputFilepath returns the path to generated file
func (db *Db) OutputFilepath() string {
	return db.File.Path
}

// Run runs to generate the type methods in Datastore interface for this type
func (db *Db) Run() error {
	fileToLoad := db.Template

	// check if output file exists
	_, err := os.Stat(db.File.Path)
	if err == nil {
		// if output file exists, load it
		fileToLoad = db.File.Path
	} else {
		// create needed dirs to outputPath
		ensureDir(filepath.Dir(db.File.Path))
	}

	db.File.Content, err = io.FileToByteArray(fileToLoad)
	if err != nil {
		return err
	}
	db.File.Ast, err = io.ByteArrayToAST(db.File.Content)
	if err != nil {
		return err
	}

	//FIXME PoC to see if functions for type are added
	for _, x := range db.File.Ast.Decls {
		if x, ok := x.(*ast.GenDecl); ok {
			if x.Tok != token.TYPE {
				continue
			}
			for _, x := range x.Specs {
				if x, ok := x.(*ast.TypeSpec); ok {
					iname := x.Name
					if x, ok := x.Type.(*ast.InterfaceType); ok {
						if iname == "Datastore" {
							// Insert new functions here
							// See
							// https://stackoverflow.com/questions/33836358/parsing-go-src-trying-to-convert-ast-gendecl-to-types-interface
							astFunc := GenerateASTFunction()
							
						}
					}
				}
			}
		}
	  }

	foundFirstFunc := false
	for i, decl := range db.File.Ast.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			if ast.FuncDecl.Name == 
			outputAst.Decls = append(outputAst.Decls[:i], append([]ast.Decl{ds.TypeHolders[index].Decl}, outputAst.Decls[i:]...)...)
			foundFirstFunc = true
		}

		if foundFirstFunc {
			break
		}
	}

}

// A test for DeleteMyType(id int) error
func GenerateASTFunction() *ast.FuncDecl {
	f := &ast.FuncDecl{
		Name: ast.NewIdent("DeleteMyType"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("id")},
						Type:  ast.NewIdent("int"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{ast.NewIdent("error")},
						Type:  ast.NewIdent("error"),
					},
				},
			},
		},
	}
}
