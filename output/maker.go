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
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
)

// Maker generates a Go output file
type Maker interface {
	Make() error
	OutputFilepath() string
}

// BasicMaker represents common members for a a maker
type BasicMaker struct {
	TypeHolder *src.TypeHolder
	Output     *io.GoFile
	Template   string
}

// NewMaker returns a maker for a certain type and template
func NewMaker(typeHolder *src.TypeHolder, templatePath string) (Maker, error) {
	var outputPath string

	templateID := templateIdentifier(templatePath)
	outputPath = createOutputPath(config.Config.Output, templateID, strings.ToLower(typeHolder.Name))

	switch templateID {
	case "datastore":
		return &Datastore{
			BasicMaker{
				TypeHolder: typeHolder,
				Output: &io.GoFile{
					Path: outputPath,
				},
				Template: templatePath,
			},
		}, nil
	case "db":
		return &Db{
			BasicMaker{
				TypeHolder: typeHolder,
				Output: &io.GoFile{
					Path: outputPath,
				},
				Template: templatePath,
			},
		}, nil
	case "handler":
		return &Handler{
			BasicMaker{
				TypeHolder: typeHolder,
				Output: &io.GoFile{
					Path: outputPath,
				},
				Template: templatePath,
			},
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
		return filepath.Join(outputFolder, "datastore/db.go")
	case "datastore":
		return filepath.Join(outputFolder, "datastore", typeIdentifierInLower+".go")
	case "handler":
		return filepath.Join(outputFolder, "service", typeIdentifierInLower+".go")
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
