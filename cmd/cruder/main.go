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

package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	logging "github.com/op/go-logging"
	"github.com/rmescandon/cruder"
)

type commonOptions struct {
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
	File    string `short:"f" long:"file" description:"file with struct types to consider for generating the skeletom code"`
}

var options commonOptions
var parser = flags.NewParser(&options, flags.HelpFlag)

func addCommand(name string, shortHelp string, longHelp string, data interface{}) (*flags.Command, error) {
	cmd, err := parser.AddCommand(name, shortHelp, longHelp, data)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp || e.Type == flags.ErrCommandRequired {
				parser.WriteHelp(os.Stdout)
				return nil
			}
		}
		fmt.Println(err)
	}

	if len(options.File) == 0 {
		parser.WriteHelp(os.Stdout)
		return nil
	}

	if len(options.Verbose) > 0 {
		cruder.InitLogger(logging.DEBUG)
	} else {
		cruder.InitLogger(logging.WARNING)
	}

	err = generateSkeletonCode(options.File)

	return err
}

func generateSkeletonCode(sourceFile string) error {
	cruder.Log.Debug("Generating Skeleton Code...")

	//TODO IMPLEMENT

	return nil
}
