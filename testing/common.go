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

package testing

import (
	"io/ioutil"
	"os"
	"testing"

	check "gopkg.in/check.v1"
)

const (
	TestTypeFileContent = `
	package mytype
	
	// MyType test type to generate skeletom code
	type MyType struct {
		ID          int
		Name        string
		Description string
		SubTypes    []string
	}	
	`
)

// Test rewrites testing in a suite
func Test(t *testing.T) { check.TestingT(t) }

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
