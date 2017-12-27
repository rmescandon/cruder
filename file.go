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

package cruder

import (
	"bytes"
	"os"
)

func fileContentsAsString(filepath string) (string, error) {
	b, err := fileBuffer(filepath)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func fileContentsAsByteArray(filepath string) ([]byte, error) {
	b, err := fileBuffer(filepath)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func fileBuffer(filepath string) (*bytes.Buffer, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
