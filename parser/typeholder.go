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
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"github.com/rmescandon/cruder/io"
)

// Categories for output generated files
const (
	Datastore = iota + 1
	Handler
	Router
)

// TypeHolder holds a type previously read from file
type TypeHolder struct {
	Name   string
	Source *io.GoFile
	Fields []TypeField
	Decl   *ast.GenDecl
}

// TypeField holds a field in a type
type TypeField struct {
	Name string
	Type string
}

// Identifier returns type name in camel case, except first letter, which is lower case:
// "theType"
func (holder *TypeHolder) Identifier() string {
	if len(holder.Name) > 0 {
		return strings.ToLower(string(holder.Name[0])) + holder.Name[1:len(holder.Name)]
	}
	return ""
}

// FieldsEnum returns enum of the fields including type indentifier and field name:
// "theType.Field1, theType.Field2, theType.FieldN"
func (holder *TypeHolder) FieldsEnum() string {
	return holder.fieldsEnum(false)
}

// FieldsEnumRef returns enum of type fields, including type identifiera and field name reference:
// "&theType.Field1, &theType.Field2, &theType.FieldN"
func (holder *TypeHolder) FieldsEnumRef() string {
	return holder.fieldsEnum(true)
}

// depending on the bool param returns the same as typeRefFieldsEnum or typeFieldsEnum
func (holder *TypeHolder) fieldsEnum(asRef bool) string {
	ref := ""
	if asRef {
		ref = "&"
	}

	tokens := []string{}
	for _, field := range holder.Fields {
		// skip ID field
		if field.Name == holder.IDFieldName() {
			continue
		}

		token := ref + holder.identifierDotField(field.Name)
		tokens = append(tokens, token)
	}

	return strings.Join(tokens, ", ")
}

// IDFieldName returns the name of the first field, taken as ID
func (holder *TypeHolder) IDFieldName() string {
	return holder.Fields[0].Name
}

// IDFieldType returns the type of the first field, taken as ID
func (holder *TypeHolder) IDFieldType() string {
	return holder.Fields[0].Type
}

// FindFieldName return the name of the field used for searches
func (holder *TypeHolder) FindFieldName() string {
	return holder.Fields[1].Name
}

// identifierDotField returns type identifier plus dot plus parameter fieldname, like:
// "theType.Field1"
func (holder *TypeHolder) identifierDotField(fieldName string) string {
	return holder.Identifier() + "." + fieldName
}

// IDFieldInDDL returns the IDField as seen in SQL DDL operations
func (holder *TypeHolder) IDFieldInDDL() string {
	result := strings.ToLower(holder.IDFieldName()) + " "
	if holder.IDFieldType() == "int" {
		result = result + "integer "
	} else {
		result = result + holder.IDFieldType() + " "
	}

	result = result + "primary key not null,"
	return result
}

// FieldsInDDL returns the type fields as they are used for SQL DDL operations, like:
// "Field1 varchar(200),
// Field2 int,
// FieldN varchar(200)"
func (holder *TypeHolder) FieldsInDDL() string {
	tokens := []string{}
	for _, field := range holder.Fields {
		// skip the ID field, show the rest
		if field.Name == holder.IDFieldName() {
			continue
		}

		var t string
		switch field.Type {
		case "string":
			t = "varchar(200)"
		case "int":
			t = "integer"
		case "decimal":
			t = "decimal"
		case "bool":
			t = "boolean"
		case "time.Time":
			t = "timestamp"
			// FIXME fill the rest of the types
		}

		token := fmt.Sprintf("%v %v", strings.ToLower(field.Name), t)
		tokens = append(tokens, token)
	}

	return strings.Join(tokens, ",\n")
}

// FieldsInDML returns "field1, field2, field3"
func (holder *TypeHolder) FieldsInDML() string {
	tokens := []string{}
	for _, field := range holder.Fields {
		// skip the ID field, show the rest
		if field.Name == holder.IDFieldName() {
			continue
		}

		tokens = append(tokens, strings.ToLower(field.Name))
	}

	return strings.Join(tokens, ", ")
}

// ValuesInDMLParams returns something like "$1, $2, $3"
func (holder *TypeHolder) ValuesInDMLParams() string {
	tokens := []string{}
	for i := 1; i < len(holder.Fields); i++ {
		tokens = append(tokens, "$"+strconv.Itoa(i))
	}
	return strings.Join(tokens, ", ")
}

// IDFieldAsDMLParam returns something like "id=$4"
func (holder *TypeHolder) IDFieldAsDMLParam() string {
	return holder.IDFieldName() + "=$" + strconv.Itoa(len(holder.Fields))
}

// FieldsAsDMLParams returns something like "field1=$1, field2=$2, field3=$3"
func (holder *TypeHolder) FieldsAsDMLParams() string {
	tokens := []string{}
	for i, field := range holder.Fields {
		// skip the ID field, show the rest
		if field.Name == holder.IDFieldName() {
			continue
		}

		token := field.Name + "=$" + strconv.Itoa(i)
		tokens = append(tokens, token)
	}
	return strings.Join(tokens, ", ")
}

// IDFieldTypeParse returns the type parsing instruction for ID field
func (holder *TypeHolder) IDFieldTypeParse() string {
	switch holder.IDFieldType() {
	case "int":
		return "strconv.Atoi(vars[\"" + holder.IDFieldName() + "\"])"
	case "decimal":
		return "strconv.ParseFloat(vars[\"" + holder.IDFieldName() + "\"])"
	case "bool":
		return "strconv.ParseBool(vars[\"" + holder.IDFieldName() + "\"])"
	default:
		return "vars[\"" + holder.IDFieldName() + "\"]"
	}
}

// IDFieldTypeFormat returns the type formatting instruction to have ID as string
func (holder *TypeHolder) IDFieldTypeFormat() string {
	switch holder.IDFieldType() {
	case "int":
		return "strconv.Itoa(vars[\"" + holder.IDFieldName() + "\"])"
	case "decimal":
		return "strconf.FormatFloat(vars[\"" + holder.IDFieldName() + "\"])"
	case "bool":
		return "strconf.FormatBool(vars[\"" + holder.IDFieldName() + "\"])"
	default:
		return "vars[\"" + holder.IDFieldName() + "\"]"
	}
}

// IDFieldPattern returns the pattern associated with id field type, to be used when routing REST paths
func (holder *TypeHolder) IDFieldPattern() string {
	switch holder.IDFieldType() {
	case "int":
		return "[0-9]+"
	case "decimal":
		return "^[0-9]+(\\.[0-9]{1,2})?$"
	case "bool":
		return "^(?:tru|fals)e$"
	default:
		return "[a-zA-Z0-9-_\\.]+"
	}
}

// ReplaceInTemplate replaces template marks with holder data
func (holder *TypeHolder) ReplaceInTemplate(templateContent string) (string, error) {
	replaced := templateContent

	/*
		For a type like:

		type TheType struct {
			ID          int
			Name        string
			Description string
			SubTypes    []string
		}

		here are the replacements meaning:
	*/

	// TheType
	replaced = strings.Replace(replaced, "_#TYPE#_", holder.Name, -1)

	// theType
	replaced = strings.Replace(replaced, "_#TYPE.IDENTIFIER#_", holder.Identifier(), -1)

	// thetype
	replaced = strings.Replace(replaced, "_#TYPE.LOWERCASE#_", strings.ToLower(holder.Name), -1)

	// ID
	replaced = strings.Replace(replaced, "_#ID.FIELD.NAME#_", holder.IDFieldName(), -1)

	// id
	replaced = strings.Replace(replaced, "_#ID.FIELD.NAME.LOWERCASE#_", strings.ToLower(holder.IDFieldName()), -1)

	// int
	replaced = strings.Replace(replaced, "_#ID.FIELD.TYPE#_", holder.IDFieldType(), -1)

	// name
	replaced = strings.Replace(replaced, "_#FIND.FIELD.NAME#_", holder.FindFieldName(), -1)

	// theType.Name, theType.Description, theType.Subtypes
	replaced = strings.Replace(replaced, "_#FIELDS.ENUM#_", holder.FieldsEnum(), -1)

	// &theType.Name, &theType.Description, &theType.Subtypes
	replaced = strings.Replace(replaced, "_#FIELDS.ENUM.REF#_", holder.FieldsEnumRef(), -1)

	// id integer primary_key not null
	replaced = strings.Replace(replaced, "_#ID.FIELD.DDL#_", holder.IDFieldInDDL(), -1)

	// Name			varchar,
	// Description	varchar,
	// Subtypes		varchar
	replaced = strings.Replace(replaced, "_#FIELDS.DDL#_", holder.FieldsInDDL(), -1)

	// name, description, subtypes
	replaced = strings.Replace(replaced, "_#FIELDS.DML#_", holder.FieldsInDML(), -1)

	// $1, $2, $3
	replaced = strings.Replace(replaced, "_#VALUES.DML.PARAMS#_", holder.ValuesInDMLParams(), -1)

	// id=$4
	replaced = strings.Replace(replaced, "_#ID.FIELD.DML.PARAM#_", holder.IDFieldAsDMLParam(), -1)

	// name=$1, description=$2, subtypes=$3
	replaced = strings.Replace(replaced, "_#FIELDS.DML.PARAMS#_", holder.FieldsAsDMLParams(), -1)

	// strconv.Atoi(vars["id"])
	replaced = strings.Replace(replaced, "_#ID.FIELD.TYPE.PARSE#_", holder.IDFieldTypeParse(), -1)

	// strconv.Itoa(vars["id"])
	replaced = strings.Replace(replaced, "_#ID.FIELD.TYPE.FORMAT#_", holder.IDFieldTypeFormat(), -1)

	// [a-z]+
	replaced = strings.Replace(replaced, "_#ID.FIELD.PATTERN#_", holder.IDFieldPattern(), -1)

	return replaced, nil
}
