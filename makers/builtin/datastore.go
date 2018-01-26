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
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
	"github.com/rmescandon/cruder/makers"
)

// Datastore generates datastore/<type>.go output go file
type Datastore struct {
	makers.BaseMaker
}

// ID returns 'datastore' as this maker identifier
func (ds *Datastore) ID() string {
	return "datastore"
}

// OutputFilepath returns the path to generated file
func (ds *Datastore) OutputFilepath() string {
	return filepath.Join(
		config.Config.Output,
		ds.ID(),
		strings.ToLower(ds.TypeHolder.Identifier())+".go")
}

// Make generates the results
func (ds *Datastore) Make() error {
	addOriginalType := false

	// check if output file exists
	_, err := os.Stat(ds.OutputFilepath())
	if err == nil {
		// in case if does exist, it should match the types file. Otherwise it's an error
		if ds.OutputFilepath() != ds.TypeHolder.Source.Path {
			return makers.NewErrOutputExists(ds.OutputFilepath())
		}

		// if output file is the same as types one, add the type to the generated output
		addOriginalType = true
	} else {
		// create needed dirs to outputPath
		ensureDir(filepath.Dir(ds.OutputFilepath()))
	}

	// execute the replacement
	logging.Debugf("Loadig template: %v", filepath.Base(ds.Template))
	templateContent, err := io.FileToString(ds.Template)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := ds.TypeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", ds.TypeHolder.Name, filepath.Base(ds.Template))
	}

	f, err := os.Create(ds.OutputFilepath())
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", ds.OutputFilepath(), err)
	}
	defer f.Close()

	_, err = f.WriteString(replacedStr)
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", ds.OutputFilepath(), err)
	}

	logging.Infof("Generated: %v", ds.OutputFilepath())

	// TODO IMPLEMENT if  addOriginalType....
	if addOriginalType {
		// - reload result file to AST
		// - prepend GenType AST to it
		// - write out AST to output, overwriting
		outputBytes, err := io.FileToByteArray(ds.OutputFilepath())
		if err != nil {
			return err
		}
		outputAst, err := io.ByteArrayToAST(outputBytes)
		if err != nil {
			return err
		}

		// insert GenType AST just before first function
		foundFirstFunc := false
		for i, decl := range outputAst.Decls {
			switch decl.(type) {
			case *ast.FuncDecl:
				outputAst.Decls = append(outputAst.Decls[:i], append([]ast.Decl{ds.TypeHolder.Decl}, outputAst.Decls[i:]...)...)
				foundFirstFunc = true
			}

			if foundFirstFunc {
				break
			}
		}

		err = io.ASTToFile(outputAst, ds.OutputFilepath())
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	makers.Register(&Datastore{})
}
