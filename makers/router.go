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

package makers

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/parser"
)

// Router generates service/router.go output go file
type Router struct {
	BaseMaker
}

// ID returns 'router' as this maker identifier
func (r *Router) ID() string {
	return "router"
}

// OutputFilepath returns the path to generated file
func (r *Router) OutputFilepath() string {
	return filepath.Join(config.Config.Output, "/service/router.go")
}

// Make generates the results
func (r *Router) Make() error {
	// Execute the replacement
	logging.Debugf("Loadig template: %v", filepath.Base(r.Template))
	templateContent, err := io.FileToString(r.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := r.TypeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v",
			r.TypeHolder.Name, filepath.Base(r.Template))
	}

	replacedStr, err = config.Config.ReplaceInTemplate(replacedStr)
	if err != nil {
		return fmt.Errorf("Error replacing configuration over template %v",
			filepath.Base(r.Template))
	}

	// Check if output file exists to merge current with existing output
	_, err = os.Stat(r.OutputFilepath())
	if err == nil {
		return r.mergeExistingOutput(replacedStr)
	}

	// Create needed dirs to outputPath and write out substituted string
	ensureDir(filepath.Dir(r.OutputFilepath()))

	io.StringToFile(replacedStr, r.OutputFilepath())

	logging.Infof("Generated: %v", r.OutputFilepath())
	return nil
}

func (r *Router) mergeExistingOutput(replacedStr string) error {
	logging.Infof("Merging new type into: %v", r.OutputFilepath())
	/*generatedAst*/ _, err := io.ByteArrayToAST([]byte(replacedStr))
	if err != nil {
		return err
	}

	// Load current output
	content, err := io.FileToByteArray(r.OutputFilepath())
	if err != nil {
		return err
	}

	currentAst, err := io.ByteArrayToAST(content)
	if err != nil {
		return err
	}

	existingHandlers := findHandlers(currentAst)
	for _, h := range existingHandlers {
		log.Printf("HANDLER: %v", h)
	}

	//TODO IMPLEMENT THE MERGE HERE

	return nil
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

func findHandlers(file *ast.File) []string {
	handlers := []string{}
	stmts := getRouterFunctionStatements(file)
	for _, s := range stmts {
		for _, arg := range s.X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Args {
			for _, ident := range arg.(*ast.CallExpr).Args {
				switch ident.(type) {
				case *ast.Ident:
					handlers = append(handlers, ident.(*ast.Ident).Name)
				}
			}
		}
	}
	return handlers
}

// func findHandlersInRouterFunction(file *ast.File) []string {
// 	handlers := []string{}
// 	for _, decl := range file.Decls {
// 		switch decl.(type) {
// 		case *ast.FuncDecl:
// 			if decl.(*ast.FuncDecl).Name.Name == "Router" {
// 				for _, stmt := range decl.(*ast.FuncDecl).Body.List {
// 					switch stmt.(type) {
// 					case *ast.ExprStmt:
// 						for _, arg := range stmt.(*ast.ExprStmt).X.(*ast.CallExpr).Fun.(*ast.SelectorExpr).X.(*ast.CallExpr).Args {
// 							for _, ident := range arg.(*ast.CallExpr).Args {
// 								switch ident.(type) {
// 								case *ast.Ident:
// 									handlers = append(handlers, ident.(*ast.Ident).Name)
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return handlers
// }

func init() {
	register(&Router{})
}
