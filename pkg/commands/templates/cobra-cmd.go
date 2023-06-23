package templates

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewCOBRACommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Name    *string
		Package *string
	}

	const tmpl string = `
package {{.Package}}

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)
	
func New{{.Name|upper}}Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "{{.Name}}",
		Aliases: []string{},
		Short:   "",
		Long:    "",
		Version: "0.0.1",
		Example: "",
		RunE: func(cmd *cobra.Command, args []string) error {},
	}

	set := pflag.NewFlagSet("{{.Name}}", pflag.ExitOnError)

	cmd.Flags().AddFlagSet(set)
	cmd.AddCommand()

	return cmd
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "cobra-cmd",
		Aliases: []string{},
		Short:   "Generate boilerplate for creating new go cobra commands.",
		Long:    ``,
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("cobra-cmd").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("cobra-cmd", pflag.ExitOnError)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "cobra.go", "Specify the output location for this template. If set to '-', will print to stdout."),
		Name:    set.StringP("name", "n", "cobra", "Specify the name for this template, this will be plugged into the template and be the command name."),
		Package: set.StringP("package", "p", "templates", "Specify the output package for this new template being created."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
