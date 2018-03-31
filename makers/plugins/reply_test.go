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
	replyTestContent = `
	package handler

	import (
		"encoding/json"
		"log"
		"net/http"
	)

	type emptyResponse struct{}

	type errorResponse struct {
		Code    string
		Message string
	}

	func replyWithError(statusCode int, errorBody errorResponse, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(errorBody); err != nil {
			log.Printf("Error forming the error response: %v\n", err)
		}
	}

	func reply200OK(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(emptyResponse{}); err != nil {
			log.Printf("Error forming the empty response: %v\n", err)
		}
	}

	func reply204NoContent(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNoContent)
	}

	func reply201Created(w http.ResponseWriter, location string) {
		w.Header().Set("Location", location)
		w.WriteHeader(http.StatusCreated)
	}

	func composeLocation(r *http.Request, id string) string {
		return "http://" + r.Host + r.URL.Path + "/" + id
	}
	`
)

type ReplySuite struct {
	r *Reply
}

var _ = check.Suite(&ReplySuite{})

func (s *ReplySuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.r = &Reply{makers.Base{TypeHolder: typeHolder}}
}

func (s *ReplySuite) TestID(c *check.C) {
	c.Assert(s.r.ID(), check.Equals, "reply")
}

func (s *ReplySuite) TestOutputPath(c *check.C) {
	c.Assert(s.r.OutputFilepath(),
		check.Equals,
		filepath.Join(makers.BasePath, "handler", s.r.ID()+".go"))
}

func (s *ReplySuite) TestOutputPath_emptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.r.OutputFilepath(),
		check.Equals,
		filepath.Join("handler", s.r.ID()+".go"))
}

func (s *ReplySuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(replyTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	output, err := s.r.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(len(str) > 0, check.Equals, true)
	c.Assert(output, check.Equals, generatedOutput)
}

func (s *ReplySuite) TestMake_existingOutput(c *check.C) {
	output, err := io.NewContent(replyTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.r.Make(output, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}

func (s *ReplySuite) TestMake_nilGeneratedOutput(c *check.C) {
	output, err := s.r.Make(nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}

func (s *ReplySuite) TestMake_nilGeneratedOutputButExistsOutput(c *check.C) {
	output, err := io.NewContent(replyTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	out, err := s.r.Make(nil, output)
	c.Assert(err, check.NotNil)
	c.Assert(out, check.IsNil)

	switch err.(type) {
	case errs.ErrOutputExists:
	default:
		c.Fail()
	}
}
