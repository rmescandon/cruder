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
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/src"
)

// Datastore generates datastore/<type>.go output go file
type Datastore struct {
	Type     *src.TypeHolder
	File     *io.GoFile
	Template string
}

// OutputFilepath returns the path to generated file
func (ds *Datastore) OutputFilepath() string {
	return ds.File.Path
}

// Run runs to generate the result
func (ds *Datastore) Run() error {
	// don't write if file exists
	// FIXME this should not happen if output file is same as source one.
	// IN such case, original file types should be added to output
	_, err := os.Stat(ds.File.Path)
	if err == nil {
		return fmt.Errorf("File %v already exists. Skip writting", ds.File.Path)
	}

	// execute the replacement:
	// create needed dirs to outputPath
	ensureDir(filepath.Dir(ds.File.Path))

	logging.Debugf("Loadig template: %v", filepath.Base(ds.Template))
	templateContent, err := io.FileContentAsString(ds.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := ds.Type.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", ds.Type.Name, filepath.Base(ds.Template))
	}

	f, err := os.Create(ds.File.Path)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", ds.File.Path, err)
	}
	defer f.Close()

	_, err = f.WriteString(replacedStr)
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", ds.File.Path, err)
	}

	logging.Infof("Generated: %v", ds.File.Path)
	return nil
}
