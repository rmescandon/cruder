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
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/rmescandon/myproject/datastore"
	yaml "gopkg.in/yaml.v1"
)

// ConfigFile path to the service config file
var ConfigFile string

type cfg struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Driver     string `yaml:"driver"`
	Datasource string `yaml:"datasource"`
}

var config cfg

// Launch starts the service
func Launch(configPath string) {
	err := readConfig(&config, configPath)
	if err != nil {
		log.Fatalf("Error parsing the config file: %v", err)
		return
	}

	err = datastore.OpenSysDatabase(config.Driver, config.Datasource)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	err = datastore.UpdateDatabase()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	handler := Router()
	port := strconv.Itoa(config.Port)
	address := strings.Join([]string{config.Host, ":", port}, "")

	log.Printf("Started service on port %s", port)
	log.Fatal(http.ListenAndServe(address, handler))
}

func readConfig(config *cfg, filePath string) error {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Error opening the config file.")
		return err
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Println("Error parsing the config file.")
		return err
	}

	return nil
}
