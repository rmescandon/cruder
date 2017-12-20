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
	"os"

	"github.com/jessevdk/go-flags"
)

type commonOptions struct {
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`
}

var parser = flags.NewParser(&commonOptions{}, flags.Default)

func addCommand(name string, shortHelp string, longHelp string, data interface{}) (*flags.Command, error) {
	cmd, err := parser.AddCommand(name, shortHelp, longHelp, data)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
}
