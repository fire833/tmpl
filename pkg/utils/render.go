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
	"bytes"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// RenderTemplate parses tmpl with the sprig function map and executes it
// against data, writing the result to w.
func RenderTemplate(w io.Writer, name, tmpl string, data any) error {
	tpl, tple := template.New(name).Funcs(sprig.TxtFuncMap()).Parse(tmpl)
	if tple != nil {
		return tple
	}

	return tpl.Execute(w, data)
}

// RenderTemplateToOutput renders tmpl and writes the result to the file at
// output, or to stdout if output is "-". The template is fully rendered before
// anything is written, so a failed render never truncates an existing file.
func RenderTemplateToOutput(stdout io.Writer, output, name, tmpl string, data any) error {
	var buf bytes.Buffer
	if e := RenderTemplate(&buf, name, tmpl, data); e != nil {
		return e
	}

	// Write to the caller's stdout rather than letting GetOutputWriter hand
	// back os.Stdout, which would then be closed.
	if output == "-" {
		_, e := buf.WriteTo(stdout)
		return e
	}

	file, filee := GetOutputWriter(output)
	if filee != nil {
		return filee
	}

	if _, e := buf.WriteTo(file); e != nil {
		file.Close()
		return e
	}

	return file.Close()
}
