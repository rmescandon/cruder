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

package makers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
)

// Router generates service/router.go output go file
type Router struct {
	BaseMaker
}

// ID returns 'router' as this maker identifier
func (r *Router) ID() string {
	return "router"
}

// OutputFilepath returns the path to generated file
func (r *Router) OutputFilepath() string {
	return filepath.Join(config.Config.Output, "/service/router.go")
}

// Make generates the results
func (r *Router) Make() error {
	// execute the replacement
	logging.Debugf("Loadig template: %v", filepath.Base(r.Template))
	templateContent, err := io.FileToString(r.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := r.TypeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v",
			r.TypeHolder.Name, filepath.Base(r.Template))
	}

	replacedStr, err = config.Config.ReplaceInTemplate(replacedStr)
	if err != nil {
		return fmt.Errorf("Error replacing configuration over template %v",
			filepath.Base(r.Template))
	}

	// check if output file exists
	_, err = os.Stat(r.OutputFilepath())
	if err == nil {
		r.mergeExistingOutput(replacedStr)
	} else {
		// write out generated ast
		// create needed dirs to outputPath
		ensureDir(filepath.Dir(r.OutputFilepath()))

		io.StringToFile(replacedStr, r.OutputFilepath())

		logging.Infof("Generated: %v", r.OutputFilepath())
	}

	return nil
}

func (r *Router) mergeExistingOutput(replacedStr string) {
	//TODO IMPLEMENT
}

func init() {
	register(&Router{})
}
