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
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
)

type typeField struct {
	Name string
	Type string
}

// GetTypesMaps returns a map of types. Each key is the type name and the
// value is a list of field name and type pairs
func getTypesMaps(filepath string) (map[string][]typeField, error) {
	emptyMap := make(map[string][]typeField)
	f, err := os.Open(filepath)
	if err != nil {
		return emptyMap, err
	}

	buffer, metafile, err := parse(f)
	if err != nil {
		return emptyMap, err
	}

	structsMap, err := getStructs(metafile)
	if err != nil {
		return emptyMap, err
	}

	return decomposeStructs(buffer, structsMap)
}

// Parse parses io reader to ast.File pointer
func parse(reader io.Reader) ([]byte, *ast.File, error) {
	if reader == nil {
		return nil, nil, errors.New("Reader is null")
	}

	var buf bytes.Buffer
	_, err := io.Copy(&buf, reader)
	if err != nil {
		return nil, nil, err
	}

	buffer := buf.Bytes()

	fs := token.NewFileSet()
	// TODO use parser.Trace Mode (last param) instead of 0 to see what is being parsed
	astFile, err := parser.ParseFile(fs, "", buffer, 0)
	return buffer, astFile, err
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

func decomposeStructs(buffer []byte, structs map[string]*ast.StructType) (map[string][]typeField, error) {
	structsMap := make(map[string][]typeField)
	for structName := range structs {
		structMembers, err := decomposeStruct(buffer, structs[structName])
		if err != nil {
			return structsMap, err
		}
		structsMap[structName] = structMembers
	}
	return structsMap, nil
}

func decomposeStruct(buffer []byte, structType *ast.StructType) ([]typeField, error) {
	var fields []typeField

	for _, field := range structType.Fields.List {
		if len(field.Names) != 1 {
			return fields, fmt.Errorf("Unexpected length of %v for a field", len(field.Names))
		}

		fields = append(fields, typeField{
			Name: field.Names[0].Name,
			Type: string(buffer[field.Type.Pos()-1 : field.Type.End()-1]),
		})
	}
	return fields, nil
}
