package templates

import (
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/fire833/tmpl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewGOSUBSYSTEMCommand() *cobra.Command {
	type cmdOpts struct {
		Output         *string
		Header         *string
		Package        *string
		SubsystemName  *string
		SubsystemShort *string
	}

	const tmpl string = `
{{.Header}}

package {{.Package}}

import (
	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/prometheus/client_golang/prometheus"
)

const {{ .SubsystemName }}SubsystemName string = "{{ .SubsystemShort }}"

var {{ .SubsystemShort | upper }} *{{ .SubsystemName }}Manager

type {{ .SubsystemName }}Manager struct {
	manager.DefaultSubsystem

	isInitialized bool
	isShutdown    bool
}

func New{{ .SubsystemName }}Manager() *{{ .SubsystemName }}Manager {
	return &{{ .SubsystemName }}Manager{
		isInitialized: false,
		isShutdown:    false,
	}
}

func (s *{{ .SubsystemName }}Manager) Name() string { return {{ .SubsystemName }}SubsystemName }

func (s *{{ .SubsystemName }}Manager) SetGlobal() { {{ .SubsystemShort | upper }} = s }

func (s *{{ .SubsystemName }}Manager) Initialize(reg *manager.SystemRegistrar) error {
	s.isInitialized = true
	return nil
}

// NOP PreInit
func (s *{{ .SubsystemName }}Manager) PreInit() {}

// NOP SyncStart
func (s *{{ .SubsystemName }}Manager) SyncStart() {}

// NOP PostInit
func (s *{{ .SubsystemName }}Manager) PostInit() {}

func (s *{{ .SubsystemName }}Manager) Configs() *[]*manager.ConfigValue {
	return &[]*manager.ConfigValue{}
}

func (s *{{ .SubsystemName }}Manager) Secrets() *[]*manager.SecretValue {
	return &[]*manager.SecretValue{}
}

// NOP to reload the subsystem
func (s *{{ .SubsystemName }}Manager) Reload() {}

// NOP to shutdown the subsystem
func (s *{{ .SubsystemName }}Manager) Shutdown() {
	s.isShutdown = true
}

// Return nothing since this subsystem does nothing, but you should be able to fill this out at runtime
// so the APIManager can effectively make decision on process lifecycle.
func (s *{{ .SubsystemName }}Manager) Status() *manager.SubsystemStatus {
	return &manager.SubsystemStatus{
		Name:          {{ .SubsystemName }}SubsystemName,
		IsInitialized: s.isInitialized,
		IsShutdown:    s.isShutdown,
		Meta:          nil,
	}
}

// NOP to implement prometheus Collector interface.
func (s *{{ .SubsystemName }}Manager) Describe(chan<- *prometheus.Desc) {}

// NOP to implement prometheus Collector interface.
func (s *{{ .SubsystemName }}Manager) Collect(chan<- prometheus.Metric) {}
`

	var opts cmdOpts

	cmd := &cobra.Command{
		Use:     "gosubsystem",
		Aliases: []string{"gss"},
		Short:   "Generate boilerplate for creating new gosubsystem templates.",
		Long:    "",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, oute := utils.GetOutputWriter(*opts.Output)
			if oute != nil {
				return oute
			}

			defer output.Close()

			tpl, tple := template.New("gosubsystem").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
			if tple != nil {
				return tple
			}

			return tpl.Execute(output, opts)
		},
	}

	set := pflag.NewFlagSet("gosubsystem", pflag.ExitOnError)

	data, _ := os.ReadFile(*set.StringP("header", "f", "hack/boilerplate.go.txt", "Specify an optional header to apply to generated files."))
	str := string(data)

	o := cmdOpts{
		Output:         set.StringP("output", "o", "tmpl.tmpl", "Specify the output location for this template. If set to '-', will print to stdout."),
		Header:         &str,
		Package:        set.StringP("package", "p", "templates", "Specify the output package for this new template being created."),
		SubsystemName:  set.StringP("name", "n", "", "Specify the name for this subsystem."),
		SubsystemShort: set.StringP("short", "s", "", "Specify the short name for the subsystem."),
	}

	cmd.Flags().AddFlagSet(set)
	opts = o

	return cmd
}
