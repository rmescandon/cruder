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

package builtin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/makers"
)

// Handler makes the controller
type Handler struct {
	makers.BaseMaker
}

// ID returns the identifier 'handler' for this maker
func (h *Handler) ID() string {
	return "handler"
}

// OutputFilepath returns the path to the generated file
func (h *Handler) OutputFilepath() string {
	return filepath.Join(
		config.Config.Output,
		h.ID(),
		strings.ToLower(h.TypeHolder.Identifier())+".go")
}

// Make generates the output
func (h *Handler) Make() error {
	// check if output file exists
	_, err := os.Stat(h.OutputFilepath())
	if err == nil {
		return makers.NewErrOutputExists(h.OutputFilepath())
	}

	ensureDir(filepath.Dir(h.OutputFilepath()))

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

	f, err := os.Create(h.OutputFilepath())
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", h.OutputFilepath(), err)
	}
	defer f.Close()

	_, err = f.WriteString(replacedStr)
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", h.OutputFilepath(), err)
	}

	logging.Infof("Generated: %v", h.OutputFilepath())
	return nil
}

func init() {
	makers.Register(&Handler{})
}
