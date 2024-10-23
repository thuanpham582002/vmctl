package cmd

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"vmctl/cmd/config"
	"vmctl/cmd/create"
	del "vmctl/cmd/delete"
	"vmctl/cmd/execute"
	"vmctl/cmd/restart"
	"vmctl/cmd/start"
	"vmctl/cmd/stop"
	"vmctl/global"
	"vmctl/util/resource"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:     "vmctl",
		Short:   "A CLI for managing virtual machines fit with lima",
		Version: "0.0.1",
		Run:     runHelp,
	}
)

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		handleExitCoder(err)
		logrus.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vmctl/config.yaml)")

	rootCmd.AddCommand(
		config.NewCmdConfig(),
		create.NewCmdCreate(),
		del.NewCmdDelete(),
		restart.NewCmdRestart(),
		start.NewCmdStart(),
		stop.NewCmdStop(),
		execute.NewCmdExecute(),
	)
}

func initConfig() {
	if cfgFile != "" {
		viper.Set("current-context", cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.vmctl")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		// The config file is optional, we shouldn't exit when the config is not found
		break
	default:
		vmManager, err := resource.GetVmManager()
		if err != nil {
			logrus.Error(err)
		}
		global.VmManager = vmManager
		_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			logrus.Error(err)
		}
	}
}

type ExitCoder interface {
	error
	ExitCode() int
}

func handleExitCoder(err error) {
	if err == nil {
		return
	}

	var exitErr ExitCoder
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode()) //nolint:revive // it's intentional to call os.Exit in this function
		return
	}
}
