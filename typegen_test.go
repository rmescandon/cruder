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
	"strings"
	"testing"

	check "gopkg.in/check.v1"
)

// gopkg.in/check.v1 stuff
func Test(t *testing.T) { check.TestingT(t) }

type TypeGenSuite struct{}

var _ = check.Suite(&TypeGenSuite{})

func (s *TypeGenSuite) SetUpTest(c *check.C) {}

func (s *TypeGenSuite) TestParseContent(c *check.C) {
	r := strings.NewReader(testTypeFileContent)
	_, file, err := parse(r)

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
	f, err := testTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(f, check.NotNil)

	b, syntaxTree, err := fileToSyntaxTree(f.Name())
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Not(check.HasLen), 0)
	c.Assert(syntaxTree, check.NotNil)

	source := &goFile{
		Path:    f.Name(),
		Content: b,
		Ast:     syntaxTree,
	}

	theMap, err := composeTypesMaps(source)
	c.Assert(err, check.IsNil)

	c.Assert(len(theMap), check.Equals, 1)
	for structName := range theMap {
		c.Assert(structName, check.Equals, "MyType")

		theFields := theMap[structName]
		for _, field := range theFields {
			switch field.Name {
			case "ID":
				c.Assert(field.Type, check.Equals, "int")
			case "Name":
				c.Assert(field.Type, check.Equals, "string")
			case "Description":
				c.Assert(field.Type, check.Equals, "string")
			case "SubTypes":
				c.Assert(field.Type, check.Equals, "[]string")
			default:
				c.Fail()
			}
		}
	}
}
