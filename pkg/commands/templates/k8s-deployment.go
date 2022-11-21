
package templates

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewK8SDEPLOYMENTCommand() *cobra.Command {
	type cmdOpts struct {
		output *string
	}

	const tmpl string = ""

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "k8sDeployment",
		Aliases: []string{},
		Short:   "Generate boilerplate for creating new k8sDeployment templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("k8sDeployment").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("k8sDeployment", pflag.ExitOnError)

	o := cmdOpts{
		output: set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
