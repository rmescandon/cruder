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

package decl

import (
	"strings"
	tst "testing"

	check "gopkg.in/check.v1"
)

type TypeHolderSuite struct {
	typeHolder TypeHolder
}

var _ = check.Suite(&TypeHolderSuite{})

// Test rewrites testing in a suite
func Test(t *tst.T) { check.TestingT(t) }

func (s *TypeHolderSuite) SetUpTest(c *check.C) {
	s.typeHolder = TypeHolder{
		Name: "MyType",
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
	c.Assert(strings.ToLower(s.typeHolder.Name), check.Equals, "mytype")
}

func (s *TypeHolderSuite) TestIDFieldName(c *check.C) {
	c.Assert(s.typeHolder.IDFieldName(), check.Equals, "ID")
}

func (s *TypeHolderSuite) TestIDFieldType(c *check.C) {
	c.Assert(s.typeHolder.IDFieldType(), check.Equals, "int")
}

func (s *TypeHolderSuite) TestFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.FieldsEnum(), check.Equals, "myType.Field1, myType.Field2, myType.Field3")
}

func (s *TypeHolderSuite) TestFieldsEnumRef(c *check.C) {
	c.Assert(s.typeHolder.FieldsEnumRef(), check.Equals, "&myType.Field1, &myType.Field2, &myType.Field3")
}

func (s *TypeHolderSuite) TestIDFieldInDDL(c *check.C) {
	c.Assert(s.typeHolder.IDFieldInDDL(), check.Equals, "id serial primary key not null,")
}

func (s *TypeHolderSuite) TestFieldsInDDL(c *check.C) {
	c.Assert(s.typeHolder.FieldsInDDL(), check.Equals, "field1 varchar(200),\nfield2 decimal,\nfield3 integer")
}

func (s *TypeHolderSuite) TestFieldsInDML(c *check.C) {
	c.Assert(s.typeHolder.FieldsInDML(), check.Equals, "Field1, Field2, Field3")
}

func (s *TypeHolderSuite) TestValuesInDMLParams(c *check.C) {
	c.Assert(s.typeHolder.ValuesInDMLParams(), check.Equals, "$1, $2, $3")
}

func (s *TypeHolderSuite) TestIDFieldAsDMLParam(c *check.C) {
	c.Assert(s.typeHolder.IDFieldAsDMLParam(), check.Equals, "ID=$4")
}

func (s *TypeHolderSuite) TestFieldsAsDMLParams(c *check.C) {
	c.Assert(s.typeHolder.FieldsAsDMLParams(), check.Equals, "Field1=$1, Field2=$2, Field3=$3")
}
