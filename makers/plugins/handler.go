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

package main

import (
	"path/filepath"
	"strings"

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
)

// Handler makes the controller
type Handler struct {
	makers.Base
}

// ID returns the identifier 'handler' for this maker
func (h *Handler) ID() string {
	return "handler"
}

// OutputFilepath returns the path to the generated file
func (h *Handler) OutputFilepath() string {
	return filepath.Join(
		makers.BasePath,
		h.ID(),
		strings.ToLower(h.TypeHolder.Identifier())+".go")
}

// Make generates the results
func (db *Db) Make(generatedOutput *io.Content, currentOutput *io.Content) (string, error) {
	if currentOutput != nil {
		return "", errs.NewErrOutputExists(h.OutputFilepath())
	}

	return string(generatedOutput.Ast)
}

func init() {
	makers.Register(&Handler{})
}
