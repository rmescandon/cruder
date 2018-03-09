// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
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

package builtin

import (
	"io/ioutil"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/testdata"
	check "gopkg.in/check.v1"
)

const (
	mainContent = `package main

	import (
		"fmt"
		"os"
	
		flags "github.com/jessevdk/go-flags"
		"github.com/rmescandon/myproject/service"
	)
	
	type opts struct {
		ConfigFile string
	}
	
	var options opts
	
	func main() {
		err := run()
		if err != nil {
			fmt.Printf("error parsing parameters: %v\r\n", err)
			return
		}
	
		service.Launch(options.ConfigFile)
	}
	
	func run() error {
		// Parse the command line arguments and execute the command
		parser := flags.NewParser(&options, flags.HelpFlag)
		_, err := parser.Parse()
	
		if err != nil {
			if e, ok := err.(*flags.Error); ok {
				if e.Type == flags.ErrHelp || e.Type == flags.ErrCommandRequired {
					parser.WriteHelp(os.Stdout)
					return nil
				}
			}
		}
	
		return err
	}
	 `
)

type MainSuite struct {
	m *Main
}

var _ = check.Suite(&MainSuite{})

func (s *MainSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.m = &Main{makers.Base{TypeHolder: typeHolder}}
}

func (s *MainSuite) TestID(c *check.C) {
	c.Assert(s.m.ID(), check.Equals, "main")
}

func (s *MainSuite) TestOutputPath(c *check.C) {
	c.Assert(s.m.OutputFilepath(),
		check.Equals,
		filepath.Join(makers.BasePath, "cmd/service", s.m.ID()+".go"))
}

func (s *MainSuite) TestOutputPath_emptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.m.OutputFilepath(),
		check.Equals,
		filepath.Join("cmd/service", s.m.ID()+".go"))
}

func (s *MainSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(mainContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	output, err := s.m.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify that the output contains create table for both types
	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(len(str) > 0, check.Equals, true)
	c.Assert(output, check.Equals, generatedOutput)
}

func (s *MainSuite) TestMake_existingOutput(c *check.C) {
	output, err := io.NewContent(mainContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.m.Make(output, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}

func (s *MainSuite) TestMake_nilGeneratedOutput(c *check.C) {
	output, err := s.m.Make(nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}

func (s *MainSuite) TestMake_nilGeneratedOutputButExistsOutput(c *check.C) {
	output, err := io.NewContent(mainContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.m.Make(nil, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}
