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

package engine

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/makers/builtin"
	"github.com/rmescandon/cruder/parser"
)

// Run generates the code, based on loaded configuration and available templates
func Run() error {
	log.Info("Generating code...")

	//TODO TEST
	builtin.DoNothing()

	makers.BasePath = config.Config.Output

	err := loadPlugins()
	if err != nil {
		return err
	}

	source, err := io.NewGoFile(config.Config.TypesFile)
	if err != nil {
		return fmt.Errorf("Error reading go source file: %v", err)
	}

	typeHolders, err := parser.ComposeTypeHolders(source)
	if err != nil {
		return fmt.Errorf("Error composing type holders from types file: %v", err)
	}

	templates, err := availableTemplates()
	if err != nil {
		return fmt.Errorf("Error listing available templates: %v", err)
	}

	processMakers(typeHolders, templates)

	return nil
}

func loadPlugins() error {
	plugins, err := filepath.Glob(filepath.Join(config.Config.BuiltinPlugins, "*.so"))
	if err != nil {
		return err
	}

	userPlugins, err := filepath.Glob(filepath.Join(config.Config.UserPlugins, "*.so"))
	if err != nil {
		return err
	}

	plugins = append(plugins, userPlugins...)

	for _, p := range plugins {
		// Once the plugin is open, its init() func is called and
		// the plugins register themselves as makers
		_, err = plugin.Open(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func availableTemplates() ([]string, error) {
	log.Infof("searching for available templates at %v", config.Config.TemplatesPath)
	return filepath.Glob(filepath.Join(config.Config.TemplatesPath, "*.template"))
}

func processMakers(holders []*parser.TypeHolder, templates []string) {
	for _, t := range templates {
		log.Infof("Found template: %v", filepath.Base(t))
		for _, h := range holders {
			err := processMaker(h, t)
			if err != nil {
				// warn of the error but continue with next maker
				log.Warning(err)
			}
		}

	}
}

func processMaker(typeHolder *parser.TypeHolder, template string) error {
	maker, err := makers.Get(template)
	switch err.(type) {
	case errs.ErrNotFound:
		log.Warningf("Maker not found, skipped - %v", err)
		return nil
	}
	if err != nil {
		return err
	}

	maker.(makers.Registrant).SetTypeHolder(typeHolder)

	merged, err := merge(typeHolder, template)
	if err != nil {
		return err
	}

	generatedOutput, err := io.NewContent(merged)
	if err != nil {
		return err
	}

	var currentOutput *io.Content
	currentOutputFile, err := io.NewGoFile(maker.OutputFilepath())
	if err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			currentOutput = nil
		default:
			return err
		}
	} else {
		currentOutput = &currentOutputFile.Content
	}

	result, err := maker.Make(generatedOutput, currentOutput)
	if err != nil {
		return err
	}

	if result != nil {
		err = io.EnsureDir(filepath.Dir(maker.OutputFilepath()))
		if err != nil {
			return err
		}

		err = io.ASTToFile(result.Ast, maker.OutputFilepath())
		if err != nil {
			return err
		}

		log.Infof("Generated: %v", maker.OutputFilepath())
	}

	return nil
}

// Merges type, config and template, returning the result as a string
func merge(typeHolder *parser.TypeHolder, templateFilepath string) (string, error) {
	// execute the replacement
	log.Debugf("Loading template: %v", filepath.Base(templateFilepath))
	templateContent, err := io.FileToString(templateFilepath)
	if err != nil {
		return "", fmt.Errorf("Error reading template file: %v", err)
	}

	replacedStr, err := typeHolder.ReplaceInTemplate(templateContent)
	if err != nil {
		return "", fmt.Errorf("Error replacing type %v over template %v",
			typeHolder.Name, filepath.Base(templateFilepath))
	}

	replacedStr, err = config.Config.ReplaceInTemplate(replacedStr)
	if err != nil {
		return "", fmt.Errorf("Error replacing configuration over template %v",
			filepath.Base(templateFilepath))
	}

	return replacedStr, err
}
