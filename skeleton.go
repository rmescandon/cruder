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

import "path/filepath"

const (
	templatesPath = "templates"
)

// GenerateSkeletonCode generates the skeleton code based on loaded configuration and available templates
func GenerateSkeletonCode() error {
	Log.Debug("Generating Skeleton Code...")

	availableTemplates := []string{"datastore.template"}
	for _, template := range availableTemplates {
		templateFullPath := filepath.Join(templatesPath, template)
		err := replace(templateFullPath, Config.TypesFile)
		if err != nil {
			return err
		}
	}

	return nil
}
