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

// Router generates service/router.go output go file
type Router struct {
	makers.Base
}

// ID returns 'router' as this maker identifier
func (r *Router) ID() string {
	return "router"
}

// OutputFilepath returns the path to generated file
func (r *Router) OutputFilepath() string {
	return filepath.Join(makers.BasePath, "/service/router.go")
}

// Make copies template to output path
func (r *Router) Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error) {
	if currentOutput != nil {
		// Search generated handlers amongst existing ones and add only new ones
		stmts := getRouterFunctionStatements(currentOutput.Ast)
		existingHandlers := findHandlersInStatements(stmts)
		generatedHandlers := findHandlers(generatedOutput.Ast)
		stmtsToAdd := []ast.Stmt{}
		for k := range generatedHandlers {
			if _, ok := existingHandlers[k]; !ok {
				stmtsToAdd = append(stmtsToAdd, generatedHandlers[k])
			}
		}

		addStatements(currentOutput.Ast, stmtsToAdd)

		return currentOutput, nil
	}

	return generatedOutput, nil
}

func addStatements(file *ast.File, stmts []ast.Stmt) {
	r := findRouterFunction(file)
	for i, stmt := range r.Body.List {
		switch stmt.(type) {
		case *ast.ExprStmt:
			// once found first expr statement, insert all new here
			r.Body.List = append(r.Body.List[:i],
				append(stmts, r.Body.List[i:]...)...)
			return
		}
	}
}

func findRouterFunction(file *ast.File) *ast.FuncDecl {
	funcs := parser.GetFuncDecls(file)
	for _, f := range funcs {
		if f.Name.Name == "Router" {
			return f
		}
	}
	return nil
}

func getRouterFunctionStatements(file *ast.File) []*ast.ExprStmt {
	stmts := []*ast.ExprStmt{}
	r := findRouterFunction(file)
	for _, stmt := range r.Body.List {
		switch stmt.(type) {
		case *ast.ExprStmt:
			stmts = append(stmts, stmt.(*ast.ExprStmt))
		}
	}
	return stmts
}

func findHandlers(file *ast.File) map[string]*ast.ExprStmt {
	stmts := getRouterFunctionStatements(file)
	return findHandlersInStatements(stmts)
}

func findHandlersInStatements(stmts []*ast.ExprStmt) map[string]*ast.ExprStmt {
	handlers := make(map[string]*ast.ExprStmt)
	for _, s := range stmts {
		for _, arg := range s.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Args {
			for _, ident := range arg.(*ast.CallExpr).Args {
				switch ident.(type) {
				case *ast.Ident:
					handlers[ident.(*ast.Ident).Name] = s
				}
			}
		}
	}
	return handlers
}

func init() {
	makers.Register(&Router{})
}
