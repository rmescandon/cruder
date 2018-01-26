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

package builtin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/parser"
)

// Db maker to include types in datastore interface
type Db struct {
	makers.BaseMaker
}

// ID returns 'db'
func (db *Db) ID() string {
	return "db"
}

// OutputFilepath returns the path to generated file
func (db *Db) OutputFilepath() string {
	return filepath.Join(config.Config.Output, fmt.Sprintf("datastore/%v.go", db.ID()))
}

// Make generates the results
func (db *Db) Make() error {
	// execute the replacement
	logging.Debugf("Loadig template: %v", filepath.Base(db.Template))
	templateContent, err := io.FileToString(db.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := db.TypeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", db.TypeHolder.Name, filepath.Base(db.Template))
	}

	// check if output file exists
	_, err = os.Stat(db.OutputFilepath())
	if err == nil {
		db.mergeExistingOutput(replacedStr)
	} else {
		// write out generated ast
		// create needed dirs to outputPath
		ensureDir(filepath.Dir(db.OutputFilepath()))

		io.StringToFile(replacedStr, db.OutputFilepath())

		logging.Infof("Generated: %v", db.OutputFilepath())
	}

	return nil
}

// mergeExistingOutput resolves the conflict when already exists an output file
func (db *Db) mergeExistingOutput(replacedStr string) error {
	logging.Infof("Merging new type into: %v", db.OutputFilepath())
	generatedAst, err := io.ByteArrayToAST([]byte(replacedStr))
	if err != nil {
		return err
	}

	// load current output
	content, err := io.FileToByteArray(db.OutputFilepath())
	if err != nil {
		return err
	}

	currentAst, err := io.ByteArrayToAST(content)
	if err != nil {
		return err
	}

	generatedIface := parser.GetInterface(generatedAst, "Datastore")
	currentIface := parser.GetInterface(currentAst, "Datastore")

	// search for generatedIface methods into currentIface and add them if not found
	for _, method := range parser.GetInterfaceMethods(generatedIface) {
		if !parser.HasMethod(currentIface, method.Names[0].Name) {
			parser.AddMethod(currentIface, method)
		}
	}

	// write out the resultant modified Datastore interface to output
	// TODO VERIFY that using pointers is enough to alter generatedAst before writting out
	io.ASTToFile(currentAst, db.OutputFilepath())
	logging.Infof("Merged into: %v successfully", db.OutputFilepath())

	return nil
}

func init() {
	makers.Register(&Db{})
}
