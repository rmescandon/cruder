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

package testdata

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
)

// TestTypeFileContent testing type
const (
	TestTypeFileContent = `
	package mytype
	
	// MyType test type to generate skeletom code
	type MyType struct {
		ID            int
		Name          string
		Description   string
		TheBoolThing  bool
		TheFloatThing float
	}	
	`

	TestOtherTypeFileContent = `
	package myothertype
	
	// MyOtherType test type to generate skeletom code
	type MyOtherType struct {
		AnID           int
		AName          string
		ADescription   string
		ABoolThing     bool
		AFloatingThing float
	}	
	`

	TestTemplateContent = `
	package pkg

	func anyFunction() error {
		if err := Db.Do_#TYPE#_Thing(); err != nil {
			return err
		}

		return nil
	}

	`
)

// TestTemplate returns a file
func TestTemplate(id string) (string, error) {
	dir, err := ioutil.TempDir(os.TempDir(), "")
	if err != nil {
		return "", err
	}

	f, err := os.Create(filepath.Join(dir, id+".template"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(TestTemplateContent)
	return f.Name(), err
}

// TestTypeFile returns a temporary file with a test type into it
func TestTypeFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return f, err
	}
	defer f.Close()

	_, err = f.WriteString(TestTypeFileContent)
	return f, err
}

// TestOtherTypeFile returns a temporary file with another different test type into it
func TestOtherTypeFile() (*os.File, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return f, err
	}
	defer f.Close()

	_, err = f.WriteString(TestOtherTypeFileContent)
	return f, err
}

// TestTypeHolder returns a type holder for testing purposes
func TestTypeHolder() (*parser.TypeHolder, error) {
	typeFile, err := TestTypeFile()
	if err != nil {
		return nil, err
	}

	source, err := io.NewGoFile(typeFile.Name())
	if err != nil {
		return nil, err
	}

	typeHolders, err := parser.ComposeTypeHolders(source)
	if err != nil {
		return nil, err
	}

	if len(typeHolders) != 1 {
		return nil, errors.New("Generated a number of type holders different than 1")
	}

	return typeHolders[0], nil
}
