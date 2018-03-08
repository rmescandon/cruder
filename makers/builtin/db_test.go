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

package builtin

import (
	"fmt"
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
	oneTypeContent = `
	package datastore

	import (
		"database/sql"
		"fmt"

		// Import the sqlite3 database driver
		_ "github.com/mattn/go-sqlite3"
	)

	// Datastore interface for different data storages
	type Datastore interface {
		CreateMyTypeTable() error
		ListMyTypes() ([]MyType, error)
		GetMyType(id int) (MyType, error)
		FindMyType(name string) (MyType, error)
		CreateMyType(myType MyType) (int, error)
		UpdateMyType(id int, myType MyType)
		DeleteMyType(id int) error
	}

	// DB struct holding database implementation for datastore
	type DB struct {
		*sql.DB
	}

	// Db pointer to database hander
	var Db *DB

	// OpenSysDatabase Return an open database connection
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

	otherTypeContent = `
	package datastore

	import (
		"database/sql"
		"fmt"

		// Import the sqlite3 database driver
		_ "github.com/mattn/go-sqlite3"
	)

	// Datastore interface for different data storages
	type Datastore interface {
		CreateOtherTypeTable() error
		ListOtherTypes() ([]OtherType, error)
		GetOtherType(id int) (OtherType, error)
		FindOtherType(name string) (OtherType, error)
		CreateOtherType(myType OtherType) (int, error)
		UpdateOtherType(id int, myType OtherType)
		DeleteOtherType(id int) error
	}

	// DB struct holding database implementation for datastore
	type DB struct {
		*sql.DB
	}

	// Db pointer to database hander
	var Db *DB

	// OpenSysDatabase Return an open database connection
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

	contentWithoutDatastoreIface = `
	package datastore

	import (
		"database/sql"
		"fmt"

		// Import the sqlite3 database driver
		_ "github.com/mattn/go-sqlite3"
	)

	// DB struct holding database implementation for datastore
	type DB struct {
		*sql.DB
	}

	// Db pointer to database hander
	var Db *DB

	// OpenSysDatabase Return an open database connection
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
)

type DbSuite struct {
	db *Db
}

var _ = check.Suite(&DbSuite{})

func (d *DbSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	d.db = &Db{makers.Base{TypeHolder: typeHolder}}
}

func (d *DbSuite) TestID(c *check.C) {
	c.Assert(d.db.ID(), check.Equals, "db")
}

func (d *DbSuite) TestOutputPath(c *check.C) {
	c.Assert(d.db.OutputFilepath(),
		check.Equals,
		filepath.Join(
			makers.BasePath,
			fmt.Sprintf("datastore/%v.go", d.db.ID())))
}

func (d *DbSuite) TestOutputPathWhenEmptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(d.db.OutputFilepath(),
		check.Equals,
		fmt.Sprintf("datastore/%v.go", d.db.ID()))
}

func (d *DbSuite) TestMake(c *check.C) {
	generatedOutput, err := io.NewContent(oneTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(otherTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := d.db.Make(generatedOutput, currentOutput)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify the output contains the target type declaration
	str, err := output.String()
	c.Assert(err, check.IsNil)

	c.Assert(strings.Contains(str, "CreateMyTypeTable() error"), check.Equals, true)
	c.Assert(strings.Contains(str, "ListMyTypes() ([]MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "GetMyType(id int) (MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "FindMyType(name string) (MyType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "CreateMyType(myType MyType) (int, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "UpdateMyType(id int, myType MyType)"), check.Equals, true)
	c.Assert(strings.Contains(str, "DeleteMyType(id int) error"), check.Equals, true)

	c.Assert(strings.Contains(str, "CreateOtherTypeTable() error"), check.Equals, true)
	c.Assert(strings.Contains(str, "ListOtherTypes() ([]OtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "GetOtherType(id int) (OtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "FindOtherType(name string) (OtherType, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "CreateOtherType(myType OtherType) (int, error)"), check.Equals, true)
	c.Assert(strings.Contains(str, "UpdateOtherType(id int, myType OtherType)"), check.Equals, true)
	c.Assert(strings.Contains(str, "DeleteOtherType(id int) error"), check.Equals, true)
}

func (d *DbSuite) TestMakeWhenNilParams(c *check.C) {
	output, err := d.db.Make(nil, nil)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (d *DbSuite) TestMakeWhenNilGeneratedOutput(c *check.C) {
	currentOutput, err := io.NewContent(oneTypeContent)
	c.Assert(err, check.IsNil)

	output, err := d.db.Make(nil, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (d *DbSuite) TestMakeWhenNilCurrentOutput(c *check.C) {
	generatedOutput, err := io.NewContent(oneTypeContent)
	c.Assert(err, check.IsNil)

	output, err := d.db.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)
	c.Assert(output, check.Equals, generatedOutput)
}

func (d *DbSuite) TestMakeWhenGeneratedOutputHasntDatastoreInterface(c *check.C) {
	generatedOutput, err := io.NewContent(contentWithoutDatastoreIface)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(otherTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := d.db.Make(generatedOutput, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)

	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}

func (d *DbSuite) TestMakeWhenCurrentOutputHasntDatastoreInterface(c *check.C) {
	generatedOutput, err := io.NewContent(oneTypeContent)
	c.Assert(err, check.IsNil)
	c.Assert(generatedOutput, check.NotNil)

	currentOutput, err := io.NewContent(contentWithoutDatastoreIface)
	c.Assert(err, check.IsNil)
	c.Assert(currentOutput, check.NotNil)

	output, err := d.db.Make(generatedOutput, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)

	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}
