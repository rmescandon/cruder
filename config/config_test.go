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

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	tst "testing"

	check "gopkg.in/check.v1"
)

const settingsTestContent = `
version: 1.2-3

templates: /local/path/templates
`

type ConfigSuite struct{}

var _ = check.Suite(&ConfigSuite{})

// Test rewrites testing in a suite
func Test(t *tst.T) { check.TestingT(t) }

func (s *ConfigSuite) SetUpTest(c *check.C) {
	// Reset GOPATH to prevent taking cruder project as target project when
	// calculating projectURL
	os.Setenv("GOPATH", "")
}

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
	c.Assert(Config.Settings, check.Equals, f.Name())
	c.Assert(Config.Output, check.Equals, "")
	c.Assert(Config.Args.TypesFile, check.Equals, "")
	c.Assert(len(Config.Verbose), check.Equals, 0)
	c.Assert(Config.Version, check.Equals, "1.2-3")
	c.Assert(Config.TemplatesPath, check.Equals, "/local/path/templates")
}

func (s *ConfigSuite) TestProjectURL(c *check.C) {
	c.Assert(Config.setDefaultValuesWhenNeeded(), check.IsNil)

	template := "_#PROJECT#_"

	result := Config.ReplaceInTemplate(template)
	c.Assert(result, check.Equals, defaultProjectURL)

	Config.ProjectURL = "launchpad.com/myuser/myproject"

	result = Config.ReplaceInTemplate(template)
	c.Assert(result, check.Equals, "launchpad.com/myuser/myproject")
}

func (s *ConfigSuite) TestSetDefaultValues(c *check.C) {
	Config = Options{}

	c.Assert(Config.setDefaultValuesWhenNeeded(), check.IsNil)
	curr, err := os.Getwd()
	c.Assert(err, check.IsNil)

	c.Assert(Config.Output, check.Equals, curr)
	c.Assert(Config.Settings, check.Equals, filepath.Join(curr, defaultSettingsFile))
	c.Assert(Config.ProjectURL, check.Equals, defaultProjectURL)
	c.Assert(Config.APIVersion, check.Equals, defaultAPIVersion)
}

func (s *ConfigSuite) TestCalculateProjectURL(c *check.C) {
	projectPath := "theserver.com/theuser/theproject"

	d, err := ioutil.TempDir("", "")
	c.Assert(err, check.IsNil)

	os.Setenv("GOPATH", d)
	dst := filepath.Join(d, "src", projectPath)

	err = os.MkdirAll(dst, 0755)
	c.Assert(err, check.IsNil)

	err = os.Chdir(dst)
	c.Assert(err, check.IsNil)

	url := calculateProjectURL()
	c.Assert(url, check.Equals, projectPath)
}

func (s *ConfigSuite) TestCalculateProjectURL_notGopathSubdir(c *check.C) {
	projectPath := "theserver.com/theuser/theproject"

	d, err := ioutil.TempDir("", "")
	c.Assert(err, check.IsNil)

	os.Setenv("GOPATH", d)
	dst := filepath.Join(d, "src", projectPath)

	err = os.MkdirAll(dst, 0755)
	c.Assert(err, check.IsNil)

	url := calculateProjectURL()
	c.Assert(url, check.Equals, defaultProjectURL)
}
