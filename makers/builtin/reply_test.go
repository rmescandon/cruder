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

package builtin

import (
	check "gopkg.in/check.v1"
)

type ReplySuite struct{}

var _ = check.Suite(&ReplySuite{})

/*
func (s *ReplySuite) TestCopyReply(c *check.C) {
	//--------------------------------------------------------------------------
	// 1.- Create an output file, not having a previous existing file
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

	r := &Reply{
		makers.CopyMaker{
			BaseMaker: makers.BaseMaker{
				TypeHolder: typeHolders[0],
				Template:   "../testdata/templates/reply.template",
			},
		},
	}

	c.Assert(r.Make(), check.IsNil)

	srcContent, err := io.FileToString(r.Template)
	c.Assert(err, check.IsNil)
	dstContent, err := io.FileToString(r.OutputFilepath())
	c.Assert(err, check.IsNil)
	c.Assert(srcContent, check.Equals, dstContent)

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

	err = r.Make()
	c.Assert(err, check.NotNil)
	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}
*/
