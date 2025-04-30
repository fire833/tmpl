package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewPULUMIGOK8SAPPCommand() *cobra.Command {
	type cmdOpts struct {
		Output  *string
		Header  *string
		Name    *string
		Package *string
		Args    *bool
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import (
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type {{ .Name }}App struct {
	pulumi.ResourceState

	resources *appresources.SimpleAppResources
}

{{- if .Args }}
type {{ .Name }}AppArgs struct {
}
{{- end }}

func new{{.Name}}Container(name namers.AppNamer, volumeMounts v1.VolumeMountArrayInput) *v1.ContainerArgs {
	return k8sutils.NewContainer(name, "", versions.{{.Name}}AppVersion, consts.ImagePullPolicyIfNotPresent, []string{}, []string{}, map[string]uint16{}, k8sutils.NewPodEnvVarsWithBasics(map[string]string{}), k8sutils.NewContainerResources("", "", "", ""), k8sutils.NewSecuritySettingsLockdown(), volumeMounts, nil, nil, nil)
}

func new{{.Name}}Pod(name namers.AppNamer, volumes v1.VolumeArrayInput, volumeMounts v1.VolumeMountArrayInput) *v1.PodSpecArgs {
	bools := k8sutils.PodBooleans(0)
	return k8sutils.NewPod(name, "", "", bools, volumes, v1.ContainerArray{new{{.Name}}Container(name, volumeMounts)}, v1.ContainerArray{})
}

// Instantiate a new instance of {{ .Name }}App
func New{{.Name}}App(ctx *pulumi.Context, name namers.AppNamer, namespace string{{ if .Args }}, args {{ .Name }}AppArgs{{ end }}, opts ...pulumi.ResourceOption) *{{ .Name }}App {
	return &{{ .Name }}App{
		resources: &appresources.SimpleAppResources{},
	}
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "pulumigoK8sApp",
		Aliases: []string{"pgk8s"},
		Short:   "Generate boilerplate for creating new pulumigoK8sApp templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("pulumigoK8sApp").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("pulumigoK8sApp", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:  set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:  &str,
		Name:    set.StringP("name", "n", "Example", "Specify the name of the component resource you wish to create."),
		Package: set.StringP("package", "p", "unknown", "Specify the package name the component resource is a part of."),
		Args:    set.BoolP("args", "a", false, "Specify whether an additional arguments struct should be generated for your component resource."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
