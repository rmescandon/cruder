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

package cruder

// This file contains the model construction by parsing source files.

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

// GetTypesMaps returns a map of types
func GetTypesMaps(filepath string) (map[string]map[string]string, error) {
	emptyMap := make(map[string]map[string]string)
	f, err := os.Open(filepath)
	if err != nil {
		return emptyMap, err
	}

	metafile, err := parse(f)
	if err != nil {
		return emptyMap, err
	}

	structsList, err := getStructs(metafile)
	if err != nil {
		return emptyMap, err
	}

	return decomposeStructs(structsList)
}

// Parse parses io reader to ast.File pointer
func parse(reader io.Reader) (*ast.File, error) {
	if reader == nil {
		return nil, errors.New("Reader is null")
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, err
	}

	buffer := buf.Bytes()

	fs := token.NewFileSet()
	return parser.ParseFile(fs, "", buffer, parser.Trace)
}

func getStructs(file *ast.File) (map[string]*ast.StructType, error) {
	structs := make(map[string]*ast.StructType)

	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			genDecl := decl.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structName := typeSpec.Name.Name
						structType := typeSpec.Type.(*ast.StructType)
						structs[structName] = structType
					}
				}
			}
		}
	}

	return structs, nil
}

func decomposeStructs(structs map[string]*ast.StructType) (map[string]map[string]string, error) {
	structsMap := make(map[string]map[string]string)
	for structName := range structs {
		structMembers, err := decomposeStruct(structs[structName])
		if err != nil {
			return structsMap, err
		}
		structsMap[structName] = structMembers
	}
	return structsMap, nil
}

func decomposeStruct(structType *ast.StructType) (map[string]string, error) {
	fields := make(map[string]string)
	for _, field := range structType.Fields.List {
		if len(field.Names) == 2 {
			fields[field.Names[0].Name] = field.Names[1].Name
		}
	}
	return fields, nil
}
