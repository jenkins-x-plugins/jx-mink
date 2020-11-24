package resolve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jenkins-x-plugins/jx-mink/pkg/cmd/initcmd"
	"github.com/jenkins-x-plugins/jx-mink/pkg/rootcmd"
	"github.com/jenkins-x-plugins/jx-mink/pkg/variables"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/templates"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/termcolor"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	useKaniko = "PLAIN_KANIKO"
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
	Env           map[string]string
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
	var err error
	path := filepath.Join(o.Dir, ".jx", "variables.sh")
	o.Env, err = variables.ParseVariables(path)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s", path)
	}

	err = o.copyKanikoDockerSecrets()
	if err != nil {
		return errors.Wrapf(err, "failed to copy kaniko docker secrets")
	}

	err = o.createMinkEnv(o.Env)
	if err != nil {
		return errors.Wrapf(err, "failed to ")
	}

	useKaniko := o.getEnv(useKaniko)
	if useKaniko == "true" || useKaniko == "yes" {
		return o.invokeKaniko()
	}

	err = o.Options.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to initialise mink")
	}

	if !o.MinkEnabled {
		return nil
	}

	if o.CommandRunner == nil {
		o.CommandRunner = cmdrunner.DefaultCommandRunner
	}

	log.Logger().Infof("using environment: %v", info(o.Env))

	c := &cmdrunner.Command{
		Name: "mink",
		Args: append([]string{"resolve"}, o.Args...),
		Out:  os.Stdout,
		Err:  os.Stderr,
		In:   os.Stdin,
		Env:  o.Env,
	}
	_, err = o.CommandRunner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s", c.CLI())
	}
	return nil
}

func (o *Options) createMinkEnv(m map[string]string) error {
	home := o.getEnv("MINK_HOME")
	if home != "" {
		m["HOME"] = home
	}
	gitURL := o.getEnv("MINK_GIT_URL")
	if gitURL == "" {
		m["MINK_GIT_URL"] = o.getEnv("REPO_URL")
	}
	kanikoFlags := o.getEnv("MINK_KANIKO_FLAGS")
	if kanikoFlags == "" {
		m["MINK_KANIKO_FLAGS"] = o.getEnv("KANIKO_FLAGS")
	}
	output := o.getEnv("MINK_OUTPUT")
	if output == "" {
		m["MINK_OUTPUT"] = "."
	}
	return nil
}

func (o *Options) getEnv(key string) string {
	value := o.Env[key]
	if value == "" {
		value = os.Getenv(key)
	}
	return value
}

func (o *Options) invokeKaniko() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrapf(err, "failed to get current working directory")
	}

	image := o.getEnv("MINK_IMAGE")
	context := o.getEnv("KANIKO_CONTEXT")
	if context == "" {
		context = filepath.Join(wd, "Dockerfile")
	}

	args := []string{"--destination", image, "--context", context}
	flags := o.getEnv("KANIKO_FLAGS")
	args = append(args, strings.Split(flags, " ")...)

	c := &cmdrunner.Command{
		Name: "/kaniko/executor",
		Args: args,
		Out:  os.Stdout,
		Err:  os.Stderr,
		In:   os.Stdin,
		Env:  o.Env,
	}
	log.Logger().Infof("running command: %s", info(c.CLI()))
	_, err = o.CommandRunner(c)
	if err != nil {
		return errors.Wrapf(err, "failed to run %s", c.CLI())
	}
	return nil

}

func (o *Options) copyKanikoDockerSecrets() error {
	glob := filepath.Join("tekton", "cred-secrets", "*", ".dockerconfigjson")
	fs, err := filepath.Glob(glob)
	if err != nil {
		return errors.Wrapf(err, "failed to find tekton secrets")
	}
	if len(fs) == 0 {
		log.Logger().Warnf("failed to find docker secrets %s", glob)
		return nil
	}
	srcFile := fs[0]

	outDir := filepath.Join("kaniko", ".docker")
	err = os.MkdirAll(outDir, files.DefaultDirWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to create dir %s", outDir)
	}
	outFile := filepath.Join(outDir, "config.json")
	err = files.CopyFile(srcFile, outFile)
	if err != nil {
		return errors.Wrapf(err, "failed to copy file %s to %s", srcFile, outFile)
	}

	log.Logger().Infof("copied secret %s to %s", srcFile, outFile)
	return nil
}
