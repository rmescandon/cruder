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
	"os"

	check "gopkg.in/check.v1"
)

const badFormattedContent = `
package t

import "fmt"

const a = 16

func aFunction() { 
	fmt.Println("An instruction") }

`

const wellFormattedContent = `package t

import "fmt"

const a = 16

func aFunction() {
	fmt.Println("An instruction")
}
`

type PostprocessSuite struct{}

var _ = check.Suite(&PostprocessSuite{})

func (s *PostprocessSuite) SetUpTest(c *check.C) {}

func (s *PostprocessSuite) TestCheckFmt(c *check.C) {
	f, err := ioutil.TempFile("", "")
	c.Assert(err, check.IsNil)

	_, err = f.WriteString(badFormattedContent)
	c.Assert(err, check.IsNil)

	err = f.Close()
	c.Assert(err, check.IsNil)

	err = gofmt(f.Name())
	c.Assert(err, check.IsNil)

	f, err = os.Open(f.Name())
	c.Assert(err, check.IsNil)

	content, err := ioutil.ReadAll(f)
	c.Assert(err, check.IsNil)

	c.Assert(wellFormattedContent, check.Equals, string(content))
}
