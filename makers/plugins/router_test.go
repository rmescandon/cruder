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

package main

import (
	"io/ioutil"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

type RouterSuite struct{}

var _ = check.Suite(&RouterSuite{})

func (s *RouterSuite) TestMakeRouter(c *check.C) {
	//--------------------------------------------------------------------------
	// 1.- Create an output file for MyType, not having a previous existing file
	typeFile, err := testdata.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	source, err := io.NewGoFile(typeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err := parser.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	r := &Router{
		BaseMaker{
			TypeHolder: typeHolders[0],
			Template:   "../testdata/templates/router.template",
		},
	}

	c.Assert(r.Make(), check.IsNil)

	content, err := io.FileToString(r.OutputFilepath())
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)

	// -----------------------------------------------------------------------
	// 2.- Reset typeHolders and load now OtherType. Create the output and see
	// if both MyType and OtherType are included into
	otherTypeFile, err := testdata.TestOtherTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(otherTypeFile, check.NotNil)

	source, err = io.NewGoFile(otherTypeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err = parser.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	r.TypeHolder = typeHolders[0]

	c.Assert(r.Make(), check.IsNil)

	content, err = io.FileToString(r.OutputFilepath())

	// Verify here if both types are included in output
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"mytype\"), http.HandlerFunc(CreateMyType)).Methods(\"POST\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"mytype\"), http.HandlerFunc(ListMyTypes)).Methods(\"GET\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"mytype/{ID:[0-9]+}\"), http.HandlerFunc(GetMyType)).Methods(\"GET\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"mytype/{ID:[0-9]+}\"), http.HandlerFunc(UpdateMyType)).Methods(\"PUT\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"mytype/{ID:[0-9]+}\"), http.HandlerFunc(DeleteMyType)).Methods(\"DELETE\")"), check.Equals, true)

	c.Assert(strings.Contains(content, "router.Handle(composePath(\"myothertype\"), http.HandlerFunc(CreateMyOtherType)).Methods(\"POST\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"myothertype\"), http.HandlerFunc(ListMyOtherTypes)).Methods(\"GET\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"myothertype/{AnID:[0-9]+}\"), http.HandlerFunc(GetMyOtherType)).Methods(\"GET\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"myothertype/{AnID:[0-9]+}\"), http.HandlerFunc(UpdateMyOtherType)).Methods(\"PUT\")"), check.Equals, true)
	c.Assert(strings.Contains(content, "router.Handle(composePath(\"myothertype/{AnID:[0-9]+}\"), http.HandlerFunc(DeleteMyOtherType)).Methods(\"DELETE\")"), check.Equals, true)

	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)
}
