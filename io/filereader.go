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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// ByteContentAsAST composes syntax tree from a byte array content
func ByteContentAsAST(content []byte) (*ast.File, error) {
	// TODO use parser.Trace Mode (last param) instead of 0 to see what is being parsed
	return parser.ParseFile(token.NewFileSet(), "", content, 0)
}

// FileContentAsString reads file content and stores it in a string
func FileContentAsString(filepath string) (string, error) {
	b, err := fileBuffer(filepath)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

// FileContentAsByteArray reads file content and stores it in a byte array
func FileContentAsByteArray(filepath string) ([]byte, error) {
	b, err := fileBuffer(filepath)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func fileBuffer(filepath string) (*bytes.Buffer, error) {
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
