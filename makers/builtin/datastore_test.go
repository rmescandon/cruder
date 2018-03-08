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
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/io"
	"github.com/rmescandon/cruder/makers"
	"github.com/rmescandon/cruder/testdata"

	check "gopkg.in/check.v1"
)

const (
	testContent = `
	package datastore

	import (
		"database/sql"
		"fmt"
	)

	const listMyTypesSQL = "select id, name, description, subtypes from mytype order by id"
	const getMyTypeSQL = "select id, name, description, subtypes from mytype where id=$1"
	const findMyTypeSQL = "select id, name, description, subtypes from mytype where Name like '%$1%'"
	const createMyTypeSQL = "insert into mytype (name, description, subtypes) values ($1, $2, $3)"
	const updateMyTypeSQL = "update mytype set Name=$1, Description=$2, SubTypes=$3 where ID=$4"
	const deleteMyTypeSQL = "delete from mytype where id=$1"

	func (db *DB) CreateMyTypeTable() error {
		_, err := db.Exec(createMyTypeTableSQL)
		return err
	}
	func (db *DB) ListMyTypes() ([]MyType, error) {
		rows, err := db.Query(listMyTypesSQL)
		if err != nil {
			return []MyType{}, fmt.Errorf("Error retrieving database users: %v", err)
		}
		defer rows.Close()
		return db.rowsToMyTypes(rows)
	}
`

	testContentWithoutFunctions = `
	package datastore

	import (
		"database/sql"
		"fmt"
	)

	const listMyTypesSQL = "select id, name, description, subtypes from mytype order by id"
	const getMyTypeSQL = "select id, name, description, subtypes from mytype where id=$1"
	const findMyTypeSQL = "select id, name, description, subtypes from mytype where Name like '%$1%'"
	const createMyTypeSQL = "insert into mytype (name, description, subtypes) values ($1, $2, $3)"
	const updateMyTypeSQL = "update mytype set Name=$1, Description=$2, SubTypes=$3 where ID=$4"
	const deleteMyTypeSQL = "delete from mytype where id=$1"
`
)

type DatastoreSuite struct {
	datastore *Datastore
}

var _ = check.Suite(&DatastoreSuite{})

func Test(t *testing.T) { check.TestingT(t) }

func (d *DatastoreSuite) SetUpTest(c *check.C) {
	typeHolder, err := testdata.TestTypeHolder()
	c.Assert(err, check.IsNil)

	config.Config.Output, err = ioutil.TempDir("", "cruder_")
	c.Assert(err, check.IsNil)

	makers.BasePath = config.Config.Output

	d.datastore = &Datastore{makers.Base{TypeHolder: typeHolder}}
}

func (d *DatastoreSuite) TestID(c *check.C) {
	c.Assert(d.datastore.ID(), check.Equals, "datastore")
}

func (d *DatastoreSuite) TestOutputPath(c *check.C) {
	c.Assert(d.datastore.OutputFilepath(),
		check.Equals,
		filepath.Join(
			makers.BasePath,
			d.datastore.ID(),
			strings.ToLower(d.datastore.TypeHolder.Name)+".go"))
}

func (d *DatastoreSuite) TestOutputPathWhenNilType(c *check.C) {
	d.datastore.TypeHolder = nil
	c.Assert(d.datastore.OutputFilepath(), check.Equals, "")
}

func (d *DatastoreSuite) TestOutputPathWhenEmptyTypeName(c *check.C) {
	d.datastore.TypeHolder.Name = ""
	c.Assert(d.datastore.OutputFilepath(), check.Equals, "")
}

func (d *DatastoreSuite) TestOutputPathWhenEmptyBasePath(c *check.C) {
	makers.BasePath = ""
	c.Assert(d.datastore.OutputFilepath(),
		check.Equals,
		filepath.Join(
			d.datastore.ID(),
			strings.ToLower(d.datastore.TypeHolder.Name)+".go"))
}

func (d *DatastoreSuite) TestMake(c *check.C) {
	content, err := io.NewContent(testContent)
	c.Assert(err, check.IsNil)

	output, err := d.datastore.Make(content, content)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify the output contains the target type declaration
	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(str, "type MyType struct {"), check.Equals, true)
}

func (d *DatastoreSuite) TestMakeWhenNilParams(c *check.C) {
	output, err := d.datastore.Make(nil, nil)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (d *DatastoreSuite) TestMakeWhenNilGeneratedOutput(c *check.C) {
	currentOutput, err := io.NewContent(testContent)
	c.Assert(err, check.IsNil)

	output, err := d.datastore.Make(nil, currentOutput)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)
	c.Assert(err, check.Equals, errs.ErrNoContent)
}

func (d *DatastoreSuite) TestMakeWhenNilCurrentOutput(c *check.C) {
	generatedOutput, err := io.NewContent(testContent)
	c.Assert(err, check.IsNil)

	output, err := d.datastore.Make(generatedOutput, nil)
	c.Assert(err, check.IsNil)
	c.Assert(output, check.NotNil)

	// verify the output contains the target type declaration
	str, err := output.String()
	c.Assert(err, check.IsNil)
	c.Assert(strings.Contains(str, "type MyType struct {"), check.Equals, true)
}

func (d *DatastoreSuite) TestMakeWhenGeneratedOutputHasntFunctions(c *check.C) {
	generatedOutput, err := io.NewContent(testContentWithoutFunctions)
	c.Assert(err, check.IsNil)

	output, err := d.datastore.Make(generatedOutput, nil)
	c.Assert(err, check.NotNil)
	c.Assert(output, check.IsNil)

	switch err.(type) {
	case errs.ErrNotFound:
	default:
		c.Fail()
	}
}
