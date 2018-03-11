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

package parser

import (
	"io/ioutil"

	"github.com/rmescandon/cruder/io"

	check "gopkg.in/check.v1"
)

const (
	astTestContent = `
package datastore

import (
	"database/sql"
	"fmt"
)

type MyInterface interface {
	MyFunc() string
}

type MyOtherInterface interface {
	MyOtherFunc() string
}

// MyType test type to generate skeletom code
type MyType struct {
	ID          int
	Name        string
	Description string
	SubTypes    bool
}

const listMyTypesSQL = "select id, name, description, subtypes from mytype order by id"
const getMyTypeSQL = "select id, name, description, subtypes from mytype where id=$1"
const findMyTypeSQL = "select id, name, description, subtypes, from mytype where name like '%$1%'"
const createMyTypeSQL = "insert into mytype (name, description, subtypes) values ($1,$2,$3)"
const updateMyTypeSQL = "update mytype set name=$1, description=$2, subtypes=$3 where id=$4"
const deleteMyTypeSQL = "delete from mytype where id=$1"

// CreateMyTypeTable creates the database table
func (db *DB) CreateMyTypeTable() error {
	_, err := db.Exec(createMyTypeTableSQL)
	return err
}

// ListMyTypes returns all the registers of the table
func (db *DB) ListMyTypes() ([]MyType, error) {
	rows, err := db.Query(listMyTypesSQL)
	if err != nil {
		return []MyType{}, fmt.Errorf("Error retrieving database users: %v", err)
	}
	defer rows.Close()

	return db.rowsToMyTypes(rows)
}

// GetMyType returns a specific register
func (db *DB) GetMyType(id int) (MyType, error) {
	row := db.QueryRow(getMyTypeSQL, id)
	myType, err := db.rowToMyType(row)
	if err != nil {
		return MyType{}, fmt.Errorf("Error retrieving mytype register: %v", err)
	}
	return myType, err
}

// FindMyType searches for a specific register
func (db *DB) FindMyType(name string) (MyType, error) {
	row := db.QueryRow(findMyTypeSQL, name)
	myType, err := db.rowToMyType(row)
	if err != nil {
		return MyType{}, fmt.Errorf("Error searching mytype registers: %v", err)
	}
	return myType, err
}

// CreateMyType Inserts a new register
func (db *DB) CreateMyType(myType MyType) (int, error) {
	result, err := db.Exec(createMyTypeSQL, myType.Name, myType.Description, myType.SubTypes)
	if err != nil {
		return -1, fmt.Errorf("Error creating mytype register: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

// UpdateMyType updates a register
func (db *DB) UpdateMyType(id int, myType MyType) error {
	_, err := db.Exec(updateMyTypeSQL, myType.Name, myType.Description, myType.SubTypes, id)
	if err != nil {
		return fmt.Errorf("Error updating mytype register: %v", err)
	}
	return nil
}

// DeleteMyType deletes a register
func (db *DB) DeleteMyType(id int) error {
	_, err := db.Exec(deleteMyTypeSQL, id)
	if err != nil {
		return fmt.Errorf("Error deleting mytype register: %v", err)
	}
	return nil
}

func (db *DB) rowToMyType(row *sql.Row) (MyType, error) {
	myType := MyType{}
	err := row.Scan(&myType.ID, &myType.Name, &myType.Description, &myType.SubTypes)
	if err != nil {
		return MyType{}, err
	}

	return myType, nil
}

func (db *DB) nextRowToMyType(rows *sql.Rows) (MyType, error) {
	myType := MyType{}
	err := rows.Scan(&myType.ID, &myType.Name, &myType.Description, &myType.SubTypes)
	if err != nil {
		return MyType{}, err
	}

	return myType, nil
}

func (db *DB) rowsToMyTypes(rows *sql.Rows) ([]MyType, error) {
	myTypeList := []MyType{}

	for rows.Next() {
		myType, err := db.nextRowToMyType(rows)
		if err != nil {
			return nil, err
		}
		myTypeList = append(myTypeList, myType)
	}

	return myTypeList, nil
}
`
)

// gopkg.in/check.v1 stuff
type AstSuite struct{}

var _ = check.Suite(&AstSuite{})

func (s *AstSuite) SetUpTest(c *check.C) {}

func (s *AstSuite) TestGetTypeDecls(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)
	decls := getTypeDecls(content.Ast)
	c.Assert(decls, check.HasLen, 3)
}

func (s *AstSuite) TestGetStructs(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)
	decls := getStructs(content.Ast)
	c.Assert(decls, check.HasLen, 1)
}

func (s *AstSuite) TestGetFuncDecls(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)
	decls := GetFuncDecls(content.Ast)
	c.Assert(decls, check.HasLen, 10)
}

func (s *AstSuite) TestGetInterfaces(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)
	ifaces := getInterfaces(content.Ast)
	c.Assert(ifaces, check.HasLen, 2)
}

func (s *AstSuite) TestInterface(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)

	iface := GetInterface(content.Ast, "MyInterface")
	c.Assert(iface, check.NotNil)

	methods := GetInterfaceMethods(iface)
	c.Assert(methods, check.HasLen, 1)

	c.Assert(HasMethod(iface, "MyFunc"), check.Equals, true)
	c.Assert(HasMethod(iface, "NotExistingFunc"), check.Equals, false)
}

func (s *AstSuite) TestUpdateInterfaceMethods(c *check.C) {
	content, err := io.NewContent(astTestContent)
	c.Assert(err, check.IsNil)

	iface := GetInterface(content.Ast, "MyInterface")
	c.Assert(iface, check.NotNil)

	iface2 := GetInterface(content.Ast, "MyOtherInterface")
	c.Assert(iface2, check.NotNil)

	fields2 := GetInterfaceMethods(iface2)
	c.Assert(fields2, check.HasLen, 1)
	AddMethod(iface, fields2[0])

	fields := GetInterfaceMethods(iface)
	c.Assert(fields, check.HasLen, 2)
}

func (s *AstSuite) TestComposeTypeHolder(c *check.C) {
	f, err := ioutil.TempFile("", "")
	c.Assert(err, check.IsNil)
	defer f.Close()

	testTypeFileContent := `
	package mytype
	
	// MyType test type to generate skeletom code
	type MyType struct {
		ID            int
		Name          string
		Description   string
		TheBoolThing  bool
		TheFloatThing float
	}	
	`

	_, err = f.WriteString(testTypeFileContent)
	c.Assert(err, check.IsNil)
	gof, err := io.NewGoFile(f.Name())
	c.Assert(err, check.IsNil)
	th, err := ComposeTypeHolders(gof)
	c.Assert(err, check.IsNil)
	c.Assert(th, check.HasLen, 1)
	c.Assert(th[0].Name, check.Equals, "MyType")
}
