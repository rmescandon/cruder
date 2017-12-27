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
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	logging "github.com/op/go-logging"
)

// Options type holding possible cli params
type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	File    string `short:"f" long:"file" description:"file with struct types to consider for generating the skeletom code"`
	Output  string `short:"o" long:"output" description:"folder where building output structure of generated files"`
}

// Config holds received configuration from command line
var Config Options

// ValidateAndInitialize check received params and initialize default ones
func (c *Options) ValidateAndInitialize() error {

	if len(c.File) == 0 {
		return &flags.Error{
			Type:    flags.ErrHelp,
			Message: "Type file not provided",
		}
	}

	if len(c.Verbose) > 0 {
		initLogger(logging.DEBUG)
	} else {
		initLogger(logging.WARNING)
	}

	if len(c.Output) == 0 {
		// calculate current dir and set it as default output path
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return &flags.Error{
				Type:    flags.ErrUnknown,
				Message: "Internal server error when setting default output path",
			}
		}
		c.Output = dir
	}

	return nil
}
