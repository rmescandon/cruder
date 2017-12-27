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

	check "gopkg.in/check.v1"
)

const settingsTestContent = `
version: 1.2-3

templates: /local/path/templates
`

type ConfigSuite struct{}

var _ = check.Suite(&ConfigSuite{})

func (s *ConfigSuite) SetUpTest(c *check.C) {}

func (s *ConfigSuite) TestLoadConfig(c *check.C) {

	f, err := ioutil.TempFile("", "")
	c.Assert(err, check.IsNil)

	_, err = f.WriteString(settingsTestContent)
	c.Assert(err, check.IsNil)

	err = f.Close()
	c.Assert(err, check.IsNil)

	Config.Settings = f.Name()
	err = Config.loadSettings()
	c.Assert(err, check.IsNil)
	c.Assert(Config.Version, check.Equals, "1.2-3")
	c.Assert(Config.TemplatesPath, check.Equals, "/local/path/templates")
}
