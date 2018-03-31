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

package engine

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/parser"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

const (
	mockContent = `package pkg

	import (
		"one"
		"two"
		"three"
		"fmt"
	)
	
	type TheStruct struct {
		AnyValue string
		OtherValue int
	}
	
	var variable TheStruct
	
	func theFunction() {
		err := do()
		if err != nil {
			fmt.Printf("error parsing parameters: %v\r\n", err)
			return
		}
	}
	`

	mockName       = "mock"
	mock2Name      = "mock2"
	mockOutputpath = "mock/path"
)

type mockMaker struct {
	id       string
	basePath string
}

func (m *mockMaker) ID() string {
	return m.id
}

func (m *mockMaker) OutputFilepath() string {
	if len(m.basePath) == 0 {
		var err error
		m.basePath, err = ioutil.TempDir("", "cruder_test")
		if err != nil {
			log.Error(err)
			return ""
		}
	}
	return filepath.Join(m.basePath, mockOutputpath)
}

func (m *mockMaker) Make(g *io.Content, c *io.Content) (*io.Content, error) {
	return io.NewContent(mockContent)
}

func (m *mockMaker) SetTypeHolder(*parser.TypeHolder) {}

func newMockMaker(id string) *mockMaker {
	return &mockMaker{id: id}
}

type EngineSuite struct{}

var _ = check.Suite(&EngineSuite{})

// Test rewrites testing in a suite
func Test(t *testing.T) { check.TestingT(t) }

func (s *EngineSuite) TestMerge(c *check.C) {
	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.TemplatesPath = "../testdata/templates/"
	io.NormalizePath(&config.Config.TemplatesPath)
	templates, err := availableTemplates()
	c.Assert(err, check.IsNil)
	c.Assert(templates, check.HasLen, 8)

	config.Config.ProjectURL = "server.dom/namespace/project"
	config.Config.APIVersion = "v1.0"

	for _, t := range templates {
		str, err := merge(h, t)
		c.Assert(err, check.IsNil)
		c.Assert(strings.Contains(str, "_#"), check.Equals, false)
		c.Assert(strings.Contains(str, "#_"), check.Equals, false)
	}
}

func (s *EngineSuite) TestMerge_cannotReplace(c *check.C) {
	f, err := ioutil.TempFile("", "")
	c.Assert(f, check.NotNil)
	c.Assert(err, check.IsNil)

	defer f.Close()

	_, err = f.WriteString("_#")
	c.Assert(err, check.IsNil)

	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	_, err = merge(h, f.Name())
	c.Assert(err, check.ErrorMatches, ".*type did not replace all.*template symbols")
}

func (s *EngineSuite) TestMerge_cannotReadTemplate(c *check.C) {
	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	_, err = merge(h, "/tmp/randominventedfilename")
	c.Assert(err, check.ErrorMatches, "Error reading template file.*")
}

func (s *EngineSuite) TestProcessMaker(c *check.C) {
	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	makers.Register(&mockMaker{id: mockName})

	t, err := testdata.TestTemplate(mockName)
	c.Assert(err, check.IsNil)

	c.Assert(processMaker(h, t), check.IsNil)
}

func (s *EngineSuite) TestProcessMakers(c *check.C) {
	h, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	makers.Register(&mockMaker{id: mockName})
	makers.Register(&mockMaker{id: mock2Name})

	t, err := testdata.TestTemplate(mockName)
	c.Assert(err, check.IsNil)
	t2, err := testdata.TestTemplate(mock2Name)
	c.Assert(err, check.IsNil)
	templates := []string{t, t2}

	processMakers([]*parser.TypeHolder{h}, templates)
}
