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

package makers

import (
	"fmt"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
)

// CopyMaker base structf for makers that only copies input template to output file
type CopyMaker struct {
	BaseMaker
}

// Make does the copy
func (c *CopyMaker) copy(output string) error {
	if len(output) == 0 {
		return NewErrEmptyString("Output filepath")
	}

	logging.Debugf("Loadig template: %v", filepath.Base(c.Template))
	templateContent, err := io.FileToString(c.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	ensureDir(filepath.Dir(output))

	io.StringToFile(templateContent, output)

	logging.Infof("Generated: %v", output)
	return nil
}
