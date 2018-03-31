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

package io

import (
	"os"
	"path/filepath"

	check "gopkg.in/check.v1"
)

const (
	fakePath = "a/random/path"
)

type DirSuite struct{}

var _ = check.Suite(&DirSuite{})

func (s *DirSuite) SetUpSuite(c *check.C) {
	os.Remove(filepath.Join(os.TempDir(), fakePath))
}

func (s *DirSuite) TestEnsureDir(c *check.C) {
	path := filepath.Join(os.TempDir(), fakePath)
	_, err := os.Stat(path)
	c.Assert(os.IsNotExist(err), check.Equals, true)

	err = EnsureDir(path)
	c.Assert(err, check.IsNil)

	_, err = os.Stat(path)
	c.Assert(err, check.IsNil)
}

func (s *DirSuite) TestNormalizePath_noChange(c *check.C) {
	path := "/any/path"
	err := NormalizePath(&path)
	c.Assert(err, check.IsNil)
	c.Assert(path, check.Equals, "/any/path")
}

func (s *DirSuite) TestNormalizePath_endingSlash(c *check.C) {
	path := "/any/path/"
	err := NormalizePath(&path)
	c.Assert(err, check.IsNil)
	c.Assert(path, check.Equals, "/any/path")
}

func (s *DirSuite) TestNormalizePath_home(c *check.C) {
	path := "~/any/path"
	err := NormalizePath(&path)
	c.Assert(err, check.IsNil)
	c.Assert(path, check.Equals, filepath.Join(os.Getenv("HOME"), "any/path"))
}

func (s *DirSuite) TestNormalizePath_relativePath(c *check.C) {
	path := "any/path"
	err := NormalizePath(&path)
	c.Assert(err, check.IsNil)
	current, err := os.Getwd()
	c.Assert(err, check.IsNil)
	c.Assert(path, check.Equals, filepath.Join(current, "any/path"))
}
