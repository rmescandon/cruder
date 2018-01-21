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

import check "gopkg.in/check.v1"

type ErrorSuite struct{}

var _ = check.Suite(&ErrorSuite{})

func (s *ErrorSuite) TestErrOutputExistsMessage(c *check.C) {
	err := NewErrOutputExists("/any/random/path")
	c.Assert(err.Error(), check.Equals, "File /any/random/path already exists. Skip writting")
}
