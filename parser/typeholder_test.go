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

package parser

import (
	"strings"
	"testing"

	check "gopkg.in/check.v1"
)

type TypeHolderSuite struct {
	typeHolder      TypeHolder
	emptyTypeHolder TypeHolder
}

var _ = check.Suite(&TypeHolderSuite{})

// Test rewrites testing in a suite
func Test(t *testing.T) { check.TestingT(t) }

func (s *TypeHolderSuite) SetUpSuite(c *check.C) {
	s.typeHolder = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "int"},
			{Name: "Field1", Type: "string"},
			{Name: "Field2", Type: "decimal"},
			{Name: "Field3", Type: "int"},
		},
	}
	s.emptyTypeHolder = TypeHolder{}
}

func (s *TypeHolderSuite) TestIdentifier(c *check.C) {
	c.Assert(s.typeHolder.Identifier(), check.Equals, "myType")
}

func (s *TypeHolderSuite) TestIdentifier_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.Identifier(), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldName(c *check.C) {
	c.Assert(s.typeHolder.IDFieldName(), check.Equals, "ID")
}

func (s *TypeHolderSuite) TestIDFieldName_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.IDFieldName(), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldType(c *check.C) {
	c.Assert(s.typeHolder.IDFieldType(), check.Equals, "int")
}

func (s *TypeHolderSuite) TestIDFieldType_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.IDFieldType(), check.Equals, "")
}

func (s *TypeHolderSuite) TestFindFieldName(c *check.C) {
	c.Assert(s.typeHolder.FindFieldName(), check.Equals, "Field1")
}

func (s *TypeHolderSuite) TestFindFieldName_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FindFieldName(), check.Equals, "")
}

func (s *TypeHolderSuite) TestFindFieldName_oneField(c *check.C) {
	t := TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "int"},
		},
	}
	c.Assert(t.FindFieldName(), check.Equals, "ID")
}

func (s *TypeHolderSuite) TestFieldsEnum(c *check.C) {
	c.Assert(s.typeHolder.FieldsEnum(), check.Equals, "myType.Field1, myType.Field2, myType.Field3")
}

func (s *TypeHolderSuite) TestFieldsEnum_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FieldsEnum(), check.Equals, "")
}

func (s *TypeHolderSuite) TestFieldsEnumRef(c *check.C) {
	c.Assert(s.typeHolder.FieldsEnumRef(), check.Equals, "&myType.Field1, &myType.Field2, &myType.Field3")
}

func (s *TypeHolderSuite) TestFieldsEnumRef_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FieldsEnumRef(), check.Equals, "")
}

func (s *TypeHolderSuite) TestTypeInComments(c *check.C) {
	c.Assert(strings.ToLower(s.typeHolder.Name), check.Equals, "mytype")
}

func (s *TypeHolderSuite) TestTypeInComments_empty(c *check.C) {
	c.Assert(strings.ToLower(s.emptyTypeHolder.Name), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldInDDL(c *check.C) {
	c.Assert(s.typeHolder.IDFieldInDDL(), check.Equals, "id integer primary key not null,")
}

func (s *TypeHolderSuite) TestIDFieldInDDL_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.IDFieldInDDL(), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldInDDL_otherType(c *check.C) {
	t := TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "otherID", Type: "string"},
		},
	}
	c.Assert(t.IDFieldInDDL(), check.Equals, "otherid varchar primary key not null,")
}

func (s *TypeHolderSuite) TestDDLTypeConversion(c *check.C) {
	c.Assert(ddlType("string"), check.Equals, "varchar")
	c.Assert(ddlType("int"), check.Equals, "integer")
	c.Assert(ddlType("decimal"), check.Equals, "decimal")
	c.Assert(ddlType("bool"), check.Equals, "boolean")
	c.Assert(ddlType("time.Time"), check.Equals, "timestamp")
	c.Assert(ddlType("other"), check.Equals, "")
}

func (s *TypeHolderSuite) TestFieldsInDDL(c *check.C) {
	c.Assert(s.typeHolder.FieldsInDDL(), check.Equals, "field1 varchar,\nfield2 decimal,\nfield3 integer")
}

func (s *TypeHolderSuite) TestFieldsInDDL_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FieldsInDDL(), check.Equals, "")
}

func (s *TypeHolderSuite) TestFieldsInDML(c *check.C) {
	c.Assert(s.typeHolder.FieldsInDML(), check.Equals, "field1, field2, field3")
}

func (s *TypeHolderSuite) TestFieldsInDML_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FieldsInDML(), check.Equals, "")
}

func (s *TypeHolderSuite) TestValuesInDMLParams(c *check.C) {
	c.Assert(s.typeHolder.ValuesInDMLParams(), check.Equals, "$1, $2, $3")
}

func (s *TypeHolderSuite) TestValuesInDMLParams_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.ValuesInDMLParams(), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldAsDMLParam(c *check.C) {
	c.Assert(s.typeHolder.IDFieldAsDMLParam(), check.Equals, "id=$4")
}

func (s *TypeHolderSuite) TestIDFieldAsDMLParam_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.IDFieldAsDMLParam(), check.Equals, "")
}

func (s *TypeHolderSuite) TestFieldsAsDMLParams(c *check.C) {
	c.Assert(s.typeHolder.FieldsAsDMLParams(), check.Equals, "field1=$1, field2=$2, field3=$3")
}

func (s *TypeHolderSuite) TestFieldsAsDMLParams_empty(c *check.C) {
	c.Assert(s.emptyTypeHolder.FieldsAsDMLParams(), check.Equals, "")
}

func (s *TypeHolderSuite) TestIDFieldTypeParse(c *check.C) {
	c.Assert(s.typeHolder.IDFieldTypeParse(),
		check.Equals,
		"strconv.Atoi(vars[\""+strings.ToLower(s.typeHolder.IDFieldName())+"\"])")

	t := TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "decimal"},
		},
	}
	c.Assert(t.IDFieldTypeParse(),
		check.Equals,
		"strconv.ParseFloat(vars[\""+strings.ToLower(t.IDFieldName())+"\"])")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "bool"},
		},
	}
	c.Assert(t.IDFieldTypeParse(),
		check.Equals,
		"strconv.ParseBool(vars[\""+strings.ToLower(t.IDFieldName())+"\"])")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "other"},
		},
	}
	c.Assert(t.IDFieldTypeParse(),
		check.Equals,
		"vars[\""+strings.ToLower(t.IDFieldName())+"\"]")

}

func (s *TypeHolderSuite) TestIDFieldTypeFormat(c *check.C) {
	c.Assert(s.typeHolder.IDFieldTypeFormat(),
		check.Equals,
		"strconv.Itoa("+strings.ToLower(s.typeHolder.IDFieldName())+")")

	t := TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "decimal"},
		},
	}
	c.Assert(t.IDFieldTypeFormat(),
		check.Equals,
		"strconv.FormatFloat("+strings.ToLower(t.IDFieldName())+")")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "bool"},
		},
	}
	c.Assert(t.IDFieldTypeFormat(),
		check.Equals,
		"strconv.FormatBool("+strings.ToLower(t.IDFieldName())+")")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "other"},
		},
	}
	c.Assert(t.IDFieldTypeFormat(),
		check.Equals,
		strings.ToLower(t.IDFieldName()))
}

func (s *TypeHolderSuite) TestIDFieldPattern(c *check.C) {
	c.Assert(s.typeHolder.IDFieldPattern(), check.Equals, "[0-9]+")

	t := TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "decimal"},
		},
	}
	c.Assert(t.IDFieldPattern(),
		check.Equals,
		"^[0-9]+(\\.[0-9]{1,2})?$")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "bool"},
		},
	}
	c.Assert(t.IDFieldPattern(),
		check.Equals,
		"^(?:tru|fals)e$")

	t = TypeHolder{
		Name: "MyType",
		Fields: []TypeField{
			{Name: "ID", Type: "other"},
		},
	}
	c.Assert(t.IDFieldPattern(),
		check.Equals,
		"[a-zA-Z0-9-_\\.]+")
}

func (s *TypeHolderSuite) TestReplaceInTemplate(c *check.C) {
	c.Assert(s.typeHolder.ReplaceInTemplate("_#TYPE#_"), check.Equals, "MyType")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#TYPE.IDENTIFIER#_"), check.Equals, "myType")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#TYPE.LOWERCASE#_"), check.Equals, "mytype")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.NAME#_"), check.Equals, "ID")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.NAME.LOWERCASE#_"), check.Equals, "id")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.TYPE#_"), check.Equals, "int")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIND.FIELD.NAME#_"), check.Equals, "Field1")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIELDS.ENUM#_"), check.Equals, "myType.Field1, myType.Field2, myType.Field3")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIELDS.ENUM.REF#_"), check.Equals, "&myType.Field1, &myType.Field2, &myType.Field3")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.DDL#_"), check.Equals, "id integer primary key not null,")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIELDS.DDL#_"), check.Equals, "field1 varchar,\nfield2 decimal,\nfield3 integer")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIELDS.DML#_"), check.Equals, "field1, field2, field3")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#VALUES.DML.PARAMS#_"), check.Equals, "$1, $2, $3")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.DML.PARAM#_"), check.Equals, "id=$4")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#FIELDS.DML.PARAMS#_"), check.Equals, "field1=$1, field2=$2, field3=$3")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.TYPE.PARSE#_"), check.Equals, "strconv.Atoi(vars[\"id\"])")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.TYPE.FORMAT#_"), check.Equals, "strconv.Itoa(id)")
	c.Assert(s.typeHolder.ReplaceInTemplate("_#ID.FIELD.PATTERN#_"), check.Equals, "[0-9]+")
}
