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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	yaml "gopkg.in/yaml.v1"
	"launchpad.net/indore/ams/datastore"
)

type opts struct {
	ConfigFile string `short:"c" long:"config" description:"Path to the config file" default:"./settings.yaml"`
}

type cfg struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Driver     string `yaml:"driver"`
	Datasource string `yaml:"datasource"`
}

var options opts
var config cfg

func main() {
	err := run()
	if err != nil {
		fmt.Printf("error parsing parameters: %v\r\n", err)
		return
	}

	err := config.ReadConfig(&config, options.ConfigFile)
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

func run() error {
	// Parse the command line arguments and execute the command
	parser := flags.NewParser(&options, flags.HelpFlag)
	_, err := parser.Parse()

	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp || e.Type == flags.ErrCommandRequired {
				parser.WriteHelp(os.Stdout)
				return nil
			}
		}
	}

	return err
}

// ReadConfig parses the config file
func ReadConfig(config *cfg, filePath string) error {
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
