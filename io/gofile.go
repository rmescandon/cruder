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

// GoFile represents a go source code disk resource. Each struct
// member is a different way of having same info
type GoFile struct {
	Path    string
	Content []byte
	Ast     *ast.File
}

// NewGoFile returns a brand new GoFile instance
func NewGoFile(filepath string) (*GoFile, error) {
	content, err := FileToByteArray(filepath)
	if err != nil {
		return nil, err
	}

	ast, err := ByteArrayToAST(content)
	if err != nil {
		return nil, err
	}

	return &GoFile{
		Path:    filepath,
		Content: content,
		Ast:     ast,
	}, nil
}
