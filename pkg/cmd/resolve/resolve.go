package resolve

import (
	"fmt"
	"os"

	"github.com/jenkins-x-plugins/jx-mink/pkg/cmd/initcmd"
	"github.com/jenkins-x-plugins/jx-mink/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	info = termcolor.ColorInfo

	cmdLong = templates.LongDesc(`
		Builds any images and resolves the references in the yaml files in .mink.yaml

		If there is no .mink.yaml and there is no Dockerfile or build pack overrides.yaml file then this step does nothing
`)

	cmdExample = templates.Examples(`
		# build any images and resolve their references in the YAML files
		%s resolve
	`)
)

// Options the options for the command
type Options struct {
	initcmd.Options
	Args          []string
	CommandRunner cmdrunner.CommandRunner
}

// NewCmdMinkResolve creates a command object for the command
func NewCmdMinkResolve() (*cobra.Command, *Options) {
	o := &Options{}

	cmd := &cobra.Command{
		Use: "resolve",

		Short:   "Builds any images and resolves the references in the yaml files in .mink.yaml",
		Long:    cmdLong,
		Example: fmt.Sprintf(cmdExample, rootcmd.BinaryName),
		Run: func(cmd *cobra.Command, args []string) {
			o.Args = args
			err := o.Run()
			helper.CheckErr(err)
		},
	}
	cmd.Flags().StringVarP(&o.Dir, "dir", "d", ".", "the directory to use")
	cmd.Flags().StringVarP(&o.Dockerfile, "dockerfile", "f", "Dockerfile", "the name of the dockerfile to use")
	return cmd, o
}

// Run transforms the YAML files
func (o *Options) Run() error {
	err := o.Options.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to initialise mink")
	}

	if !o.MinkEnabled {
		return nil
	}

	if o.CommandRunner == nil {
		o.CommandRunner = cmdrunner.DefaultCommandRunner
	}

	env := o.createMinkEnv()

	c := &cmdrunner.Command{
		Name: "mink",
		Args: append([]string{"resolve"}, o.Args...),
		Out:  os.Stdout,
		Err:  os.Stderr,
		In:   os.Stdin,
		Env:  env,
	}
	_, err = o.CommandRunner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s", c.CLI())
	}
	return nil
}

func (o *Options) createMinkEnv() map[string]string {
	m := map[string]string{}
	gitURL := os.Getenv("MINK_GIT_URL")
	if gitURL == "" {
		m["MINK_GIT_URL"] = os.Getenv("REPO_URL")
	}
	kanikoFlags := os.Getenv("MINK_KANIKO_FLAGS")
	if kanikoFlags == "" {
		m["MINK_KANIKO_FLAGS"] = os.Getenv("KANIKO_FLAGS")
	}
	output := os.Getenv("MINK_OUTPUT")
	if output == "" {
		m["MINK_OUTPUT"] = "."
	}
	return m
}
