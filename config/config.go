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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v1"

	flags "github.com/jessevdk/go-flags"
	gologging "github.com/op/go-logging"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/logging"
)

const (
	defaultSettingsFile = "settings.yaml"
)

// Options type holding possible cli params
type Options struct {
	Verbose       []bool `short:"v" long:"verbose" description:"Verbose output"`
	TypesFile     string `short:"t" long:"types" description:"file with struct types to consider for generating the skeletom code"`
	Output        string `short:"o" long:"output" description:"folder where building output structure of generated files"`
	Settings      string `short:"c" long:"config" description:"settings file path"`
	Version       string `yaml:"version"`
	TemplatesPath string `yaml:"templates"`
}

// Config holds received configuration from command line
var Config Options

// ValidateAndInitialize check received params and initialize default ones
func (c *Options) ValidateAndInitialize() error {

	if len(c.TypesFile) == 0 {
		return &flags.Error{
			Type:    flags.ErrHelp,
			Message: "Types file not provided",
		}
	}

	if len(c.Verbose) > 0 {
		logging.InitLogger(gologging.DEBUG)
	} else {
		logging.InitLogger(gologging.WARNING)
	}

	if len(c.Output) == 0 {
		// calculate current dir and set it as default output path
		dir, err := currentDir()
		if err != nil {
			return &flags.Error{
				Type:    flags.ErrUnknown,
				Message: "Internal server error when setting default output path",
			}
		}
		c.Output = dir
	}

	if len(c.Settings) == 0 {
		// calculate current dir and set it as default settings path
		dir, err := currentDir()
		if err != nil {
			return &flags.Error{
				Type:    flags.ErrUnknown,
				Message: "Internal server error when setting default settings path",
			}
		}
		c.Settings = filepath.Join(dir, defaultSettingsFile)
	}

	//normalize settings file path
	err := normalizePath(&c.Settings)
	if err != nil {
		return err
	}

	err = c.loadSettings()
	if err != nil {
		return fmt.Errorf("Error loading settings: %v", err)
	}

	err = c.normalizePaths()
	if err != nil {
		return err
	}

	return nil
}

func currentDir() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func (c *Options) loadSettings() error {
	b, err := io.FileContentAsByteArray(c.Settings)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		return fmt.Errorf("Error parsing the settngs file: %v", err)
	}

	return nil
}

func normalizePath(ptrStr *string) (err error) {
	if strings.Contains(*ptrStr, "~") {
		*ptrStr = strings.Replace(*ptrStr, "~", os.Getenv("HOME"), -1)
	}

	*ptrStr, err = filepath.Abs(*ptrStr)
	return
}

func (c *Options) normalizePaths() error {
	err := normalizePath(&c.TemplatesPath)
	if err != nil {
		return err
	}

	err = normalizePath(&c.TypesFile)
	if err != nil {
		return err
	}

	return nil
}
