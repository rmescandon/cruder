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

package core

import (
	"fmt"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/parser"
)

// GenerateSkeletonCode generates the skeleton code based on loaded configuration and available templates
func GenerateSkeletonCode() error {

	log.Info("Generating Skeleton Code...")

	source, err := io.NewGoFile(config.Config.TypesFile)
	if err != nil {
		return fmt.Errorf("Error reading go source file: %v", err)
	}

	typeHolders, err := parser.ComposeTypeHolders(source)
	if err != nil {
		return fmt.Errorf("Error composing type holders from types file: %v", err)
	}

	templates, err := availableTemplates()
	if err != nil {
		return fmt.Errorf("Error listing available templates: %v", err)
	}

	makers, err := buildMakers(typeHolders, templates)
	if err != nil {
		return err
	}

	for _, maker := range makers {
		err := maker.Make()
		if err != nil {
			log.Warningf("Could not run maker: %v", err)
			continue
		}
	}

	return nil
}

func availableTemplates() ([]string, error) {
	log.Infof("searching for available templates at %v", config.Config.TemplatesPath)
	return filepath.Glob(filepath.Join(config.Config.TemplatesPath, "*.template"))
}

func buildMakers(holders []*parser.TypeHolder, templates []string) ([]makers.Maker, error) {
	var mks []makers.Maker
	for _, t := range templates {
		log.Infof("Found template: %v", filepath.Base(t))
		for _, h := range holders {
			//FIXME: this won't work as every maker associated with a type is reused for the next type
			// Execute Run for every maker got until next one or
			// Create dynamic objects by reflection into makers.Get
			m, err := makers.New(h, t)
			if err != nil {
				return []makers.Maker{}, err
			}
			mks = append(mks, m)
		}

	}
	return mks, nil
}
