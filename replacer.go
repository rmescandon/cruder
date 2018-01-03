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
	"os"
	"path/filepath"
)

// GenerateSkeletonCode generates the skeleton code based on loaded configuration and available templates
func GenerateSkeletonCode() error {
	Log.Debug("Generating Skeleton Code...")

	source, err := newGoFile(Config.TypesFile)
	if err != nil {
		return fmt.Errorf("Error reading go source file: %v", err)
	}

	typeHolders, err := composeTypeHolders(source)
	if err != nil {
		return fmt.Errorf("Error composing type holders from types file: %v", err)
	}

	for _, typeHolder := range typeHolders {
		err = typeHolder.appendOutputs()
		if err != nil {
			Log.Warningf("Could not append output: %v", err)
			continue
		}
	}

	return nil
}

func templateIdentifier(templateAbsPath string) string {
	filename := filepath.Base(templateAbsPath)
	var extension = filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
