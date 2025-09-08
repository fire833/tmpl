package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewGODOCKERFILECommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Header  *string
		AppName *string
		Repo    *string
		Author  *string
	}

	const tmpl string = `
{{ .Header }}

FROM docker.io/golang:alpine AS builder

RUN apk add --update make wget git

ARG VERSION=
ARG COMMIT=
ARG DATE=

# Get dumb-init in image
RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 && chmod +x /usr/local/bin/dumb-init
# Install mage
RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go

ADD . /usr/src/{{ .AppName }}
WORKDIR /usr/src/{{ .AppName }}

# Only build for amd64 and {{ .AppName }}.
RUN mage build{{ .AppName }}

FROM docker.io/alpine:latest

COPY --from=builder /usr/local/bin/dumb-init /usr/local/bin/dumb-init
COPY --from=builder /usr/src/{{ .AppName }}/bin/{{ .AppName }} /usr/local/bin/{{ .AppName }}
USER 1000:1000

LABEL org.opencontainers.image.author="{{ .Author }}"
LABEL org.opencontainers.image.title="{{ .AppName }}"
LABEL org.opencontainers.image.source="{{ .Repo }}"
LABEL org.opencontainers.image.version=${VERSION}
LABEL org.opencontainers.image.commit=${COMMIT}
LABEL org.opencontainers.image.builddate=${DATE}

ENTRYPOINT [ "dumb-init", "{{ .AppName }}" ]
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "godockerfile",
		Aliases: []string{"dockerfilego", "dgo", "dfgo"},
		Short:   "Generate boilerplate for creating new GoDockerfile templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("GoDockerfile").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("GoDockerfile", pflag.ExitOnError)

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
