package templates

import (
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
			return utils.RenderTemplateToOutput(cmd.OutOrStdout(), *opts.Output, "latex-article", tmpl, opts)
		},
	}

	set := pflag.NewFlagSet("latex-article", pflag.ExitOnError)

	o := cmdOpts{
		Output: set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header: utils.HeaderFlag(set),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
