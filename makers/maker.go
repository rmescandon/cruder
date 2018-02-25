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

package makers

import (
	"fmt"
	"path/filepath"

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
)

// BasePath local folder taken as base path by Makers to write their results
var BasePath string

// Maker generates a Go output file
type Maker interface {
	Make(generatedOutput *io.Content, currentOutput *io.Content) (*io.Content, error)
	OutputFilepath() string
}

// Base represents common members for any maker
type Base struct {
	TypeHolder *parser.TypeHolder
}

// SetTypeHolder sets the holder for the type this maker uses
func (b *Base) SetTypeHolder(t *parser.TypeHolder) {
	b.TypeHolder = t
}

// Get returns the maker related with the template ID
func Get(template string) (Maker, error) {
	if registeredMakers == nil {
		return nil, errs.ErrNoMakerRegistered
	}

	templateID := templateIdentifier(template)
	maker, ok := registeredMakers[templateID]
	if !ok {
		return nil, errs.NewErrNotFound(fmt.Sprintf("Maker with id '%v'", templateID))
	}

	return maker, nil
}

func templateIdentifier(templateAbsPath string) string {
	filename := filepath.Base(templateAbsPath)
	var extension = filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}
