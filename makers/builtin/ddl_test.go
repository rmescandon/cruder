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
	"fmt"
	"go/ast"
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
	ddlMyTypeContent = `package datastore

	// UpdateDatabase creates or updates tables by using DDL
	func UpdateDatabase() error {
		if err := Db.CreateMyTypeTable(); err != nil {
			return err
		}
	
		return nil
	}
	`

	ddlMyOtherTypeContent = `package datastore

	// UpdateDatabase creates or updates tables by using DDL
	func UpdateDatabase() error {
		if err := Db.CreateMyOtherTypeTable(); err != nil {
			return err
		}
	
		return nil
	}
	`

	ddlContentWithoutUpdateDatabaseFunc = `package datastore

	func otherMethod() error {
		if err := Db.CreateMyOtherTypeTable(); err != nil {
			return err
		}
	
		return nil
	}
	`

	ddlContentWithSeveralStatements = `package datastore

	// UpdateDatabase creates or updates tables by using DDL
	func UpdateDatabase() error {
		if err := Db.CreateMyTypeTable(); err != nil {
			return err
		}

		if err := Db.CreateMyOtherTypeTable(); err != nil {
			return err
		}

		if err := Db.CreateAnotherTypeTable(); err != nil {
			return err
		}
		return nil
	}
	`

	ddlContentWithoutStatements = `package datastore

	// UpdateDatabase creates or updates tables by using DDL
	func UpdateDatabase() error {
		return nil
	}
	`
)

type DDLSuite struct {
	ddl *DDL
}

var _ = check.Suite(&DDLSuite{})

func (s *DDLSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	s.ddl = &DDL{makers.Base{TypeHolder: typeHolder}}
}

func (s *DDLSuite) TestID(c *check.C) {
	c.Assert(s.ddl.ID(), check.Equals, "ddl")
}

func (s *DDLSuite) TestOutputPath(c *check.C) {
	c.Assert(s.ddl.OutputFilepath(),
		check.Equals,
		filepath.Join(
			makers.BasePath,
			fmt.Sprintf("datastore/%v.go", s.ddl.ID())))
}

func (s *DDLSuite) TestOutputPath_EmptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(s.ddl.OutputFilepath(),
		check.Equals,
		fmt.Sprintf("datastore/%v.go", s.ddl.ID()))
}

func (s *DDLSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(ddlMyOtherTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := s.ddl.Make(generatedOutput, currentOutput)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify that the output contains create table for both types
	str, err := output.String()
	c.Assert(err, check.IsNil)

	c.Assert(strings.Count(str, "if err := Db.CreateMyTypeTable(); err != nil {"), check.Equals, 1)
	c.Assert(strings.Count(str, "if err := Db.CreateMyOtherTypeTable(); err != nil {"), check.Equals, 1)
}

func (s *DDLSuite) TestMake_currentOutputWithoutStmts(c *check.C) {
	generatedOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(ddlContentWithoutStatements)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := s.ddl.Make(generatedOutput, currentOutput)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify that the output contains create table for both types
	str, err := output.String()
	c.Assert(err, check.IsNil)

	c.Assert(strings.Count(str, "if err := Db.CreateMyTypeTable(); err != nil {"), check.Equals, 1)
}

func (s *DDLSuite) TestMake_generatedOutputWithoutStmts(c *check.C) {
	generatedOutput, err := io.NewContent(ddlContentWithoutStatements)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := s.ddl.Make(generatedOutput, currentOutput)
	c.Assert(err, check.IsNil)
	// As not found target type in generated output, nothing is returned
	c.Assert(output, check.IsNil)
}

func (s *DDLSuite) TestMake_targetTypeExistsInOutput(c *check.C) {
	generatedOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := s.ddl.Make(generatedOutput, currentOutput)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.IsNil)
}

func (s *DDLSuite) TestMake_nilParams(c *check.C) {
	output, err := s.ddl.Make(nil, nil)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (s *DDLSuite) TestMake_nilGeneratedOutput(c *check.C) {
	currentOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)

	output, err := s.ddl.Make(nil, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (s *DDLSuite) TestMake_nilCurrentOutput(c *check.C) {
	generatedOutput, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)

	output, err := s.ddl.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)
	c.Assert(output, check.Equals, generatedOutput)
}

func (s *DDLSuite) TestMake_generatedOutputHasntUpdateDatabaseFunc(c *check.C) {
	generatedOutput, err := io.NewContent(ddlContentWithoutUpdateDatabaseFunc)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(ddlMyOtherTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := s.ddl.Make(generatedOutput, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)

	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}

func (s *DDLSuite) TestUpdateDatabaseFunction(c *check.C) {
	content, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)

	decl := findUpdateDatabaseFunction(content.Ast)
	c.Assert(decl, check.NotNil)
}

func (s *DDLSuite) TestUpdateDatabaseFunction_notFound(c *check.C) {
	content, err := io.NewContent(ddlContentWithoutUpdateDatabaseFunc)
	c.Assert(err, check.IsNil)

	decl := findUpdateDatabaseFunction(content.Ast)
	c.Assert(decl, check.IsNil)
}

func (s *DDLSuite) TestGetUpdateDatabaseStmts(c *check.C) {
	content, err := io.NewContent(ddlContentWithSeveralStatements)
	c.Assert(err, check.IsNil)

	stmts, err := getUpdateDatabaseStmts(content.Ast)
	c.Assert(err, check.IsNil)
	// 'return nil' is included
	c.Assert(stmts, check.HasLen, 4)
}

func (s *DDLSuite) TestGetUpdateDatabaseStmts_notFoundFunc(c *check.C) {
	content, err := io.NewContent(ddlContentWithoutUpdateDatabaseFunc)
	c.Assert(err, check.IsNil)

	stmts, err := getUpdateDatabaseStmts(content.Ast)
	c.Assert(err, check.NotNil)
	c.Assert(stmts, check.HasLen, 0)

	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}

func (s *DDLSuite) TestGetUpdateDatabaseTargetStatement(c *check.C) {
	content, err := io.NewContent(ddlContentWithSeveralStatements)
	c.Assert(err, check.IsNil)

	stmt, err := getUpdateDatabaseTargetStatement(content.Ast, "MyType")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.NotNil)

	stmt, err = getUpdateDatabaseTargetStatement(content.Ast, "MyOtherType")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.NotNil)

	stmt, err = getUpdateDatabaseTargetStatement(content.Ast, "AnotherType")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.NotNil)
}

func (s *DDLSuite) TestGetUpdateDatabaseTargetStatement_notFound(c *check.C) {
	content, err := io.NewContent(ddlContentWithSeveralStatements)
	c.Assert(err, check.IsNil)

	stmt, err := getUpdateDatabaseTargetStatement(content.Ast, "Foo")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.IsNil)
}

func (s *DDLSuite) TestSetStatements(c *check.C) {
	content1, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(content1, check.NotNil)

	content2, err := io.NewContent(ddlMyOtherTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(content2, check.NotNil)

	stmt, err := getUpdateDatabaseTargetStatement(content1.Ast, "MyType")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.NotNil)

	err = setStatements(content2.Ast, []ast.Stmt{stmt})
	c.Assert(err, check.IsNil)

	// verify that the output contains create table only for the first type
	str, err := content2.String()
	c.Assert(err, check.IsNil)

	c.Assert(strings.Count(str, "if err := Db.CreateMyTypeTable(); err != nil {"), check.Equals, 1)
	c.Assert(strings.Count(str, "if err := Db.CreateMyOtherTypeTable(); err != nil {"), check.Equals, 0)
}

func (s *DDLSuite) TestSetStatements_updateDatabaseFuncNotFound(c *check.C) {
	content1, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(content1, check.NotNil)

	content2, err := io.NewContent(ddlContentWithoutUpdateDatabaseFunc)
	c.Assert(err, check.IsNil)
	c.Assert(content2, check.NotNil)

	stmt, err := getUpdateDatabaseTargetStatement(content1.Ast, "MyType")
	c.Assert(err, check.IsNil)
	c.Assert(stmt, check.NotNil)

	err = setStatements(content2.Ast, []ast.Stmt{stmt})
	c.Assert(err, check.NotNil)
}

func (s *DDLSuite) TestSetStatements_setNone(c *check.C) {
	content, err := io.NewContent(ddlMyTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(content, check.NotNil)

	err = setStatements(content.Ast, []ast.Stmt{})
	c.Assert(err, check.IsNil)

	// verify that the output contains no statement
	r := findUpdateDatabaseFunction(content.Ast)
	c.Assert(r, check.NotNil)
	c.Assert(r.Body.List, check.HasLen, 0)
}
