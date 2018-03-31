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
	"strings"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/testdata"
	check "gopkg.in/check.v1"
)

const (
	handlerTestContent = `package handler

	import (
		"encoding/json"
		"fmt"
		"io"
		"log"
		"net/http"
		"strconv"

		"github.com/gorilla/mux"
		"github.com/rmescandon/myproject/datastore"
	)

	type myTypesResponse struct {
		MyTypes []datastore.MyType
	}

	// ListMyTypes handles listing mytypes API operation
	func ListMyTypes(w http.ResponseWriter, r *http.Request) {
		myTypes, err := datastore.Db.ListMyTypes()
		if err != nil {
			log.Printf("Service error: %v", err)
			replyWithError(
				http.StatusInternalServerError,
				errorResponse{
					Code:    "list-mytypes-failed",
					Message: "Could not list available mytypes due to a server error",
				},
				w,
			)
			return
		}

		response := myTypesResponse{MyTypes: myTypes}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Service error: %v", err)
			replyWithError(
				http.StatusInternalServerError,
				errorResponse{
					Code:    "list-applications-failed",
					Message: "A server error has happened when encoding the response",
				},
				w,
			)
			return
		}
	}
	`
)

type HandlerSuite struct {
	handler *Handler
}

var _ = check.Suite(&HandlerSuite{})

func (s *HandlerSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.handler = &Handler{makers.Base{TypeHolder: typeHolder}}
}

func (s *HandlerSuite) TestID(c *check.C) {
	c.Assert(s.handler.ID(), check.Equals, "handler")
}

func (s *HandlerSuite) TestOutputPath(c *check.C) {
	c.Assert(s.handler.OutputFilepath(),
		check.Equals,
		filepath.Join(
			makers.BasePath,
			s.handler.ID(),
			strings.ToLower(s.handler.TypeHolder.Name)+".go"))
}

func (s *HandlerSuite) TestOutputPath_nilType(c *check.C) {
	s.handler.TypeHolder = nil
	c.Assert(s.handler.OutputFilepath(), check.Equals, "")
}

func (s *HandlerSuite) TestOutputPath_emptyTypeName(c *check.C) {
	s.handler.TypeHolder.Name = ""
	c.Assert(s.handler.OutputFilepath(), check.Equals, "")
}

func (s *HandlerSuite) TestOutputPath_emptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.handler.OutputFilepath(),
		check.Equals,
		filepath.Join(
			s.handler.ID(),
			strings.ToLower(s.handler.TypeHolder.Name)+".go"))
}

func (s *HandlerSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(handlerTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	output, err := s.handler.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify that the output contains create table for both types
	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(len(str) > 0, check.Equals, true)
	c.Assert(output, check.Equals, generatedOutput)
}

func (s *HandlerSuite) TestMake_existingOutput(c *check.C) {
	output, err := io.NewContent(handlerTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.handler.Make(output, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}

func (s *HandlerSuite) TestMake_nilGeneratedOutput(c *check.C) {
	output, err := s.handler.Make(nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}

func (s *HandlerSuite) TestMake_nilGeneratedOutputButExistsOutput(c *check.C) {
	output, err := io.NewContent(handlerTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.handler.Make(nil, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}
