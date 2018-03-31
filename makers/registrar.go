// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018 Roberto Mier Escandon <rmescandon@gmail.com>
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

package makers

import (
	"github.com/rmescandon/cruder/errs"
	"github.com/rmescandon/cruder/log"
	"github.com/rmescandon/cruder/parser"
)

// Registrant adds a method to set the type holder to Maker interface.
// It is used to allow registering new makers
type Registrant interface {
	Maker
	SetTypeHolder(*parser.TypeHolder)
}

var registeredMakers map[string]Registrant

// Register registers a builtin maker
func Register(m Registrant) error {
	if m == nil {
		return errs.ErrNilObject
	}

	if registeredMakers[m.ID()] != nil {
		return errs.NewErrDuplicatedMaker(m.ID())
	}

	if registeredMakers == nil {
		registeredMakers = make(map[string]Registrant)
	}

	log.Infof("Registering plugin: %v", m.ID())
	registeredMakers[m.ID()] = m
	return nil
}
