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
	"go/ast"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/parser"
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
	if currentOutput != nil {
		// search create table for target type statement in current output. If not found, add it
		stmt := getUpdateDatabaseTargetStatement(currentOutput.Ast, d.TypeHolder.Name)
		if stmt == nil {
			existingStmts := getUpdateDatabaseStmts(currentOutput.Ast)
			newStmt := getUpdateDatabaseTargetStatement(generatedOutput.Ast, d.TypeHolder.Name)
			// prepend to existing statements and update
			stmtsToSet := append([]ast.Stmt{newStmt}, existingStmts...)
			setStatements(currentOutput.Ast, stmtsToSet)

			return currentOutput, nil
		}

		// if statement exists, existing output is valid
		return nil, nil
	}

	return generatedOutput, nil
}

func findUpdateDatabaseFunction(file *ast.File) *ast.FuncDecl {
	funcs := parser.GetFuncDecls(file)
	for _, f := range funcs {
		if f.Name.Name == "UpdateDatabase" {
			return f
		}
	}
	return nil
}

func getUpdateDatabaseStmts(file *ast.File) []ast.Stmt {
	r := findUpdateDatabaseFunction(file)
	return r.Body.List
}

func getUpdateDatabaseTargetStatement(file *ast.File, typeName string) ast.Stmt {
	r := findUpdateDatabaseFunction(file)
	for _, stmt := range r.Body.List {
		switch stmt.(type) {
		case *ast.IfStmt:
			init := stmt.(*ast.IfStmt).Init
			switch init.(type) {
			case *ast.AssignStmt:
				for _, e := range init.(*ast.AssignStmt).Rhs {
					switch e.(type) {
					case *ast.CallExpr:
						switch e.(*ast.CallExpr).Fun.(type) {
						case *ast.SelectorExpr:
							if e.(*ast.CallExpr).Fun.(*ast.SelectorExpr).Sel.Name == "Create"+typeName+"Table" {
								return stmt
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func setStatements(file *ast.File, stmts []ast.Stmt) {
	r := findUpdateDatabaseFunction(file)
	r.Body.List = stmts
}

func init() {
	makers.Register(&DDL{})
}
