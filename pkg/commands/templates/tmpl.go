/*
*	Copyright (C) 2022  Kendall Tauser
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
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewTMPLTMPLCommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Name    *string
		Package *string
	}

	const tmpl string = `
package {{.Package}}

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func New{{.Name|upper}}Command() *cobra.Command {
	type cmdOpts struct {
		Output *string
		Header *string
	}

	const tmpl string = ""

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "{{.Name}}",
		Aliases: []string{},
		Short:   "Generate boilerplate for creating new {{.Name}} templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("{{.Name}}").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("{{.Name}}", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output: set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header: &str,
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "tmpl",
		Aliases: []string{"t", "tpl"},
		Short:   "Generate boilerplate for creating new tmpl templates.",
		Long:    ``,
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("tmpltmpl").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("tmpltmpl", pflag.ExitOnError)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Name:    set.StringP("name", "n", "example", "Specify the name for this template, this will be plugged into the template and be the command name."),
		Package: set.StringP("package", "p", "templates", "Specify the output package for this new template being created."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
