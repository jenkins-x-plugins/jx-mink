package cmd

import (
	"github.com/jenkins-x-plugins/jx-mink/pkg/cmd/initcmd"
	"github.com/jenkins-x-plugins/jx-mink/pkg/cmd/resolve"
	"github.com/jenkins-x-plugins/jx-mink/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Use:   rootcmd.TopLevelCommand,
		Short: "commands for stashing results",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}

	cmd.AddCommand(cobras.SplitCommand(initcmd.NewCmdMinkInit()))
	cmd.AddCommand(cobras.SplitCommand(resolve.NewCmdMinkResolve()))
	return cmd
}
