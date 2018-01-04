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

package src

import (
	check "gopkg.in/check.v1"
)

type TypeHolderSuite struct {
	typeHolder TypeHolder
}

var _ = check.Suite(&TypeHolderSuite{})

func (s *TypeHolderSuite) SetUpTest(c *check.C) {
	s.typeHolder = TypeHolder{
		Name: "MyType",
		IDField: TypeField{
			Name: "ID",
			Type: "int",
		},
		Fields: []TypeField{
			TypeField{Name: "ID", Type: "int"},
			TypeField{Name: "Field1", Type: "string"},
			TypeField{Name: "Field2", Type: "decimal"},
			TypeField{Name: "Field3", Type: "int"},
		},
	}
}

func (s *TypeHolderSuite) TestTypeIdentifier(c *check.C) {
	c.Assert(s.typeHolder.Identifier(), check.Equals, "myType")
}

func (s *TypeHolderSuite) TestTypeInComments(c *check.C) {
	c.Assert(s.typeHolder.InComments(), check.Equals, "mytype")
}

func (s *TypeHolderSuite) TestTypeFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.FieldsEnum(), check.Equals, "myType.Field1, myType.Field2, myType.Field3")
}

func (s *TypeHolderSuite) TestTypeRefFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.RefFieldsEnum(), check.Equals, "&myType.Field1, &myType.Field2, &myType.Field3")
}

func (s *TypeHolderSuite) TestTypeDBIDField(c *check.C) {
	c.Assert(s.typeHolder.DbIDField(), check.Equals, "id serial primary key not null,")
}

func (s *TypeHolderSuite) TestTypeDbFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.DbFieldsEnum(), check.Equals, "field1 varchar(200),\nfield2 decimal,\nfield3 integer")
}
