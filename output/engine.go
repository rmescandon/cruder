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

package output

import (
	"fmt"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/src"
)

// GenerateSkeletonCode generates the skeleton code based on loaded configuration and available templates
func GenerateSkeletonCode() error {
	logging.Debug("Generating Skeleton Code...")

	source, err := io.NewGoFile(config.Config.TypesFile)
	if err != nil {
		return fmt.Errorf("Error reading go source file: %v", err)
	}

	typeHolders, err := src.ComposeTypeHolders(source)
	if err != nil {
		return fmt.Errorf("Error composing type holders from types file: %v", err)
	}

	templates, err := availableTemplates()
	if err != nil {
		return fmt.Errorf("Error listing available templates: %v", err)
	}

	makers, err := makers(typeHolders, templates)
	if err != nil {
		return err
	}

	for _, maker := range makers {
		err := maker.Make()
		if err != nil {
			logging.Warningf("Could not run maker: %v", err)
			continue
		}
	}

	return nil
}

func availableTemplates() ([]string, error) {
	logging.Debugf("searching for available templates at %v", config.Config.TemplatesPath)
	return filepath.Glob(filepath.Join(config.Config.TemplatesPath, "*.template"))
}

func makers(typeHolders []*src.TypeHolder, availableTemplates []string) ([]Maker, error) {
	var makers []Maker
	for _, template := range availableTemplates {
		logging.Debugf("Found template: %v", filepath.Base(template))
		for _, t := range typeHolders {
			maker, err := NewMaker(t, template)
			if err != nil {
				return []Maker{}, err
			}
			makers = append(makers, maker)
		}

	}
	return makers, nil
}
