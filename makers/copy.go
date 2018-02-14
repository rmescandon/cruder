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

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
)

// CopyMaker base struct for makers that only copies input template to output file
type CopyMaker struct {
	BaseMaker
}

// Copy does the copy of raw input template to output
func (c *CopyMaker) Copy(output string) error {
	if len(output) == 0 {
		return errs.NewErrEmptyString("Output filepath")
	}

	log.Debugf("Loadig template: %v", filepath.Base(c.Template))
	templateContent, err := io.FileToString(c.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	io.EnsureDir(filepath.Dir(output))

	io.StringToFile(templateContent, output)

	log.Infof("Generated: %v", output)
	return nil
}
