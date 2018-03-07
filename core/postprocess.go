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

package core

import (
	"os/exec"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
	"github.com/rmescandon/cruder/log"
)

// Defined cli commands
const (
	GoLintCmd = iota
	GoFmtCmd
	GoVetCmd
)

func checkAndFormatResults() error {
	err := gofmt(config.Config.Output)
	if err != nil {
		return err
	}

	err = golint(filepath.Join(config.Config.Output, "..."))
	if err != nil {
		return err
	}

	err = govet(filepath.Join(config.Config.Output, "..."))
	if err != nil {
		return err
	}

	return nil
}

func golint(target string) error {
	return execCmd(GoLintCmd, target)
}

func gofmt(target string) error {
	return execCmd(GoFmtCmd, target)
}

func govet(target string) error {
	return execCmd(GoVetCmd, target)
}

func execCmd(cmd int, target string) error {
	var bs []byte
	var err error

	switch cmd {
	case GoLintCmd:
		bs, err = exec.Command("golint", target).CombinedOutput()
		if err != nil {
			return err
		}
	case GoFmtCmd:
		bs, err = exec.Command("gofmt", "-l", "-w", target).CombinedOutput()
		if err != nil {
			return err
		}
	case GoVetCmd:
		bs, err = exec.Command("go vet", target).CombinedOutput()
		if err != nil {
			return err
		}
	}

	log.Infof("%v", string(bs))
	return nil
}
