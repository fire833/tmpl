/*
*	Copyright (C) 2025 Kendall Tauser
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
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewMCPTOOLCommand() *cobra.Command {
	type cmdOpts struct {
		Output      *string
		Header      *string
		Name        *string
		Package     *string
		Readonly    bool
		Destructive bool
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func New{{ .Name }}Tool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("{{ .Name }}", mcp.WithDescription(""), 
			mcp.WithReadOnlyHintAnnotation({{ if .Readonly }}true{{ else }}false{{ end }}),
			mcp.WithDestructiveHintAnnotation({{ if .Destructive }}true{{ else }}false{{ end }}),
		),
		Handler: {{ .Name | lower }}Tool,
	}
}

func {{ .Name | lower }}Tool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError("not implemented"), nil
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "mcptool",
		Aliases: []string{"mcpt"},
		Short:   "Generate boilerplate for creating new mcptool templates.",
		Long:    "",
		Version: "0.2.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.RenderTemplateToOutput(cmd.OutOrStdout(), *opts.Output, "mcptool", tmpl, opts)
		},
	}

	set := pflag.NewFlagSet("mcptool", pflag.ExitOnError)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:  utils.HeaderFlag(set),
		Name:    set.StringP("name", "n", "", "Specify the name of this tool."),
		Package: set.StringP("package", "p", "templates", "Specify the output package for this new tool template being created."),
	}

	// Bools are bound by value, since a template treats any non-nil *bool as true.
	opts = o
	set.BoolVarP(&opts.Destructive, "destructive", "d", false, "Is this tool capable of performing destructive actions?")
	set.BoolVarP(&opts.Readonly, "readonly", "r", false, "Is this tool readonly?")

	cmd.Flags().AddFlagSet(set)

	return cmd
}
