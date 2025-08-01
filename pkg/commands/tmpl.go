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

package commands

import (
	"github.com/fire833/tmpl/pkg/commands/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewTMPLCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tmpl",
		Aliases: []string{"templates", "tmplt", "t", "tpl"},
		Short:   "CLI for generating boilerplate for different code things.",
		Long:    ``,
		Version: "0.1.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	set := pflag.NewFlagSet("tmpl", pflag.ExitOnError)

	cmd.Flags().AddFlagSet(set)

	cmd.AddCommand(
		templates.NewTMPLTMPLCommand(),
		templates.NewCOBRACommand(),
		templates.NewK8SDAEMONSETCommand(),
		templates.NewK8SDEPLOYMENTCommand(),
		templates.NewK8SINGRESSCommand(),
		templates.NewK8SSVCACCTCommand(),
		templates.NewREACTCLASSCommand(),
		templates.NewREACTFUNCTIONALCommand(),
		templates.NewRUSTCommand(),
		templates.NewPULUMICRCommand(),
		templates.NewSVELTECCommand(),
		templates.NewJAVAOPENGLCommand(),
		templates.NewLATEXARTICLECommand(),
		templates.NewPULUMIJAVACRCommand(),
		templates.NewPULUMIGOK8SAPPCommand(),
		templates.NewADVENTOFCODECHALLENGECommand(),
		templates.NewPULUMISTATELESSWEBAPPCommand(),
		templates.NewBURNMODULECommand(),
		templates.NewMCPTOOLCommand(),
		templates.NewMCPPromptCommand(),
	)

	return cmd
}
