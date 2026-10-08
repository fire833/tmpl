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

package templates

import (
	"io"
	"testing"

	"github.com/spf13/cobra"
)

// render runs cmd with args, sending its output to stdout rather than a file,
// and fails t if the template does not render.
func render(t *testing.T, cmd *cobra.Command, args ...string) {
	t.Helper()

	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs(append([]string{"--output", "-"}, args...))
	if e := cmd.Execute(); e != nil {
		t.Fatalf("%s %v: %v", cmd.Name(), args, e)
	}
}
