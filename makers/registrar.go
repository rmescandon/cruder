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

package makers

import (
	"fmt"

	"github.com/rmescandon/cruder/parser"
)

// BuiltinMaker adds a method to set the type holder to Maker interface.
// It is used internally to register builtin makers
type BuiltinMaker interface {
	Maker
	ID() string
	SetTypeHolder(*parser.TypeHolder)
	SetTemplate(string)
}

var builtinMakers map[string]BuiltinMaker

// Register registers a builtin maker
func Register(m BuiltinMaker) error {
	if builtinMakers[m.ID()] != nil {
		return fmt.Errorf("cannot register duplicated maker %q", m.ID())
	}
	if builtinMakers == nil {
		builtinMakers = make(map[string]BuiltinMaker)
	}
	builtinMakers[m.ID()] = m
	return nil
}
