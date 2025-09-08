package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewRUSTDOCKERFILECommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Header  *string
		AppName *string
		Repo    *string
		Author  *string
	}

	const tmpl string = `
{{ .Header }}

FROM docker.io/rust:latest AS builder

ADD src /usr/src/{{ .AppName }}/src
ADD rs_crates /usr/src/{{ .AppName }}/rs_crates
ADD Cargo.lock /usr/src/{{ .AppName }}
ADD Cargo.toml /usr/src/{{ .AppName }}
WORKDIR /usr/src/{{ .AppName }}

RUN wget -O /usr/local/bin/dumb-init https://static.tauser.us/init/dumb-init-x86 && chmod +x /usr/local/bin/dumb-init

# Compile for musl so it works on alpine: https://stackoverflow.com/questions/59766239/how-to-build-a-rust-app-free-of-shared-libraries#59766875
RUN RUSTFLAGS='-C link-arg=-s' cargo build --release --bin {{ .AppName }}

FROM docker.io/ubuntu:noble

COPY --from=builder /usr/local/bin/dumb-init /usr/local/bin/dumb-init
COPY --from=builder /usr/src/{{ .AppName }}/target/release/{{ .AppName }} /usr/local/bin/{{ .AppName }}
USER 1000:1000

LABEL org.opencontainers.image.author="{{ .Author }}"
LABEL org.opencontainers.image.title="{{ .AppName }}"
LABEL org.opencontainers.image.source="{{ .Repo }}"

ENTRYPOINT [ "dumb-init", "{{ .AppName }}" ]
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "rustdockerfile",
		Aliases: []string{"dockerfilers", "drs", "dfrs"},
		Short:   "Generate boilerplate for creating new RustDockerfile templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("RustDockerfile").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("RustDockerfile", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:  &str,
		AppName: set.StringP("name", "n", "", "Specify the name of the app/binary of the app."),
		Repo:    set.StringP("repo", "r", "github.com/fire833/KTAPI", "Specify the repo to place in the image metadata."),
		Author:  set.StringP("author", "a", "Kendall Tauser", "Specify the author of the image."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
