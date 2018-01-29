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

package makers

import (
	"io/ioutil"
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/parser"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

type DbSuite struct{}

var _ = check.Suite(&DbSuite{})

func (s *DbSuite) TestMakeDb(c *check.C) {
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

	db := &Db{
		BaseMaker{
			TypeHolder: typeHolders[0],
			Template:   "../testdata/templates/db.template",
		},
	}

	c.Assert(db.Make(), check.IsNil)

	content, err := io.FileToString(db.OutputFilepath())
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)

	// -----------------------------------------------------------------------
	// 2.- Reset typeHolders and load now OtherType. Create the output and see
	// if both MyType and OtherType are included in Datastore interface
	otherTypeFile, err := testdata.TestOtherTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(otherTypeFile, check.NotNil)

	source, err = io.NewGoFile(otherTypeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err = parser.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	db.TypeHolder = typeHolders[0]

	c.Assert(db.Make(), check.IsNil)

	content, err = io.FileToString(db.OutputFilepath())

	// Verify here if both types are included in output
	c.Assert(strings.Contains(content, "CreateMyTypeTable() error"), check.Equals, true)
	c.Assert(strings.Contains(content, "ListMyTypes() ([]MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "GetMyType(ID int) (MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "FindMyType(query string) (MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "CreateMyType(myType MyType) (int, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "UpdateMyType(ID int, myType MyType)"), check.Equals, true)
	c.Assert(strings.Contains(content, "DeleteMyType(ID int) error"), check.Equals, true)

	c.Assert(strings.Contains(content, "CreateMyOtherTypeTable() error"), check.Equals, true)
	c.Assert(strings.Contains(content, "ListMyOtherTypes() ([]MyOtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "GetMyOtherType(AnID int) (MyOtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "FindMyOtherType(query string) (MyOtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "CreateMyOtherType(myOtherType MyOtherType) (int, error)"), check.Equals, true)
	c.Assert(strings.Contains(content, "UpdateMyOtherType(AnID int, myOtherType MyOtherType)"), check.Equals, true)
	c.Assert(strings.Contains(content, "DeleteMyOtherType(AnID int) error"), check.Equals, true)

	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(content, "_#"), check.Equals, false)
	c.Assert(strings.Contains(content, "#_"), check.Equals, false)
}
