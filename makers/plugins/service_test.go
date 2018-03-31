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

package main

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
	serviceTestContent = `
	 package service

	import (
		"io/ioutil"
		"log"
		"net/http"
		"strconv"
		"strings"

		"github.com/rmescandon/myproject/datastore"
		yaml "gopkg.in/yaml.v1"
	)

	// ConfigFile path to the service config file
	var ConfigFile string

	type cfg struct {
		Host       string
		Port       int   
		Driver     string
		Datasource string
	}

	var config cfg

	// Launch starts the service
	func Launch(configPath string) {
		err := readConfig(&config, configPath)
		if err != nil {
			log.Fatalf("Error parsing the config file: %v", err)
			return
		}

		err = datastore.OpenSysDatabase(config.Driver, config.Datasource)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		err = datastore.UpdateDatabase()
		if err != nil {
			log.Printf("%v", err)
			return
		}

		handler := Router()
		port := strconv.Itoa(config.Port)
		address := strings.Join([]string{config.Host, ":", port}, "")

		log.Printf("Started service on port %s", port)
		log.Fatal(http.ListenAndServe(address, handler))
	}

	func readConfig(config *cfg, filePath string) error {
		source, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println("Error opening the config file.")
			return err
		}

		err = yaml.Unmarshal(source, &config)
		if err != nil {
			log.Println("Error parsing the config file.")
			return err
		}

		return nil
	}
	 `
)

type ServiceSuite struct {
	service *Service
}

var _ = check.Suite(&ServiceSuite{})

func (s *ServiceSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.service = &Service{makers.Base{TypeHolder: typeHolder}}
}

func (s *ServiceSuite) TestID(c *check.C) {
	c.Assert(s.service.ID(), check.Equals, "service")
}

func (s *ServiceSuite) TestOutputPath(c *check.C) {
	c.Assert(s.service.OutputFilepath(),
		check.Equals,
		filepath.Join(makers.BasePath, s.service.ID(), s.service.ID()+".go"))
}

func (s *ServiceSuite) TestOutputPath_emptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.service.OutputFilepath(),
		check.Equals,
		filepath.Join(s.service.ID(), s.service.ID()+".go"))
}

func (s *ServiceSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(serviceTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	output, err := s.service.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(len(str) > 0, check.Equals, true)
	c.Assert(output, check.Equals, generatedOutput)
}

func (s *ServiceSuite) TestMake_existingOutput(c *check.C) {
	output, err := io.NewContent(serviceTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.service.Make(output, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}

func (s *ServiceSuite) TestMake_nilGeneratedOutput(c *check.C) {
	output, err := s.service.Make(nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}

func (s *ServiceSuite) TestMake_nilGeneratedOutputButExistsOutput(c *check.C) {
	output, err := io.NewContent(serviceTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.service.Make(nil, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}
