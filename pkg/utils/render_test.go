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
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	name := "widget"
	count := 3

	data := struct {
		Name  *string
		Count *int
		Plain string
	}{
		Name:  &name,
		Count: &count,
		Plain: "plain",
	}

	tests := []struct {
		name    string
		tmpl    string
		want    string
		wantErr bool
	}{
		{name: "plain field", tmpl: "{{ .Plain }}", want: "plain"},
		{name: "pointer fields are dereferenced", tmpl: "{{ .Name }}-{{ .Count }}", want: "widget-3"},
		{name: "sprig functions are available", tmpl: "{{ .Name | upper }}", want: "WIDGET"},
		{name: "empty template", tmpl: "", want: ""},
		{name: "parse error", tmpl: "{{ .Name ", wantErr: true},
		{name: "unknown field", tmpl: "{{ .Missing }}", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			e := RenderTemplate(&out, "test", tt.tmpl, data)
			if tt.wantErr {
				if e == nil {
					t.Fatalf("expected an error, got output %q", out.String())
				}
				return
			}

			if e != nil {
				t.Fatalf("unexpected error: %v", e)
			}

			if got := out.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
