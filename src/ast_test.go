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

package src

import (
	check "gopkg.in/check.v1"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/testdata"
)

// gopkg.in/check.v1 stuff
type AstSuite struct{}

var _ = check.Suite(&AstSuite{})

func (s *AstSuite) SetUpTest(c *check.C) {}

func (s *AstSuite) TestFileWithSimpeType(c *check.C) {
	f, err := testdata.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(f, check.NotNil)

	b, err := io.FileToByteArray(f.Name())
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Not(check.HasLen), 0)

	ast, err := io.ByteArrayToAST(b)
	c.Assert(ast, check.NotNil)
	c.Assert(err, check.IsNil)

	source := &io.GoFile{
		Path:    f.Name(),
		Content: b,
		Ast:     ast,
	}

	typeHolders, err := ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)
	c.Assert(typeHolders[0].Name, check.Equals, "MyType")

	theFields := typeHolders[0].Fields
	c.Assert(len(theFields), check.Equals, 4)
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
