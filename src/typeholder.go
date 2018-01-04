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
	"go/ast"
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
	Name    string
	Source  *io.GoFile
	IDField TypeField
	Fields  []TypeField
	Decl    *ast.GenDecl
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

// InComments returns type name in lower case:
// "thetype"
func (holder *TypeHolder) InComments() string {
	return strings.ToLower(holder.Name)
}

// FieldsEnum returns enum of the fields including type indentifier and field name:
// "theType.Field1, theType.Field2, theType.FieldN"
func (holder *TypeHolder) FieldsEnum() string {
	return holder.fieldsEnum(false)
}

// RefFieldsEnum returns enum of type fields, including type identifiera and field name reference:
// "&theType.Field1, &theType.Field2, &theType.FieldN"
func (holder *TypeHolder) RefFieldsEnum() string {
	return holder.fieldsEnum(true)
}

// depending on the bool param returns the same as typeRefFieldsEnum or typeFieldsEnum
func (holder *TypeHolder) fieldsEnum(asRef bool) string {
	ref := ""
	if asRef {
		ref = "&"
	}

	var enum string
	for _, field := range holder.Fields {
		// skip ID field
		if field.Name == holder.IDField.Name {
			continue
		}

		if len(enum) == 0 {
			enum = ref + holder.IdentifierDotField(field.Name)
			continue
		}

		enum = enum + ", " + ref + holder.IdentifierDotField(field.Name)
	}
	return enum
}

// IDFieldName returns the name of the IDField
func (holder *TypeHolder) IDFieldName() string {
	return holder.IdentifierDotField(holder.IDField.Name)
}

// IDFieldType returns the type of the IDField
func (holder *TypeHolder) IDFieldType() string {
	return holder.IdentifierDotField(holder.IDField.Type)
}

// IdentifierDotField returns type identifier plus dot plus parameter fields name, like:
// "theType.Field1"
func (holder *TypeHolder) IdentifierDotField(fieldName string) string {
	return holder.Identifier() + "." + fieldName
}

// DbIDField returns the IDField as seen in SQL operations
func (holder *TypeHolder) DbIDField() string {
	result := strings.ToLower(holder.IDField.Name) + " "
	if holder.IDField.Type == "int" {
		result = result + "serial "
	} else {
		result = result + holder.IDField.Type + " "
	}

	result = result + "primary key not null,"
	return result
}

// DbFieldsEnum returns the type fields as they are used for SQL operations, like:
// "Field1 varchar(200),
// Field2 int,
// FieldN varchar(200)"
func (holder *TypeHolder) DbFieldsEnum() string {
	var result string
	for _, field := range holder.Fields {
		if field.Name == holder.IDField.Name {
			continue
		}

		if len(result) > 0 {
			result = result + "\n"
		}
		result = result + strings.ToLower(field.Name) + " "

		switch field.Type {
		case "string":
			result = result + "varchar(200),"
		case "int":
			result = result + "integer,"
		case "decimal":
			result = result + "decimal,"
		case "bool":
			result = result + "boolean,"
		case "time.Time":
			result = result + "timestamp,"
			// FIXME fill the rest of the types
		}
	}

	// Remove last element if it is a colon
	lenResult := len(result)
	if lenResult > 0 && result[lenResult-1] == ',' {
		result = result[0 : lenResult-1]
	}

	return result
}

// ReplaceInTemplate replaces template marks with holder data
func (holder *TypeHolder) ReplaceInTemplate(templateContent string) (string, error) {
	replaced := templateContent

	replaced = strings.Replace(replaced, "_#TheType#_", holder.Name, -1)
	replaced = strings.Replace(replaced, "_#theType#_", holder.Identifier(), -1)
	replaced = strings.Replace(replaced, "_#thetype#_", holder.InComments(), -1)
	replaced = strings.Replace(replaced, "_#theType.ID#_", holder.IDFieldName(), -1)
	replaced = strings.Replace(replaced, "_#theType.ID.Type#_", holder.IDFieldType(), -1)
	replaced = strings.Replace(replaced, "_#theType.Fields#_", holder.FieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#theType.Fields.Ref#_", holder.RefFieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.ID#_", holder.DbIDField(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.Fields#_", holder.DbFieldsEnum(), -1)

	return replaced, nil
}
