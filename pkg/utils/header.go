/*
*	Copyright (C) 2026 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package utils

import (
	"os"

	"github.com/spf13/pflag"
)

const DefaultHeaderFile string = "hack/boilerplate.go.txt"

// HeaderFlag registers the --header/-f flag on set and returns a pointer to
// the contents of the header file it names. A missing default header file is
// treated as an empty header, but a file passed explicitly must exist.
func HeaderFlag(set *pflag.FlagSet) *string {
	h := &headerValue{path: DefaultHeaderFile, contents: new(string)}
	if data, e := os.ReadFile(DefaultHeaderFile); e == nil {
		*h.contents = string(data)
	}

	set.VarP(h, "header", "f", "Specify an optional header to apply to generated files.")
	return h.contents
}

// headerValue is a pflag.Value that reads the file it is set to.
type headerValue struct {
	path     string
	contents *string
}

func (h *headerValue) String() string { return h.path }

func (h *headerValue) Type() string { return "string" }

func (h *headerValue) Set(path string) error {
	data, e := os.ReadFile(path)
	if e != nil {
		return e
	}

	h.path = path
	*h.contents = string(data)
	return nil
}
