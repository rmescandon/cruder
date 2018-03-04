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
	"go/ast"
)

// Content payload in two formats, byte arrays or syntax tree
type Content struct {
	Ast *ast.File
}

// NewContent returns a pointer to a content struct from a string payload
func NewContent(str string) (*Content, error) {
	ast, err := StringToAST(str)
	if err != nil {
		return nil, err
	}

	return &Content{ast}, nil
}

// Bytes returns the content as a byte array
func (c *Content) Bytes() ([]byte, error) {
	return ASTToByteArray(c.Ast)
}

func (c *Content) String() (string, error) {
	return ASTToString(c.Ast)
}
