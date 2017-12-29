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

package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmescandon/cruder/example/myproject/datastore"
)

type myTypesResponse struct {
	MyTypes []datastore.MyType `json:"mytypes"`
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

// GetMyType handles reading mytype API operation
func GetMyType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//FIXME see how to make this conversion when generated the file. Probably with a switch
	myTypeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		replyWithError(
			http.StatusNotFound,
			errorResponse{
				Code:    "invalid-mytype-id",
				Message: "MyType was not found",
			},
			w,
		)
		return
	}

	myType, err := datastore.Db.GetMyType(myTypeID)
	if err != nil {
		log.Printf("Service error: %v", err)
		replyWithError(
			http.StatusInternalServerError,
			errorResponse{
				Code:    "get-mytype-failed",
				Message: "Could not get mytype info due to a server error",
			},
			w,
		)
		return
	}

	if err := json.NewEncoder(w).Encode(myType); err != nil {
		log.Printf("Service error: %v", err)
		replyWithError(
			http.StatusInternalServerError,
			errorResponse{
				Code:    "get-mytype-failed",
				Message: "A server error has happened when encoding the response",
			},
			w,
		)
		return
	}
}

// CreateMyType handles creating mytype API operation
func CreateMyType(w http.ResponseWriter, r *http.Request) {
	// Decode the body
	myType := datastore.MyType{}
	err := json.NewDecoder(r.Body).Decode(&myType)
	switch {
	// Check we have some data
	case err == io.EOF:
		replyWithError(
			http.StatusBadRequest,
			errorResponse{
				Code:    "empty-body-content",
				Message: "No mytype content supplied in body content",
			},
			w,
		)
		return
	// Check for parsing errors
	case err != nil:
		log.Printf("Request bad format: %v", err)
		replyWithError(
			http.StatusBadRequest,
			errorResponse{
				Code:    "invalid-body-content",
				Message: "Body content format is not valid",
			},
			w,
		)
		return
	}

	ID, err := datastore.Db.CreateMyType(myType)
	if err != nil {
		log.Printf("Service error creating mytpe: %v", err)
		replyWithError(
			http.StatusInternalServerError,
			errorResponse{
				Code:    "create-mytype-failed",
				Message: fmt.Sprintf("MyType creation failed due to a service error: %v", err),
			},
			w,
		)
		return
	}

	//FIXME see how to make this conversion when generated the file. Probably with a switch
	reply201Created(w, composeLocation(r, strconv.Itoa(ID)))
}

// UpdateMyType handles updating mytype API operation
func UpdateMyType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//FIXME see how to make this conversion when generated the file. Probably with a switch
	myTypeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		replyWithError(
			http.StatusNotFound,
			errorResponse{
				Code:    "invalid-mytype-id",
				Message: "MyType was not found",
			},
			w,
		)
		return
	}

	myType := datastore.MyType{}
	err = json.NewDecoder(r.Body).Decode(&myType)
	if err != nil {
		replyWithError(
			http.StatusBadRequest,
			errorResponse{
				Code:    "bad-body-content",
				Message: "Bad mytype supplied in body content",
			},
			w,
		)
		return
	}

	err = datastore.Db.UpdateMyType(myTypeID, myType)
	if err != nil {
		log.Printf("Service error: %v", err)
		replyWithError(
			http.StatusInternalServerError,
			errorResponse{
				Code:    "update-mytype-failed",
				Message: "Could not update requested mytype",
			},
			w,
		)
		return
	}
}

// DeleteMyType handles deleting mytype API operation
func DeleteMyType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//FIXME see how to make this conversion when generated the file. Probably with a switch
	myTypeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		replyWithError(
			http.StatusNotFound,
			errorResponse{
				Code:    "invalid-mytype-id",
				Message: "MyType was not found",
			},
			w,
		)
		return
	}

	err = datastore.Db.DeleteMyType(myTypeID)
	if err != nil {
		log.Printf("Service error: %v", err)
		replyWithError(
			http.StatusInternalServerError,
			errorResponse{
				Code:    "delete-mytype-failed",
				Message: "Could not delete requested mytype",
			},
			w,
		)
		return
	}

	reply204NoContent(w)
}
