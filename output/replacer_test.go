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
	"io/ioutil"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
	"github.com/rmescandon/cruder/testing"

	check "gopkg.in/check.v1"
)

type ReplacerSuite struct{}

var _ = check.Suite(&ReplacerSuite{})

func (s *ReplacerSuite) SetUpTest(c *check.C) {}

func (s *ReplacerSuite) TestMakersForType(c *check.C) {
	typeFile, err := testing.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	source, err := io.NewGoFile(typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(source, check.NotNil)

	typeHolders, err := src.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)

	makers, err := makersForType(typeHolders[0])
	c.Assert(err, check.IsNil)
	// TODO increase when having more makers ready
	c.Assert(makers, check.HasLen, 1)
}

// TODO this test should disappear when having specific test for every stage
func (s *ReplacerSuite) TestReplaceInAllTemplates(c *check.C) {
	typeFile, err := testing.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	config.Config.TemplatesPath = "testdata/templates"

	source, err := io.NewGoFile(typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(source, check.NotNil)

	typeHolders, err := src.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)
	c.Assert(typeHolders[0].Name, check.Equals, "MyType")
	c.Assert(typeHolders[0].Fields, check.HasLen, 4)
	c.Assert(typeHolders[0].IDField.Name, check.Equals, "ID")
	c.Assert(typeHolders[0].IDField.Type, check.Equals, "int")
	c.Assert(typeHolders[0].Source, check.NotNil)
	c.Assert(typeHolders[0].Source.Path, check.Equals, typeFile.Name())

	b, err := io.FileToByteArray(typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Not(check.HasLen), 0)

	ast, err := io.ByteArrayToAST(b)
	c.Assert(ast, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(typeHolders[0].Source.Content, check.DeepEquals, b)
	c.Assert(typeHolders[0].Source.Ast, check.NotNil)
	c.Assert(typeHolders[0].Source.Ast, check.DeepEquals, ast)

	makers, err := makersForType(typeHolders[0])
	c.Assert(err, check.IsNil)
	// TODO increase when having more makers ready
	c.Assert(makers, check.HasLen, 1)

	for _, maker := range makers {
		err = maker.Run()
		c.Assert(err, check.IsNil)

		content, err := io.FileToString(maker.OutputFilepath())
		c.Assert(err, check.IsNil)
		c.Assert(strings.Contains(content, "_#"), check.Equals, false)
		c.Assert(strings.Contains(content, "#_"), check.Equals, false)
	}
}
