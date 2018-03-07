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

package parser

// This file contains the model construction by parsing source files.

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
)

// ComposeTypeHolders composes the type holders for the types in source file
func ComposeTypeHolders(source *io.GoFile) ([]*TypeHolder, error) {
	var holders []*TypeHolder
	decls := getTypeDecls(source.Ast)

	for _, decl := range decls {
		for _, spec := range decl.Specs {
			var buf bytes.Buffer
			printer.Fprint(&buf, token.NewFileSet(), spec)
			log.Info(string(buf.Bytes()))

			fields, err := composeTypeFields(spec)
			if err != nil {
				return []*TypeHolder{}, err
			}

			// validate that there are at least two fields,
			// first will be taken as the ID and second as the search field
			name := spec.(*ast.TypeSpec).Name.Name
			if len(fields) < 2 {
				return holders, fmt.Errorf("Found less than 2 fields for type %v", name)
			}

			holders = append(holders, &TypeHolder{
				Name:   name,
				Source: source,
				Fields: fields,
				Decl:   decl,
			})
		}
	}

	return holders, nil
}

func getTypeDecls(file *ast.File) []*ast.GenDecl {
	typeDecls := []*ast.GenDecl{}
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

// GetFuncDecls returns the list of func declarations
func GetFuncDecls(file *ast.File) []*ast.FuncDecl {
	funcDecls := []*ast.FuncDecl{}
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			funcDecls = append(funcDecls, decl.(*ast.FuncDecl))
		}
	}
	return funcDecls
}

func getInterfaces(file *ast.File) []*ast.GenDecl {
	interfaces := []*ast.GenDecl{}
	typeDecls := getTypeDecls(file)
	for _, decl := range typeDecls {
		for _, spec := range decl.Specs {
			switch spec.(*ast.TypeSpec).Type.(type) {
			case *ast.InterfaceType:
				interfaces = append(interfaces, decl)
			}
		}
	}
	return interfaces
}

// GetInterface returns certain interface identified by name
func GetInterface(file *ast.File, name string) *ast.InterfaceType {
	interfaces := getInterfaces(file)
	for _, decl := range interfaces {
		for _, spec := range decl.Specs {
			if spec.(*ast.TypeSpec).Name.Name == name {
				return spec.(*ast.TypeSpec).Type.(*ast.InterfaceType)
			}
		}
	}
	return nil
}

// GetInterfaceMethods returns the list of methods in a declared interface
func GetInterfaceMethods(iface *ast.InterfaceType) []*ast.Field {
	if iface.Methods == nil {
		return nil
	}
	return iface.Methods.List
}

// HasMethod returns true if found method into iface
func HasMethod(iface *ast.InterfaceType, methodName string) bool {
	if iface.Methods == nil {
		return false
	}

	for _, field := range iface.Methods.List {
		if len(field.Names) == 0 {
			continue
		}

		if field.Names[0].Name == methodName {
			return true
		}
	}

	return false
}

// AddMethod modyfies iface by adding method
func AddMethod(iface *ast.InterfaceType, method *ast.Field) {
	if iface.Methods == nil {
		iface.Methods = &ast.FieldList{List: []*ast.Field{method}}
	} else if iface.Methods.List == nil {
		iface.Methods.List = []*ast.Field{method}
	} else {
		iface.Methods.List = append(iface.Methods.List, method)
	}
}

func composeTypeFields(spec ast.Spec) ([]TypeField, error) {
	var fields []TypeField
	for _, field := range spec.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		if len(field.Names) != 1 {
			return []TypeField{}, fmt.Errorf("Unexpected length of %v for a field", len(field.Names))
		}

		fields = append(fields, TypeField{
			Name: field.Names[0].Name,
			Type: fmt.Sprintf("%s", field.Type),
		})
	}
	return fields, nil
}
