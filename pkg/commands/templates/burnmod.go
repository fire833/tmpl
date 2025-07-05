package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewBURNMODULECommand() *cobra.Command {
	type cmdOpts struct {
		Output *string
		Header *string
		Name   *string
	}

	const tmpl string = `
{{.Header}}

use burn::{config::Config, module::Module, prelude::Backend, tensor::{Tensor, Float}};

#[derive(Debug, Config)]
pub struct {{ .Name }}Config {
}

impl {{ .Name }}Config {
	pub fn init<B: Backend>(&self, device: &B::Device) -> {{ .Name }}<B> {
        {{ .Name }} {}
    }
}

#[derive(Debug, Module)]
pub struct {{ .Name }}<B: Backend> {
}

impl<B: Backend> {{ .Name }}<B> {
	pub fn forward(&self) -> Tensor<B, 2, Float> {
		todo!()
	}
}

#[test]
fn test_forward() {}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "burnmodule",
		Aliases: []string{"burnmod", "burn"},
		Short:   "Generate boilerplate for creating new burnmodule templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("burnmodule").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("burnmodule", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output: set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header: &str,
		Name:   set.StringP("name", "n", "", "Specify the name of this module."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
