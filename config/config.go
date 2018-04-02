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
	logging "github.com/op/go-logging"

	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/log"
)

const (
	defaultSettingsFile = "settings.yaml"
	defaultProjectURL   = "github.com/myuser/myproject"
	defaultAPIVersion   = "v1"
)

// Options type holding possible cli params
type Options struct {
	Args struct {
		TypesFile string `positional-arg-name:"types_file"`
	} `positional-args:"yes" required:"yes"`

	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	//TypesFile   string `short:"t" long:"types" description:"File with struct types to consider for generating the skeletom code" required:"yes"`
	Output      string `short:"o" long:"output" description:"Folder where building output structure of generated files"`
	ProjectURL  string `short:"u" long:"url" description:"Url of this project. If not specified 'github.com/myproject' is used"`
	APIVersion  string `short:"a" long:"apiversion" description:"Version of the REST api"`
	Settings    string `short:"c" long:"config" description:"Settings file path"`
	UserPlugins string `short:"p" long:"plugins" description:"Path to the folder with .so plugin files"`

	// Options loaded from settings file
	Version        string `yaml:"version"`
	TemplatesPath  string `yaml:"templates"`
	BuiltinPlugins string `yaml:"plugins"`
}

// Config holds received configuration from command line
var Config Options

// ValidateAndInitialize check received params and initialize default ones
func (c *Options) ValidateAndInitialize() error {
	if len(c.Args.TypesFile) == 0 {
		return &flags.Error{
			Type:    flags.ErrHelp,
			Message: "Types file not provided",
		}
	}

	if len(c.Verbose) > 0 {
		log.InitLogger(logging.DEBUG)
	} else {
		log.InitLogger(logging.WARNING)
	}

	err := c.setDefaultValuesWhenNeeded()
	if err != nil {
		return err
	}

	//normalize settings file path
	err = io.NormalizePath(&c.Settings)
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

// ReplaceInTemplate replaces config values in template
func (c *Options) ReplaceInTemplate(templateContent string) string {
	replaced := templateContent

	// github.com/myaccount/myproject
	replaced = strings.Replace(replaced, "_#PROJECT#_", c.ProjectURL, -1)

	// v1
	replaced = strings.Replace(replaced, "_#API.VERSION#_", c.APIVersion, -1)

	return replaced
}

func (c *Options) setDefaultValuesWhenNeeded() error {
	if len(c.Output) == 0 {
		// calculate current dir and set it as default output path
		dir, err := os.Getwd()
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
		dir, err := os.Getwd()
		if err != nil {
			return &flags.Error{
				Type:    flags.ErrUnknown,
				Message: "Internal server error when setting default settings path",
			}
		}
		c.Settings = filepath.Join(dir, defaultSettingsFile)
	}

	if len(c.ProjectURL) == 0 {
		c.ProjectURL = calculateProjectURL()
	}

	if len(c.APIVersion) == 0 {
		c.APIVersion = defaultAPIVersion
	}

	return nil
}

func (c *Options) loadSettings() error {
	b, err := io.FileToByteArray(c.Settings)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &Config)
	if err != nil {
		return fmt.Errorf("Error parsing the settngs file: %v", err)
	}

	return nil
}

func (c *Options) normalizePaths() error {
	err := io.NormalizePath(&c.TemplatesPath)
	if err != nil {
		return err
	}

	err = io.NormalizePath(&c.Args.TypesFile)
	if err != nil {
		return err
	}

	err = io.NormalizePath(&c.UserPlugins)
	if err != nil {
		return err
	}

	err = io.NormalizePath(&c.BuiltinPlugins)
	if err != nil {
		return err
	}

	return nil
}

func calculateProjectURL() string {
	gopath := os.Getenv("GOPATH")
	if len(gopath) == 0 {
		return defaultProjectURL
	}

	src := filepath.Join(gopath, "src")
	current, err := os.Getwd()
	if err != nil {
		log.Errorf("Could not get current directory. Use default project URL: %v", defaultProjectURL)
		return defaultProjectURL
	}

	if strings.Index(current, src) == -1 {
		log.Errorf("Current folder is not a $GOPATH/src subfolder. Use default project URL: %v", defaultProjectURL)
		return defaultProjectURL
	}

	return strings.TrimPrefix(current, src+"/")
}
