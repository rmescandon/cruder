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

package io

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

// Content payload in two formats, byte arrays or syntax tree
type Content struct {
	Bytes []byte
	Ast   *ast.File
}

// NewContent returns a pointer to a content struct from a string payload
func NewContent(str string) (*Content, error) {
	ast, err := StringToAST(str)
	if err != nil {
		return nil, err
	}

	return &Content{
		Bytes: []byte(str),
		Ast:   ast,
	}, nil
}

// ByteArrayToAST composes syntax tree from a byte array content
func ByteArrayToAST(buf []byte) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), "", buf, 0)
}

// StringToAST composes syntax tree from a string
func StringToAST(str string) (*ast.File, error) {
	return parser.ParseFile(token.NewFileSet(), "", str, 0)
}

// ASTToString returns a syntax tree as a string
func ASTToString(ast *ast.File) (string, error) {
	b, err := astToBuffer(ast)
	return b.String(), err
}

// ASTToByteArray returns a syntax tree as a byte array
func ASTToByteArray(ast *ast.File) ([]byte, error) {
	b, err := astToBuffer(ast)
	return b.Bytes(), err
}

func astToBuffer(ast *ast.File) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), ast)
	return buf, err
}

// TraceAST prints out AST file content
func TraceAST(f *ast.File) error {
	return ast.Print(token.NewFileSet(), f)
}
