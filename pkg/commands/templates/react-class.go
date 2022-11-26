package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewREACTCLASSCommand() *cobra.Command {
	type cmdOpts struct {
		Output *string
		Header *string
		Name   *string
	}

	const tmpl string = `
{{.Header}}

import { Component } from 'react';

interface {{.Name}}Props {

}

interface {{.Name}}State {

}

export default class {{.Name}} extends Component<{{.Name}}Props, {{.Name}}State> {

	constructor(props: {{.Name}}Props) {
		super(props);
	}

	public render() {
		return (
			<h1>Hello, world</h1>
		)
	}
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "reactClass",
		Aliases: []string{"reactc", "rc", "compclass", "rcc"},
		Short:   "Generate boilerplate for creating new reactClass templates.",
		Long:    "",
		Version: "0.0.2",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("reactClass").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("reactClass", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output: set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header: &str,
		Name:   set.StringP("compname", "n", "KTComponent", "Specify the name of the component you wish to create."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
