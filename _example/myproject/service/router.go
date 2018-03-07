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
