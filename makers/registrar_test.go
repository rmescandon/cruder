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
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/parser"
	check "gopkg.in/check.v1"
)

type mockRegistrant struct {
	mockMaker
}

func (m *mockRegistrant) SetTypeHolder(t *parser.TypeHolder) {}

func newMockRegistrant(id, path, content string) *mockRegistrant {
	return &mockRegistrant{*newMockMaker(id, path, content)}
}

type RegistrarSuite struct{}

var _ = check.Suite(&RegistrarSuite{})

func (s *RegistrarSuite) SetUpTest(c *check.C) {
	registeredMakers = nil
}

func (s *RegistrarSuite) TestRegister(c *check.C) {
	err := Register(newMockRegistrant(mock1Name, mock1Outputpath, mockContent))
	c.Assert(err, check.IsNil)

	err = Register(newMockRegistrant(mock2Name, mock2Outputpath, mockContent))
	c.Assert(err, check.IsNil)

	err = Register(newMockRegistrant(mock3Name, mock3Outputpath, mockContent))
	c.Assert(err, check.IsNil)

	c.Assert(registeredMakers, check.HasLen, 3)
}

func (s *RegistrarSuite) TestRegister_duplicated(c *check.C) {
	err := Register(newMockRegistrant(mock1Name, mock1Outputpath, mockContent))
	c.Assert(err, check.IsNil)

	err = Register(newMockRegistrant(mock1Name, mock1Outputpath, mockContent))
	c.Assert(err, check.NotNil)
	switch err.(type) {
	case errs.ErrDuplicatedMaker:
	default:
		c.Fail()
	}

	c.Assert(registeredMakers, check.HasLen, 1)
}

func (s *RegistrarSuite) TestRegister_nilRegistrant(c *check.C) {
	err := Register(nil)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.Equals, errs.ErrNilObject)
}
