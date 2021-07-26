package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tfbackend",
		Short:         "tfbackend is a CLI tool to create terraform backend to cloud.",
		Long:          `tfbackend is a CLI tool to create terraform backend to cloud.`,
		SilenceErrors: true,
	}
	cobra.OnInitialize(initConfig)

	cmd.AddCommand(NewCmdAws())
	cmd.AddCommand(NewCmdCompletion())

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewCmdRoot()
	if err := cmd.Execute(); err != nil {
		printErrorRed(err)
	}
}

func init() {
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".tfbackend" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfbackend")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
