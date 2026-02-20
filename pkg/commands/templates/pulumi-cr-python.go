package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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
		Args      *bool
	}

	const tmpl string = `
{{.Header}}
from pulumi import ComponentResource, ResourceOptions
{{ if .Args }}
class {{ .Name }}Args:
	"""
	
	"""

	pass
{{- end }}

class {{ .Name }}(ComponentResource):
	"""
	
	"""

	def __init__(self, name: str, {{ if .Args }}args: {{ .Name }}Args, {{ end }}opts: ResourceOptions | None = None) -> None:
		super().__init__('{{.Module}}:{{.Namespace}}:{{.Name}}', name, None, opts)
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "pulumi-cr-python",
		Aliases: []string{"pcrp"},
		Short:   "Generate boilerplate for creating new PulumiPythonCR templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("PulumiPythonCR").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("PulumiPythonCR", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:    set.StringP("output", "o", "tmpl.py", "Specify the output location for this template. If set to '-', will print to stdout."),
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
