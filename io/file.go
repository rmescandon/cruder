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
)

// ByteArrayToAST composes syntax tree from a byte array content
func ByteArrayToAST(content []byte) (*ast.File, error) {
	// TODO use parser.Trace Mode (last param) instead of 0 to see what is being parsed
	return parser.ParseFile(token.NewFileSet(), "", content, 0)
}

// FileToString reads file content and stores it in a string
func FileToString(filepath string) (string, error) {
	b, err := fileToBuffer(filepath)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// FileToByteArray reads file content and stores it in a byte array
func FileToByteArray(filepath string) ([]byte, error) {
	b, err := fileToBuffer(filepath)
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
func ASTToFile(ast *ast.File, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", filepath, err)
	}
	defer f.Close()

	return printer.Fprint(f, token.NewFileSet(), ast)
}

// writeToFile writes a string content to a file
func writeToFile(content interface{}, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", filepath, err)
	}
	defer f.Close()

	switch content.(type) {
	case []byte:
		_, err = f.Write(content.([]byte))
	case string:
		_, err = f.WriteString(content.(string))
	}
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", filepath, err)
	}

	return nil
}

func fileToBuffer(filepath string) (*bytes.Buffer, error) {
	f, err := os.Open(filepath)
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
