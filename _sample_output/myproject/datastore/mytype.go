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

package datastore

import (
	"database/sql"
	"fmt"
)

// MyType test type to generate skeletom code
type MyType struct {
	ID          int
	Name        string
	Description string
	SubTypes    bool
}

const createMyTypeTableSQL = `
	CREATE TABLE IF NOT EXISTS mytype (
		id           integer primary key not null,
		name         varchar(200),
		description  varchar(200),
		subtypes     boolean
	)
`

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
