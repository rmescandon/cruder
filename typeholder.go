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
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"
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
	Source  *goFile
	IDField typeField
	Fields  []typeField
	Decl    *ast.GenDecl
	Outputs []*Output
}

// Output represents an output generated file
type Output struct {
	Name     string
	File     *goFile
	Type     int
	Template string
}

// correlation between template identifier and output category
var outputCategories = map[string]int{
	"datastore": Datastore,
	"handler":   Handler,
	"router":    Router,
}

// returns type name in camel case, except first letter, which is lower case:
// "theType"
func (holder *TypeHolder) typeIdentifier() string {
	if len(holder.Name) > 0 {
		return strings.ToLower(string(holder.Name[0])) + holder.Name[1:len(holder.Name)]
	}
	return ""
}

// returns type name in lower case:
// "thetype"
func (holder *TypeHolder) typeInComments() string {
	return strings.ToLower(holder.Name)
}

// returns enum of the fields including type indentifier and field name:
// "theType.Field1, theType.Field2, theType.FieldN"
func (holder *TypeHolder) typeFieldsEnum() string {
	return holder.fieldsEnum(false)
}

// returns enum of type fields, including type identifiera and field name reference:
// "&theType.Field1, &theType.Field2, &theType.FieldN"
func (holder *TypeHolder) typeRefFieldsEnum() string {
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
			enum = ref + holder.typeIdentifierDotField(field.Name)
			continue
		}

		enum = enum + ", " + ref + holder.typeIdentifierDotField(field.Name)
	}
	return enum
}

func (holder *TypeHolder) typeIDFieldName() string {
	return holder.typeIdentifierDotField(holder.IDField.Name)
}

func (holder *TypeHolder) typeIDFieldType() string {
	return holder.typeIdentifierDotField(holder.IDField.Type)
}

// returns type identifier plus dot plus parameter fields name, like:
// "theType.Field1"
func (holder *TypeHolder) typeIdentifierDotField(fieldName string) string {
	return holder.typeIdentifier() + "." + fieldName
}

func (holder *TypeHolder) typeDbIDField() string {
	result := strings.ToLower(holder.IDField.Name) + " "
	if holder.IDField.Type == "int" {
		result = result + "serial "
	} else {
		result = result + holder.IDField.Type + " "
	}

	result = result + "primary key not null,"
	return result
}

// returns the type fields as they are used for SQL operations, like:
// "Field1 varchar(200),
// Field2 int,
// FieldN varchar(200)"
func (holder *TypeHolder) typeDbFieldsEnum() string {
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

func (holder *TypeHolder) appendOutputs() error {
	Log.Debugf("searching for available templates at %v", Config.TemplatesPath)
	availableTemplates, err := filepath.Glob(filepath.Join(Config.TemplatesPath, "*.template"))
	if err != nil {
		return fmt.Errorf("Error listing available templates: %v", err)
	}

	for _, templateFilePath := range availableTemplates {
		Log.Debugf("Found template: %v", filepath.Base(templateFilePath))

		err := holder.appendOutput(templateFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (holder *TypeHolder) appendOutput(templateFilePath string) error {
	templateID := templateIdentifier(templateFilePath)
	category := outputCategories[templateID]
	outputFilePath, err := holder.getOutputFilePathFor(category)
	if err != nil {
		return err
	}

	err = holder.generateOutputFile(templateFilePath, outputFilePath)
	if err != nil {
		return err
	}

	// read generated file content as byte array
	content, ast, err := fileToSyntaxTree(outputFilePath)
	if err != nil {
		return err
	}

	output := &Output{
		Name: templateID,
		File: &goFile{
			Path:    outputFilePath,
			Content: content,
			Ast:     ast,
		},
		Type:     category,
		Template: templateFilePath,
	}

	holder.Outputs = append(holder.Outputs, output)
	return nil
}

func (holder *TypeHolder) getOutputFilePathFor(category int) (string, error) {
	switch category {
	case Datastore:
		return filepath.Join(Config.Output, "datastore", holder.typeInComments()+".go"), nil
	case Handler:
		return filepath.Join(Config.Output, "service", holder.typeInComments()+".go"), nil
	case Router:
		return filepath.Join(Config.Output, "service", "router.go"), nil
	default:
		return "", errors.New("Invalid output category")
	}
}

func (holder *TypeHolder) generateOutputFile(templateFilePath, outputFilePath string) error {
	// don't write if file exists
	// FIXME this should not happen if output file is same as source one.
	// IN such case, original file types should be added to output
	_, err := os.Stat(outputFilePath)
	if err == nil {
		return fmt.Errorf("File %v already exists. Skip writting", outputFilePath)
	}

	// execute the replacement:
	// create needed dirs to outputPath
	ensureDir(filepath.Dir(outputFilePath))

	Log.Debugf("Loadig template: %v", filepath.Base(templateFilePath))
	templateContent, err := fileContentsAsString(templateFilePath)
	if err != nil {
		return fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := holder.replaceIn(templateContent)
	if err != nil {
		return fmt.Errorf("Error replacing type %v over template %v", holder.Name, filepath.Base(templateFilePath))
	}

	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("Could not create %v: %v", outputFilePath, err)
	}
	defer f.Close()

	_, err = f.WriteString(replacedStr)
	if err != nil {
		return fmt.Errorf("Error writing to output %v: %v", outputFilePath, err)
	}

	Log.Infof("Generated: %v", outputFilePath)
	return nil
}

func (holder *TypeHolder) replaceIn(templateContent string) (string, error) {
	replaced := templateContent

	replaced = strings.Replace(replaced, "_#TheType#_", holder.Name, -1)
	replaced = strings.Replace(replaced, "_#theType#_", holder.typeIdentifier(), -1)
	replaced = strings.Replace(replaced, "_#thetype#_", holder.typeInComments(), -1)
	replaced = strings.Replace(replaced, "_#theType.ID#_", holder.typeIDFieldName(), -1)
	replaced = strings.Replace(replaced, "_#theType.ID.Type#_", holder.typeIDFieldType(), -1)
	replaced = strings.Replace(replaced, "_#theType.Fields#_", holder.typeFieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#theType.Fields.Ref#_", holder.typeRefFieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.ID#_", holder.typeDbIDField(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.Fields#_", holder.typeDbFieldsEnum(), -1)

	return replaced, nil
}
