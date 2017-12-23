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
	"errors"
	"path/filepath"
	"strings"
)

// TypeHolder holds a type previously read from file
type TypeHolder struct {
	Name        string
	IDFieldName string
	IDFieldType string
	Fields      []typeField
}

func newTypeHolder(typeName string, typeFields []typeField) *TypeHolder {
	return &TypeHolder{
		Name:   typeName,
		Fields: typeFields,
	}
}

func (holder *TypeHolder) typeIdentifier() string {
	if len(holder.Name) > 0 {
		return strings.ToLower(string(holder.Name[0])) + holder.Name[1:len(holder.Name)]
	}
	return ""
}

func (holder *TypeHolder) typeInComments() string {
	return strings.ToLower(holder.Name)
}

func (holder *TypeHolder) typeFieldsEnum() string {
	return holder.fieldsEnum(false)
}

func (holder *TypeHolder) typeRefFieldsEnum() string {
	return holder.fieldsEnum(true)
}

func (holder *TypeHolder) fieldsEnum(asRef bool) string {
	ref := ""
	if asRef {
		ref = "&"
	}

	var enum string
	for _, field := range holder.Fields {
		// skip ID field
		if field.Name == holder.IDFieldName {
			continue
		}

		if len(enum) == 0 {
			enum = ref + holder.typeIdentifierDotField(field.Name)
			continue
		}

		enum = enum + ", " + ref + holder.typeIdentifierDotField(field.Name)
	}
	return enum
}

func (holder *TypeHolder) typeIdentifierDotField(fieldName string) string {
	return holder.typeIdentifier() + "." + fieldName
}

func (holder *TypeHolder) typeDbIDField() string {
	result := strings.ToLower(holder.IDFieldName) + " "
	if holder.IDFieldType == "int" {
		result = result + "serial "
	} else {
		result = result + holder.IDFieldType + " "
	}

	result = result + "primary key not null,"
	return result
}

func (holder *TypeHolder) typeDbFieldsEnum() string {
	var result string
	for _, field := range holder.Fields {
		if field.Name == holder.IDFieldName {
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

func (holder *TypeHolder) getOutputFilePathFor(category int) (string, error) {
	switch category {
	case Datastore:
		return filepath.Join(Config.Output, "datastore", holder.Name, ".go"), nil
	case Handler:
		return filepath.Join(Config.Output, "service", holder.Name, ".go"), nil
	case Router:
		return filepath.Join(Config.Output, "service", "router.go"), nil
	default:
		return "", errors.New("Invalid output category")
	}
}
