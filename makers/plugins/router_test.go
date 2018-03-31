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
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/testdata"
	check "gopkg.in/check.v1"
)

const (
	routerTestContent = `
	package service

	import (
		"net/http"

		"github.com/gorilla/mux"

		"github.com/rmescandon/myproject/handler"
	)

	const apiVersion = "v1"

	func composePath(operation string) string {
		return "/" + apiVersion + "/" + operation
	}

	// Router REST path multiplexer
	func Router() *mux.Router {
		router := mux.NewRouter().StrictSlash(true)

		router.Handle(composePath("mytype"), http.HandlerFunc(handler.CreateMyType)).Methods("POST")
		router.Handle(composePath("mytype"), http.HandlerFunc(handler.ListMyTypes)).Methods("GET")
		router.Handle(composePath("mytype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.GetMyType)).Methods("GET")
		router.Handle(composePath("mytype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.UpdateMyType)).Methods("PUT")
		router.Handle(composePath("mytype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.DeleteMyType)).Methods("DELETE")

		return router
	}
	`

	routerTestExistingContent = `
	package service

	import (
		"net/http"

		"github.com/gorilla/mux"

		"github.com/rmescandon/myproject/handler"
	)

	const apiVersion = "v1"

	func composePath(operation string) string {
		return "/" + apiVersion + "/" + operation
	}

	// Router REST path multiplexer
	func Router() *mux.Router {
		router := mux.NewRouter().StrictSlash(true)

		router.Handle(composePath("myothertype"), http.HandlerFunc(handler.CreateMyOtherType)).Methods("POST")
		router.Handle(composePath("myothertype"), http.HandlerFunc(handler.ListMyOtherTypes)).Methods("GET")
		router.Handle(composePath("myothertype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.GetMyOtherType)).Methods("GET")
		router.Handle(composePath("myothertype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.UpdateMyOtherType)).Methods("PUT")
		router.Handle(composePath("myothertype/{id:[a-zA-Z0-9-_:]+}"), http.HandlerFunc(handler.DeleteMyOtherType)).Methods("DELETE")

		return router
	}
	`
)

type RouterSuite struct {
	r *Router
}

var _ = check.Suite(&RouterSuite{})

func (s *RouterSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.r = &Router{makers.Base{TypeHolder: typeHolder}}
}

func (s *RouterSuite) TestID(c *check.C) {
	c.Assert(s.r.ID(), check.Equals, "router")
}

func (s *RouterSuite) TestOutputPath(c *check.C) {
	c.Assert(s.r.OutputFilepath(),
		check.Equals,
		filepath.Join(makers.BasePath, "service", s.r.ID()+".go"))
}

func (s *RouterSuite) TestOutputPath_emptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.r.OutputFilepath(),
		check.Equals,
		filepath.Join("service", s.r.ID()+".go"))
}

func (s *RouterSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(routerTestContent)
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

func (s *RouterSuite) TestMake_existingOutput(c *check.C) {
	generatedOutput, err := io.NewContent(routerTestContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	existingOutput, err := io.NewContent(routerTestExistingContent)
	c.Assert(err, check.IsNil)
	c.Assert(existingOutput, check.NotNil)

	out, err := s.r.Make(generatedOutput, existingOutput)
	c.Assert(err, check.IsNil)
	c.Assert(out, check.NotNil)

	stmts := getRouterFunctionStatements(out.Ast)
	c.Assert(stmts, check.HasLen, 10)
	m := findHandlersInStatements(stmts)

	for key := range m {
		switch key {
		case "CreateMyType":
		case "UpdateMyType":
		case "GetMyType":
		case "DeleteMyType":
		case "ListMyTypes":
		case "CreateMyOtherType":
		case "UpdateMyOtherType":
		case "GetMyOtherType":
		case "DeleteMyOtherType":
		case "ListMyOtherTypes":
		default:
			c.Fail()
		}
	}
}

func (s *RouterSuite) TestMake_nilGeneratedOutput(c *check.C) {
	output, err := s.r.Make(nil, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}
