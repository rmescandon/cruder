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

package engine

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

type EngineSuite struct {
	typeFile    *os.File
	typeHolders []*parser.TypeHolder
	templates   []string
}

var _ = check.Suite(&EngineSuite{})

// Test rewrites testing in a suite
func Test(t *testing.T) { check.TestingT(t) }

func (s *EngineSuite) TestMerge(c *check.C) {
	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.TemplatesPath = "../testdata/templates/"
	io.NormalizePath(&config.Config.TemplatesPath)
	templates, err := availableTemplates()
	c.Assert(err, check.IsNil)
	c.Assert(templates, check.HasLen, 8)

	for _, t := range templates {
		str, err := merge(h, t)
		c.Assert(err, check.IsNil)
		c.Assert(strings.Contains(str, "_#"), check.Equals, false)
		c.Assert(strings.Contains(str, "#_"), check.Equals, false)
	}
}

func (s *EngineSuite) TestMerge_cannotReplace(c *check.C) {
	f, err := ioutil.TempFile("", "")
	c.Assert(f, check.NotNil)
	c.Assert(err, check.IsNil)

	defer f.Close()

	_, err = f.WriteString("_#")
	c.Assert(err, check.IsNil)

	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	_, err = merge(h, f.Name())
	c.Assert(err, check.ErrorMatches, ".*type did not replace all.*template symbols")
}

/* TEMPORARY COMMENTED OUT
func (s *EngineSuite) SetUpTest(c *check.C) {
	var err error
	s.typeFile, err = testdata.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(s.typeFile, check.NotNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	config.Config.TemplatesPath = "../testdata/templates"

	source, err := io.NewGoFile(s.typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(source, check.NotNil)

	s.typeHolders, err = parser.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)

	s.templates, err = availableTemplates()
	c.Assert(err, check.IsNil)
	c.Assert(s.templates, check.HasLen, 5)
}

func (s *EngineSuite) TestGetMakers(c *check.C) {

	makers, err := buildMakers(s.typeHolders, s.templates)
	c.Assert(err, check.IsNil)
	// TODO increase when having more makers ready
	c.Assert(makers, check.HasLen, 5)
}

// TODO this test should disappear when having specific test for every stage
func (s *EngineSuite) TestReplaceInAllTemplates(c *check.C) {
	c.Assert(s.typeHolders, check.HasLen, 1)
	c.Assert(s.typeHolders[0].Name, check.Equals, "MyType")
	c.Assert(s.typeHolders[0].Fields, check.HasLen, 4)
	c.Assert(s.typeHolders[0].IDFieldName(), check.Equals, "ID")
	c.Assert(s.typeHolders[0].IDFieldType(), check.Equals, "int")
	c.Assert(s.typeHolders[0].Source, check.NotNil)
	c.Assert(s.typeHolders[0].Source.Path, check.Equals, s.typeFile.Name())

	b, err := io.FileToByteArray(s.typeFile.Name())
	c.Assert(err, check.IsNil)
	c.Assert(b, check.Not(check.HasLen), 0)

	ast, err := io.ByteArrayToAST(b)
	c.Assert(ast, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(s.typeHolders[0].Source.Content, check.DeepEquals, b)
	c.Assert(s.typeHolders[0].Source.Ast, check.NotNil)
	c.Assert(s.typeHolders[0].Source.Ast, check.DeepEquals, ast)

	makers, err := buildMakers(s.typeHolders, s.templates)
	c.Assert(err, check.IsNil)
	// TODO increase when having more makers ready
	c.Assert(makers, check.HasLen, 5)

	for _, maker := range makers {
		err = maker.Make()
		c.Assert(err, check.IsNil)

		content, err := io.FileToString(maker.OutputFilepath())

		c.Assert(err, check.IsNil)
		c.Assert(strings.Contains(content, "_#"), check.Equals, false)
		c.Assert(strings.Contains(content, "#_"), check.Equals, false)
	}
}
*/
