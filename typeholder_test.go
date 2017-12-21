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
	check "gopkg.in/check.v1"
)

type TypeHolderSuite struct {
	typeHolder TypeHolder
}

var _ = check.Suite(&TypeHolderSuite{})

func (s *TypeHolderSuite) SetUpTest(c *check.C) {
	s.typeHolder = TypeHolder{
		Name:        "MyType",
		IDFieldName: "ID",
		IDFieldType: "int",
		Fields: []typeField{
			typeField{Name: "ID", Type: "int"},
			typeField{Name: "Field1", Type: "string"},
			typeField{Name: "Field2", Type: "decimal"},
			typeField{Name: "Field3", Type: "int"},
		},
	}
}

func (s *TypeHolderSuite) TestTypeIdentifier(c *check.C) {
	c.Assert(s.typeHolder.typeIdentifier(), check.Equals, "myType")
}

func (s *TypeHolderSuite) TestTypeInComments(c *check.C) {
	c.Assert(s.typeHolder.typeInComments(), check.Equals, "mytype")
}

func (s *TypeHolderSuite) TestTypeFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.typeFieldsEnum(), check.Equals, "myType.Field1, myType.Field2, myType.Field3")
}

func (s *TypeHolderSuite) TestTypeRefFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.typeRefFieldsEnum(), check.Equals, "&myType.Field1, &myType.Field2, &myType.Field3")
}

func (s *TypeHolderSuite) TestTypeDBIDField(c *check.C) {
	c.Assert(s.typeHolder.typeDbIDField(), check.Equals, "id serial primary key not null,")
}

func (s *TypeHolderSuite) TestTypeDbFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.typeDbFieldsEnum(), check.Equals, "field1 varchar(200),\nfield2 decimal,\nfield3 integer")
}
