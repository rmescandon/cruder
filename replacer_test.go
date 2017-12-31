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
	"io/ioutil"
	"strings"

	check "gopkg.in/check.v1"
)

type ReplacerSuite struct{}

var _ = check.Suite(&ReplacerSuite{})

func (s *ReplacerSuite) SetUpTest(c *check.C) {}

func (s *ReplacerSuite) TestReplaceUsingDatastoreTemplate(c *check.C) {
	typeFile, err := testTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	Config.TemplatesPath = "testdata/templates"

	typeHolders, err := typeHoldersFromFile(typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)
	c.Assert(typeHolders[0].Name, check.Equals, "MyType")
	c.Assert(typeHolders[0].Fields, check.HasLen, 4)
	c.Assert(typeHolders[0].IDField.Name, check.Equals, "ID")
	c.Assert(typeHolders[0].IDField.Type, check.Equals, "int")
	c.Assert(typeHolders[0].Source, check.NotNil)
	c.Assert(typeHolders[0].Source.Path, check.Equals, typeFile.Name())
	b, ast, err := fileToSyntaxTree(typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders[0].Source.Content, check.DeepEquals, b)
	c.Assert(typeHolders[0].Source.Ast, check.NotNil)
	c.Assert(typeHolders[0].Source.Ast, check.DeepEquals, ast)
	c.Assert(typeHolders[0].Outputs, check.IsNil)

	outputFile, err := typeHolders[0].getOutputFilePathFor(Datastore)
	c.Assert(err, check.IsNil)

	err = typeHolders[0].appendOutputs()
	c.Assert(err, check.IsNil)

	content, err := fileContentsAsString(outputFile)
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)
}
