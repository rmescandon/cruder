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

package log

import (
	"os"

	logging "github.com/op/go-logging"
)

var l = logging.MustGetLogger("cruder")

// InitLogger initializes logger for backend with the specified level
func InitLogger(level logging.Level) {
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")

	logging.SetBackend(backendLeveled)
}

// Errorf calls logger in eror level with format
func Errorf(format string, args ...interface{}) {
	l.Errorf(format, args)
}

// Error calls logger in error level
func Error(args ...interface{}) {
	l.Error(args)
}

// Warningf calls logger in warning level with format
func Warningf(format string, args ...interface{}) {
	l.Warningf(format, args)
}

// Warning calls logger in warning level
func Warning(args ...interface{}) {
	l.Warning(args)
}

// Infof calls logger in info level with format
func Infof(format string, args ...interface{}) {
	l.Infof(format, args)
}

// Info calls logger in info level
func Info(args ...interface{}) {
	l.Info(args)
}

// Debugf calls logger in debug level with format
func Debugf(format string, args ...interface{}) {
	l.Debugf(format, args)
}

// Debug calls logger in debug level
func Debug(args ...interface{}) {
	l.Debug(args)
}
