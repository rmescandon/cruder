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

package io

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

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

// FileToString reads file content and stores it in a string
func FileToString(file string) (string, error) {
	b, err := fileToBuffer(file)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// FileToByteArray reads file content and stores it in a byte array
func FileToByteArray(file string) ([]byte, error) {
	b, err := fileToBuffer(file)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// StringToFile writes a strint to a file
func StringToFile(content, filepath string) error {
	return writeToFile(content, filepath)
}

// ByteArrayToFile writes a buffer to a file
func ByteArrayToFile(content []byte, filepath string) error {
	return writeToFile(content, filepath)
}

// ASTToFile writes a syntax tree to file
func ASTToFile(ast *ast.File, file string) error {
	err := EnsureDir(filepath.Dir(file))
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", file, err)
	}
	defer f.Close()

	return printer.Fprint(f, token.NewFileSet(), ast)
}

// writeToFile writes a string content to a file
func writeToFile(content interface{}, file string) error {
	err := EnsureDir(filepath.Dir(file))
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", file, err)
	}
	defer f.Close()

	switch content.(type) {
	case []byte:
		_, err = f.Write(content.([]byte))
	case string:
		_, err = f.WriteString(content.(string))
	}
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", file, err)
	}

	return nil
}

func fileToBuffer(file string) (*bytes.Buffer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
