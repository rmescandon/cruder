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

package output

import (
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
)

// Maker generates a Go output file
type Maker interface {
	Run() error
	OutputFilepath() string
}

// NewMaker returns a maker for a certain type and template
func NewMaker(holders []*src.TypeHolder, outputFolder, templatePath string) (Maker, error) {
	var outputPath string
	templateID := templateIdentifier(templatePath)
	switch templateID {
	case "datastore":
		if len(holders) > 0 {
			outputPath = createOutputPath(outputFolder, "datastore", holders[0].InComments())
		}

		return &Datastore{
			TypeHolders: holders,
			File: &io.GoFile{
				Path: outputPath,
			},
			Template: templatePath,
		}, nil
	}

	return nil, nil
}

func templateIdentifier(templateAbsPath string) string {
	filename := filepath.Base(templateAbsPath)
	var extension = filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

func createOutputPath(outputFolder, templateID, typeIdentifierInLower string) string {
	switch templateID {
	case "db":
		fallthrough
	case "datastore":
		return filepath.Join(outputFolder, "datastore", typeIdentifierInLower+".go")
	}

	return ""
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
