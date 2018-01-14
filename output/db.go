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
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
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

// Make generates the results
func (db *Db) Make() error {
	for i := range db.TypeHolders {
		err := db.makeOne(i)
		if err != nil {
			return err
		}
	}
	return nil
}

// mergeExistingOutput resolves the conflict when already exists an output file
func (db *Db) mergeExistingOutput(replacedStr string) error {
	generatedAst, err := io.ByteArrayToAST([]byte(replacedStr))
	if err != nil {
		return err
	}

	// load current output
	content, err := io.FileToByteArray(db.File.Path)
	if err != nil {
		return err
	}

	currentAst, err := io.ByteArrayToAST(content)
	if err != nil {
		return err
	}

	generatedIface := src.GetInterface(generatedAst, "Datastore")
	currentIface := src.GetInterface(currentAst, "Datastore")

	// search for generatedIface methods into currentIface and add them if not found
	for _, method := range src.GetInterfaceMethods(generatedIface) {
		if !src.HasMethod(currentIface, method.Names[0].Name) {
			src.AddMethod(currentIface, method)
		}
	}

	// write out the resultant modified Datastore interface to output
	// TODO VERIFY that using pointers is enough to alter generatedAst before writting out
	io.ASTToFile(generatedAst, db.File.Path)

	return nil
}

// makeOne runs to generate a single output result
func (db *Db) makeOne(index int) error {
	// execute the replacement
	logging.Debugf("Loadig template: %v", filepath.Base(db.Template))
	templateContent, err := io.FileToString(db.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := db.TypeHolders[index].ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", db.TypeHolders[index].Name, filepath.Base(db.Template))
	}

	// check if output file exists
	_, err = os.Stat(db.File.Path)
	if err == nil {
		db.mergeExistingOutput(replacedStr)
	} else {
		// write out generated ast
		// create needed dirs to outputPath
		ensureDir(filepath.Dir(db.File.Path))

		io.StringToFile(replacedStr, db.File.Path)

		logging.Infof("Generated: %v", db.File.Path)
	}

	/*
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
	*/
	return nil
}

/*
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
*/
