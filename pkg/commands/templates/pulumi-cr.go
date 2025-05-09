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

func NewPULUMICRCommand() *cobra.Command {
	type cmdOpts struct {
		Output    *string
		Header    *string
		Name      *string
		Package   *string
		Module    *string
		Namespace *string
		Args      *bool
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type {{.Name}} struct {
	pulumi.ResourceState
}

{{- if .Args }}
type {{ .Name }}Args struct {
}
{{- end }}

func New{{.Name}}(ctx *pulumi.Context, name string{{ if .Args }}, args {{ .Name }}Args{{ end }}, opts ...pulumi.ResourceOption) (*{{.Name}}, error) {
	c := &{{.Name}}{}
	if e := ctx.RegisterComponentResource("{{.Module}}:{{.Namespace}}:{{.Name}}", name, c, opts...); e != nil {
		return nil, e
	}

	// TODO: Add child resources here

	return c, nil
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "pulumi-cr",
		Aliases: []string{"pcr"},
		Short:   "Generate boilerplate for creating new pulumi component resource templates.",
		Long:    ``,
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("pulumi-cr").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("pulumi-cr", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:    set.StringP("output", "o", "tmpl.go", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:    &str,
		Name:      set.StringP("name", "n", "ExampleInstance", "Specify the name of the component resource you wish to create."),
		Package:   set.StringP("package", "p", "unknown", "Specify the package name the component resource is a part of."),
		Module:    set.StringP("module", "m", "unknown", "Specify the high-level module this component resource is a part of."),
		Namespace: set.String("namespace", "unknown", "Specify the namespace for this component resource within your overall stack."),
		Args:      set.BoolP("args", "a", false, "Specify whether an additional arguments struct should be generated for your component resource."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
