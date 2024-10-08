package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewLATEXARTICLECommand() *cobra.Command {
	type cmdOpts struct {
		Output *string
		Header *string
	}

	const tmpl string = ``

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "latex-article",
		Aliases: []string{"latexa", "tex", "tex-art"},
		Short:   "Generate boilerplate for creating new Latex documents.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("latex-article").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("latex-article", pflag.ExitOnError)

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
