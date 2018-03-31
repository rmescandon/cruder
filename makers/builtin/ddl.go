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

	"github.com/rmescandon/cruder/errs"
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
	if generatedOutput == nil {
		return nil, errs.ErrNoContent
	}

	if currentOutput != nil {
		// search create table for target type statement in current output. If not found, add it
		stmt, err := getUpdateDatabaseTargetStatement(currentOutput.Ast, d.TypeHolder.Name)
		if err != nil {
			return nil, err
		}

		if stmt == nil {
			existingStmts, err := getUpdateDatabaseStmts(currentOutput.Ast)
			if err != nil {
				return nil, err
			}

			newStmt, err := getUpdateDatabaseTargetStatement(generatedOutput.Ast, d.TypeHolder.Name)
			if err != nil {
				return nil, err
			}

			// prepend to existing statements and update
			stmtsToSet := append([]ast.Stmt{newStmt}, existingStmts...)
			err = setStatements(currentOutput.Ast, stmtsToSet)
			if err != nil {
				return nil, err
			}

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

func getUpdateDatabaseStmts(file *ast.File) ([]ast.Stmt, error) {
	r := findUpdateDatabaseFunction(file)
	if r == nil {
		return []ast.Stmt{}, errs.NewErrNotFound("UpdateDatabase function")
	}
	return r.Body.List, nil
}

func getUpdateDatabaseTargetStatement(file *ast.File, typeName string) (ast.Stmt, error) {
	r := findUpdateDatabaseFunction(file)
	if r == nil {
		return nil, errs.NewErrNotFound("UpdateDatabase function")
	}

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
								return stmt, nil
							}
						}
					}
				}
			}
		}
	}
	// If not found the statement, simply return null, but it is not an error
	return nil, nil
}

func setStatements(file *ast.File, stmts []ast.Stmt) error {
	r := findUpdateDatabaseFunction(file)
	if r == nil {
		return errs.NewErrNotFound("UpdateDatabase function")
	}
	r.Body.List = stmts
	return nil
}

func init() {
	makers.Register(&DDL{})
}
