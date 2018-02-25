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

package main

import (
	"path/filepath"

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
)

// Reply struct holding data to copy reply template
type Reply struct {
	makers.Base
}

// ID returns 'reply' as this maker identifier
func (r *Reply) ID() string {
	return "reply"
}

// OutputFilepath returns the path to the output file
func (r *Reply) OutputFilepath() string {
	return filepath.Join(makers.BasePath, "handler/reply.go")
}

// Make copies template to output path
func (r *Reply) Make(generatedOutput *io.Content, currentOutput *io.Content) (string, error) {
	if currentOutput != nil {
		return "", errs.NewErrOutputExists(r.OutputFilepath())
	}

	return string(generatedOutput.Bytes), nil
}

func init() {
	makers.Register(&Reply{})
}
