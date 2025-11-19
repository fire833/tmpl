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
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewMCPPromptCommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Header  *string
		Name    *string
		Package *string
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func New{{ .Name }}Prompt() server.ServerPrompt {
	return server.ServerPrompt{
		Prompt:  mcp.NewPrompt("{{ .Name }}", mcp.WithPromptDescription("")),
		Handler: new{{ .Name }}Prompt,
	}
}

func new{{ .Name }}Prompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return mcp.NewGetPromptResult("unimplemented", []mcp.PromptMessage{

	}), nil
}

func New{{ .Name }}PromptText() string {
	return ""
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "mcpprompt",
		Aliases: []string{"mcpp"},
		Short:   "Generate boilerplate for creating new mcpp templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("mcpp").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("mcpp", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:  &str,
		Name:    set.StringP("name", "n", "", "Specify the name of this prompt."),
		Package: set.StringP("package", "p", "templates", "Specify the output package for this new prompt template being created."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
