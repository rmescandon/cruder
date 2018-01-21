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
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
)

// Handler makes the controller
type Handler struct {
	BasicMaker
}

// OutputFilepath returns the path to the generated file
func (h *Handler) OutputFilepath() string {
	return h.Output.Path
}

// Make generates the output
func (h *Handler) Make() error {
	// check if output file exists
	_, err := os.Stat(h.Output.Path)
	if err == nil {
		return NewErrOutputExists(h.Output.Path)
	}

	ensureDir(filepath.Dir(h.Output.Path))

	logging.Debugf("Loadig template: %v", filepath.Base(h.Template))
	templateContent, err := io.FileToString(h.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := h.TypeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", h.TypeHolder.Name, filepath.Base(h.Template))
	}

	replacedStr, err = config.Config.ReplaceInTemplate(replacedStr)
	if err != nil {
		return fmt.Errorf("Error replacing configuration over template %v", filepath.Base(h.Template))
	}

	f, err := os.Create(h.Output.Path)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", h.Output.Path, err)
	}
	defer f.Close()

	_, err = f.WriteString(replacedStr)
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", h.Output.Path, err)
	}

	logging.Infof("Generated: %v", h.Output.Path)
	return nil
}
