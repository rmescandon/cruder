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

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/src"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

type DbSuite struct{}

var _ = check.Suite(&DbSuite{})

func (s *DbSuite) TestDbLoading(c *check.C) {
	typeFile, err := testdata.TestTypeFile()
	c.Assert(err, check.IsNil)
	c.Assert(typeFile, check.NotNil)

	source, err := io.NewGoFile(typeFile.Name())
	c.Assert(err, check.IsNil)

	typeHolders, err := src.ComposeTypeHolders(source)
	c.Assert(err, check.IsNil)
	c.Assert(typeHolders, check.HasLen, 1)

	outputFile, err := ioutil.TempFile("", "cruder_")
	c.Assert(err, check.IsNil)

	db := &Db{TypeHolders: typeHolders,
		File: &io.GoFile{
			Path: outputFile.Name(),
		},
		Template: "../testdata/templates/db.template",
	}

	c.Assert(db.Make(), check.IsNil)
}
