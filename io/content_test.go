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

package io

import (
	"testing"

	check "gopkg.in/check.v1"
)

const (
	testContent = `
	package datastore

	import (
		"database/sql"
		"fmt"

		_ "github.com/mattn/go-sqlite3"
	)

	type Datastore interface {
		CreateMyTypeTable() error
		ListMyTypes() ([]MyType, error)
		GetMyType(id int) (MyType, error)
		FindMyType(name string) (MyType, error)
		CreateMyType(myType MyType) (int, error)
		UpdateMyType(id int, myType MyType)
		DeleteMyType(id int) error
	}

	type DB struct {
		*sql.DB
	}

	var Db *DB

	func OpenSysDatabase(driver, dataSource string) error {
		// Open the database connection
		db, err := sql.Open(driver, dataSource)
		if err != nil {
			return fmt.Errorf("Error opening the database: %v\n", err)
		}

		// Check that we have a valid database connection
		err = db.Ping()
		if err != nil {
			return fmt.Errorf("Error accessing the database: %v\n", err)
		}

		Db = &DB{db}

		return nil
	}
	`

	expectedContent = `package datastore

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Datastore interface {
	CreateMyTypeTable() error
	ListMyTypes() ([]MyType, error)
	GetMyType(id int) (MyType, error)
	FindMyType(name string) (MyType, error)
	CreateMyType(myType MyType) (int, error)
	UpdateMyType(id int, myType MyType)
	DeleteMyType(id int) error
}
type DB struct{ *sql.DB }

var Db *DB

func OpenSysDatabase(driver, dataSource string) error {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		return fmt.Errorf("Error opening the database: %v\n", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Error accessing the database: %v\n", err)
	}
	Db = &DB{db}
	return nil
}
`
)

type ContentSuite struct{}

var _ = check.Suite(&ContentSuite{})

func Test(t *testing.T) { check.TestingT(t) }

func (s *ContentSuite) TestValidSource(c *check.C) {
	content, err := NewContent(testContent)
	c.Assert(content, check.NotNil)
	c.Assert(err, check.IsNil)

	c.Assert(content.Ast, check.NotNil)

	str, err := content.String()
	c.Assert(err, check.IsNil)
	c.Assert(str, check.Equals, expectedContent)

	b, err := content.Bytes()
	c.Assert(err, check.IsNil)
	c.Assert(string(b), check.Equals, expectedContent)

	c.Assert(content.Trace(), check.IsNil)
}

func (s *ContentSuite) TestInvalidSource(c *check.C) {
	content, err := NewContent("whatever")
	c.Assert(content, check.IsNil)
	c.Assert(err, check.NotNil)
}
