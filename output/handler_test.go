// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
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
	"path/filepath"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
	"github.com/rmescandon/cruder/testdata"
	check "gopkg.in/check.v1"
)

type HandlerSuite struct{}

var _ = check.Suite(&HandlerSuite{})

func (s *HandlerSuite) SetUpSuite(c *check.C) {}
func (s *HandlerSuite) SetUpTest(c *check.C)  {}

func (s *HandlerSuite) TestMakeHandler(c *check.C) {
	config.Config.ProjectURL = "github.com/auser/aproject"
	//--------------------------------------------------------------------------
	// 1.- Create an output file for MyType, not having a previous existing file
	typeFile, err := testdata.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	source, err := io.NewGoFile(typeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err := src.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	handler := &Handler{
		BasicMaker{
			TypeHolder: typeHolders[0],
			Output: &io.GoFile{
				Path: filepath.Join(config.Config.Output, "handlertestoutput.go"),
			},
			Template: "../testdata/templates/handler.template",
		},
	}

	c.Assert(handler.Make(), check.IsNil)

	content, err := io.FileToString(handler.Output.Path)
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)

	// -----------------------------------------------------------------------
	// 2.- Reset typeHolders and load now OtherType. Create the output and see
	// if maker returns ErrOutputExists error
	otherTypeFile, err := testdata.TestOtherTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(otherTypeFile, check.NotNil)

	source, err = io.NewGoFile(otherTypeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err = src.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	handler.TypeHolder = typeHolders[0]

	c.Assert(handler.Make(), check.DeepEquals, NewErrOutputExists(handler.Output.Path))
}
