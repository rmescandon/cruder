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
	"log"
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
func OpenSysDatabase(driver, dataSource string) {
	// Open the database connection
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}

	// Check that we have a valid database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}

	Db = &DB{db}
}
