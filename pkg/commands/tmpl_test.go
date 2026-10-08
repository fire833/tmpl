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

package commands

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewTMPLCommand(t *testing.T) {
	t.Run("subcommand names and aliases are unique", func(t *testing.T) {
		seen := map[string]string{}
		for _, sub := range NewTMPLCommand().Commands() {
			if sub.RunE == nil {
				t.Errorf("%s has no RunE", sub.Name())
			}

			for _, name := range append([]string{sub.Name()}, sub.Aliases...) {
				if other, ok := seen[name]; ok {
					t.Errorf("%q is claimed by both %s and %s", name, other, sub.Name())
				}
				seen[name] = sub.Name()
			}
		}
	})

	t.Run("no arguments prints help", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewTMPLCommand()
		cmd.SetOut(&out)
		cmd.SetArgs([]string{})
		if e := cmd.Execute(); e != nil {
			t.Fatalf("unexpected error: %v", e)
		}

		if !strings.Contains(out.String(), "Available Commands:") {
			t.Errorf("expected help output, got:\n%s", out.String())
		}
	})

	t.Run("dispatches to a subcommand by alias", func(t *testing.T) {
		var out bytes.Buffer
		cmd := NewTMPLCommand()
		cmd.SetOut(&out)
		cmd.SetErr(io.Discard)
		cmd.SetArgs([]string{"mcpt", "--name", "Search", "--output", "-"})
		if e := cmd.Execute(); e != nil {
			t.Fatalf("unexpected error: %v", e)
		}

		if !strings.Contains(out.String(), "func NewSearchTool() server.ServerTool {") {
			t.Errorf("expected mcptool output, got:\n%s", out.String())
		}
	})
}
