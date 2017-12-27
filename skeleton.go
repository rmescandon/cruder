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

package cruder

import (
	"fmt"
	"path/filepath"
)

// GenerateSkeletonCode generates the skeleton code based on loaded configuration and available templates
func GenerateSkeletonCode() error {
	Log.Debug("Generating Skeleton Code...")

	Log.Debugf("searching for available templates at %v", Config.TemplatesPath)
	availableTemplates, err := filepath.Glob(filepath.Join(Config.TemplatesPath, "*.template"))
	if err != nil {
		return fmt.Errorf("Error listing available templates: %v", err)
	}

	for _, template := range availableTemplates {
		Log.Debugf("Found template: %v", filepath.Base(template))

		typeHolders, err := typeHoldersFromFile(Config.TypesFile)
		if err != nil {
			return fmt.Errorf("Error composing type holders from types file: %v", err)
		}

		err = replace(template, typeHolders)
		if err != nil {
			return err
		}
	}

	return nil
}
