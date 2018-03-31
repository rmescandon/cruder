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

package makers

import (
	"testing"

	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
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

	mock1Name       = "mock1"
	mock2Name       = "mock2"
	mock3Name       = "mock3"
	mock1Outputpath = "mock1/path"
	mock2Outputpath = "mock2/path"
	mock3Outputpath = "mock3/path"
)

type mockMaker struct {
	Base
	id      string
	path    string
	content string
}

func (m *mockMaker) ID() string {
	return m.id
}

func (m *mockMaker) OutputFilepath() string {
	return m.path
}

func (m *mockMaker) Make(g *io.Content, c *io.Content) (*io.Content, error) {
	return io.NewContent(m.content)
}

func newMockMaker(id, path, content string) *mockMaker {
	return &mockMaker{id: id, path: path, content: content}
}

type MakerSuite struct{}

var _ = check.Suite(&MakerSuite{})

func Test(t *testing.T) { check.TestingT(t) }

func (s *MakerSuite) SetUpTest(c *check.C) {
	registeredMakers = map[string]Registrant{
		mock1Name: newMockMaker(mock1Name, mock1Outputpath, mockContent),
		mock2Name: newMockMaker(mock2Name, mock2Outputpath, mockContent),
		mock3Name: newMockMaker(mock3Name, mock3Outputpath, mockContent),
	}
}

func (s *MakerSuite) TestGetMaker(c *check.C) {
	m, err := Get("whatever/path/" + mock1Name + ".template")
	c.Assert(err, check.IsNil)
	c.Assert(m.ID(), check.Equals, mock1Name)
	c.Assert(m.OutputFilepath(), check.Equals, mock1Outputpath)

	m, err = Get("whatever/path/" + mock2Name + ".template")
	c.Assert(err, check.IsNil)
	c.Assert(m.ID(), check.Equals, mock2Name)
	c.Assert(m.OutputFilepath(), check.Equals, mock2Outputpath)

	m, err = Get("whatever/path/" + mock3Name + ".template")
	c.Assert(err, check.IsNil)
	c.Assert(m.ID(), check.Equals, mock3Name)
	c.Assert(m.OutputFilepath(), check.Equals, mock3Outputpath)
}

func (s *MakerSuite) TestGetMaker_notExistingTemplate(c *check.C) {
	m, err := Get("whatever/path/notexisting.template")
	c.Assert(err, check.NotNil)
	c.Assert(m, check.IsNil)
	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}

func (s *MakerSuite) TestGetMaker_emptyTemplatePath(c *check.C) {
	m, err := Get("")
	c.Assert(err, check.NotNil)
	c.Assert(m, check.IsNil)
	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}

func (s *MakerSuite) TestGetMaker_noRegisteredMakers(c *check.C) {
	registeredMakers = nil
	m, err := Get("whatever/path/" + mock1Name + ".template")
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, errs.ErrNoMakerRegistered)
	c.Assert(m, check.IsNil)
}

func (s *MakerSuite) TestTemplateIdentifier(c *check.C) {
	c.Assert(templateIdentifier("name.template"), check.Equals, "name")
	c.Assert(templateIdentifier("name.template.template"), check.Equals, "name.template")
	c.Assert(templateIdentifier("the/path/name.template"), check.Equals, "name")
	c.Assert(templateIdentifier("the/path/.template"), check.Equals, "")
	c.Assert(templateIdentifier("the/path/name.other"), check.Equals, "")
}
