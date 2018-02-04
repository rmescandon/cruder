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
	"os"
	"path/filepath"

	"github.com/rmescandon/cruder/config"
)

// Reply struct holding data to copy reply template
type Reply struct {
	CopyMaker
}

// ID returns 'reply' as this maker identifier
func (r *Reply) ID() string {
	return "reply"
}

// OutputFilepath returns the path to the output file
func (r *Reply) OutputFilepath() string {
	return filepath.Join(config.Config.Output, "service/reply.go")
}

// Make copies template to output path
func (r *Reply) Make() error {
	_, err := os.Stat(r.OutputFilepath())
	if err == nil {
		return NewErrOutputExists(r.OutputFilepath())
	}

	return r.CopyMaker.copy(r.OutputFilepath())
}

func init() {
	register(&Reply{})
}
