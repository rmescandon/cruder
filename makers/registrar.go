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
	"fmt"

	"github.com/rmescandon/cruder/parser"
)

// Registrant adds a method to set the type holder to Maker interface.
// It is used to allow registering new makers
type Registrant interface {
	Maker
	ID() string
	SetTypeHolder(*parser.TypeHolder)
	SetTemplate(string)
}

var registeredMakers map[string]Registrant

// Register registers a builtin maker
func register(m Registrant) error {
	if registeredMakers[m.ID()] != nil {
		return fmt.Errorf("cannot register duplicated maker %q", m.ID())
	}
	if registeredMakers == nil {
		registeredMakers = make(map[string]Registrant)
	}
	registeredMakers[m.ID()] = m
	return nil
}
