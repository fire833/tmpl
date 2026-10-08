package templates

import (
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewPULUMIPYTHONCRCommand() *cobra.Command {
	type cmdOpts struct {
		Output    *string
		Header    *string
		Name      *string
		Package   *string
		Module    *string
		Namespace *string
		Args      bool
	}

	const tmpl string = `
{{.Header}}
from dataclasses import dataclass

from pulumi import ComponentResource, ResourceOptions
{{ if .Args }}
@dataclass(kw_only=True)
class {{ .Name }}Args:
	"""
	Arguments to pass to {{ .Name }}.
	"""
{{- end }}

class {{ .Name }}(ComponentResource):
	"""
	"""

	def __init__(self, name: str, {{ if .Args }}args: {{ .Name }}Args, {{ end }}opts: ResourceOptions | None = None) -> None:
		super().__init__("{{.Module}}:{{.Namespace}}:{{.Name}}", name, None, opts)
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "pulumi-cr-python",
		Aliases: []string{"pcrp"},
		Short:   "Generate boilerplate for creating new PulumiPythonCR templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.RenderTemplateToOutput(cmd.OutOrStdout(), *opts.Output, "PulumiPythonCR", tmpl, opts)
		},
	}

	set := pflag.NewFlagSet("PulumiPythonCR", pflag.ExitOnError)

	o := cmdOpts{
		Output:    set.StringP("output", "o", "tmpl.py", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:    utils.HeaderFlag(set),
		Name:      set.StringP("name", "n", "ExampleInstance", "Specify the name of the component resource you wish to create."),
		Package:   set.StringP("package", "p", "unknown", "Specify the package name the component resource is a part of."),
		Module:    set.StringP("module", "m", "unknown", "Specify the high-level module this component resource is a part of."),
		Namespace: set.String("namespace", "unknown", "Specify the namespace for this component resource within your overall stack."),
	}

	// Bools are bound by value, since a template treats any non-nil *bool as true.
	opts = o
	set.BoolVarP(&opts.Args, "args", "a", false, "Specify whether an additional arguments struct should be generated for your component resource.")

	cmd.Flags().AddFlagSet(set)

	return cmd
}
