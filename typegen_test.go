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

package cruder

import (
	"go/ast"
	"io/ioutil"
	"strings"
	"testing"

	check "gopkg.in/check.v1"
)

const (
	testTypeFileContent = `
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

// gopkg.in/check.v1 stuff
func Test(t *testing.T) { check.TestingT(t) }

type TypeGenSuite struct{}

var _ = check.Suite(&TypeGenSuite{})

func (s *TypeGenSuite) SetUpTest(c *check.C) {}

func (s *TypeGenSuite) TestParseContent(c *check.C) {
	r := strings.NewReader(testTypeFileContent)
	file, err := parse(r)

	c.Assert(err, check.IsNil)
	c.Assert(len(file.Decls), check.Equals, 1)

	genDecl := file.Decls[0].(*ast.GenDecl)
	c.Assert(len(genDecl.Specs), check.Equals, 1)

	typeDecl := genDecl.Specs[0].(*ast.TypeSpec)
	ident := typeDecl.Name
	c.Assert(ident.Name, check.Equals, "MyType")

	structType := typeDecl.Type.(*ast.StructType)
	c.Assert(len(structType.Fields.List), check.Equals, 4)
}

func (s *TypeGenSuite) TestFileWithSimpeType(c *check.C) {

	f, err := ioutil.TempFile("", "")
	c.Assert(err, check.IsNil)

	_, err = f.WriteString(testTypeFileContent)
	c.Assert(err, check.IsNil)

	theMap, err := GetTypesMaps(f.Name())
	c.Assert(err, check.IsNil)

	c.Assert(len(theMap), check.Equals, 1)
	for structName := range theMap {
		c.Assert(structName, check.Equals, "MyType")

		theFields := theMap[structName]
		for field := range theFields {
			switch field {
			case "ID":
				c.Assert(theFields[field], check.Equals, "int")
			case "Name":
				c.Assert(theFields[field], check.Equals, "string")
			case "Description":
				c.Assert(theFields[field], check.Equals, "string")
			case "Subtypes":
				c.Assert(theFields[field], check.Equals, "[]string")
			default:
				c.Fail()
			}
		}
	}
}
