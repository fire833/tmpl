/*
*	Copyright (C) 2025 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewPULUMISTATELESSWEBAPPCommand() *cobra.Command {
	type cmdOpts struct {
		Output         *string
		Header         *string
		Name           *string
		Module         *string
		Namespace      *string
		Package        *string
		Args           *bool
		NamespaceRBAC  *bool
		ClusterRBAC    *bool
		CollectMetrics *bool
		Image          *string
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import (
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type {{ .Name }}App struct {
	pulumi.ResourceState
	appresources.SimpleAppResources
}

{{- if .Args }}
type {{ .Name }}AppArgs struct {
}
{{- end }}

func new{{.Name}}Container(name namers.AppNamer, listenPort uint16, volumeMounts v1.VolumeMountArrayInput) *v1.ContainerArgs {
	return k8sutils.NewContainer(name, "{{ .Image }}", versions.{{.Name}}AppVersion, consts.ImagePullPolicyIfNotPresent, 
		[]string{}, []string{}, map[string]uint16{
			"http": listenPort,
		},
		k8sutils.NewPodEnvVarsWithBasics(map[string]string{}), 
		k8sutils.NewContainerResources("50m", "100m", "50Mi", "100Mi"), 
		k8sutils.NewSecuritySettingsLockdown(), 
		volumeMounts, nil, nil, nil)
}

func new{{.Name}}Pod(name namers.AppNamer, listenPort uint16, volumes v1.VolumeArrayInput, volumeMounts v1.VolumeMountArrayInput) *v1.PodSpecArgs {
	bools := k8sutils.PodBooleans(0){{ if .NamespaceRBAC }}.WithRBAC().WithAutomountSvcTok(){{ else if .NamespaceRBAC }}.WithRBAC().WithAutomountSvcTok(){{ end }}
	return k8sutils.NewPod(name, "", "", bools, volumes, v1.ContainerArray{new{{.Name}}Container(name, listenPort, volumeMounts)}, v1.ContainerArray{})
}

// Instantiate a new instance of {{ .Name }}App.
func New{{.Name}}App(ctx *pulumi.Context, name namers.AppNamer, namespace string{{ if .Args }}, args {{ .Name }}AppArgs{{ end }}, opts ...pulumi.ResourceOption) (*{{ .Name }}App, error) {
	c := &{{.Name}}App{}
	if e := ctx.RegisterComponentResource("{{.Module}}:{{.Namespace}}:{{.Name}}App", string(name), c, opts...); e != nil {
		return nil, e
	}

	const listenPort uint16 = 8080
	{{ if .NamespaceRBAC }}
	r, rb, sa := k8sutils.NewRoleRoleBindingSvcAcct(ctx, name, namespace, k8sutils.NewPolicyRuleArray())
	{{- end }}
	{{- if .ClusterRBAC }}
	cr, crb, sa := k8sutils.NewClusterRoleRoleBindingSvcAcct(ctx, name, namespace, k8sutils.NewPolicyRuleArray())
	{{- end }}

	c.SimpleAppResources = appresources.SimpleAppResources{
		Deployment: k8sutils.NewDeployment(ctx, name, namespace, 1, new{{.Name}}Pod(name, listenPort, v1.VolumeArray{}, v1.VolumeMountArray{})),
		Service: k8sutils.NewService(ctx, name, namespace, k8sutils.ServiceTypeClusterIP, map[string]uint16{
			"http": listenPort,
		}),
		{{- if .CollectMetrics }}
		ServiceMonitor:      k8sutils.NewServiceMonitor(ctx, name, namespace, 256, 256, 30, "metrics", "/metrics", "HTTP"),
		{{- end }}
		ServiceAccounts:     map[namers.AppNamer]*v1.ServiceAccountArgs{},
		ClusterRoles:        map[namers.AppNamer]*rbacv1.ClusterRoleArgs{},
		ClusterRoleBindings: map[namers.AppNamer]*rbacv1.ClusterRoleBindingArgs{},
		Roles:               map[namers.AppNamer]*rbacv1.RoleArgs{},
		RoleBindings:        map[namers.AppNamer]*rbacv1.RoleBindingArgs{},
	}

	return c, nil
}

func (a *{{.Name}}App) Deploy(ctx *pulumi.Context, name namers.AppNamer, opts ...pulumi.ResourceOption) error {
	return a.SimpleAppResources.Deploy(ctx, a, name, opts...)
}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "pulumi-stateless-webapp",
		Aliases: []string{"pswapp"},
		Short:   "Generate boilerplate for creating new pulumi-stateless-webapp templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("pulumi-stateless-webapp").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("PulumiStatelessWebApp", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:         set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:         &str,
		Name:           set.StringP("name", "n", "Example", "Specify the name of the component resource you wish to create."),
		Package:        set.StringP("package", "p", "unknown", "Specify the package name the component resource is a part of."),
		Module:         set.StringP("module", "m", "KTIAC", "Specify the high-level module this component resource is a part of."),
		Namespace:      set.String("namespace", "k8s", "Specify the namespace for this component resource within your overall stack."),
		Args:           set.BoolP("args", "a", false, "Specify whether an additional arguments struct should be generated for your component resource."),
		NamespaceRBAC:  set.BoolP("nsrbac", "r", true, "Toggle whether this application will need Namespaced RBAC permissions."),
		ClusterRBAC:    set.BoolP("crbac", "c", true, "Toggle whether this application will need cluster-wide RBAC permissions."),
		CollectMetrics: set.BoolP("metrics", "s", true, "Toggle whether this application should have a servicemonitor created alongside it for collecting metrics."),
		Image:          set.StringP("image", "i", "", "Specify the image for this application."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
