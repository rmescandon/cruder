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

func typeHoldersFromFile(typesFile string) ([]*TypeHolder, error) {
	var typeHolders []*TypeHolder
	typesMap, err := getTypesMaps(typesFile)
	if err != nil {
		return nil, err
	}

	for typeName := range typesMap {
		Log.Debugf("Found type: %v", typeName)
		typeHolder := newTypeHolder(typeName, typesMap[typeName])
		typeHolders = append(typeHolders, typeHolder)
	}

	return typeHolders, nil
}

func replace(templateFile string, typeHolders []*TypeHolder) error {
	var templateContent string

	for _, typeHolder := range typeHolders {
		// don't write if file exists
		// FIXME this should not happen if output file is same as source one.
		// IN such case, original file types should be added to output
		outputPath, err := getOutputFilePath(templateFile, typeHolder)
		_, err = os.Stat(outputPath)
		if err == nil {
			Log.Warningf("File %v already exists. Skip writting", outputPath)
			continue
		}

		// create needed dirs to outputPath
		ensureDir(filepath.Dir(outputPath))

		// read template content if first time
		if len(templateContent) == 0 {
			Log.Debugf("Loadig template: %v", filepath.Base(templateFile))
			templateContent, err = fileContentsAsString(templateFile)
			if err != nil {
				return fmt.Errorf("Error reading template file: %v", err)
			}
		}

		replacedStr, err := replaceOne(templateContent, typeHolder)
		if err != nil {
			return fmt.Errorf("Error replacing type %v over template %v", typeHolder.Name, templateFile)
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("Could not create %v: %v", outputPath, err)
		}
		defer f.Close()

		_, err = f.WriteString(replacedStr)
		if err != nil {
			return fmt.Errorf("Error writing to output %v: %v", outputPath, err)
		}

		Log.Infof("Generated: %v", outputPath)
	}

	return nil
}

func replaceOne(originalContent string, typeHolder *TypeHolder) (string, error) {
	replaced := originalContent

	replaced = strings.Replace(replaced, "_#TheType#_", typeHolder.Name, -1)
	replaced = strings.Replace(replaced, "_#theType#_", typeHolder.typeIdentifier(), -1)
	replaced = strings.Replace(replaced, "_#thetype#_", typeHolder.typeInComments(), -1)
	replaced = strings.Replace(replaced, "_#theType.ID#_", typeHolder.IDFieldName, -1)
	replaced = strings.Replace(replaced, "_#theType.ID.Type#_", typeHolder.IDFieldType, -1)
	replaced = strings.Replace(replaced, "_#theType.Fields#_", typeHolder.typeFieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#theType.Fields.Ref#_", typeHolder.typeRefFieldsEnum(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.ID#_", typeHolder.typeDbIDField(), -1)
	replaced = strings.Replace(replaced, "_#TheType.Db.Fields#_", typeHolder.typeDbFieldsEnum(), -1)

	return replaced, nil
}

func getOutputFilePath(templateFile string, typeHolder *TypeHolder) (string, error) {
	filename := filepath.Base(templateFile)
	switch filename {
	case "datastore.template":
		return typeHolder.getOutputFilePathFor(Datastore)
	case "handler.template":
		return typeHolder.getOutputFilePathFor(Handler)
	case "router.template":
		return typeHolder.getOutputFilePathFor(Router)
	default:
		return "", errors.New("Unknown template file")
	}
}

func ensureDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
