package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Version is used by the build system.
var Version string

// The log flag value.
var l string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "tfe-cli",
	Short:   "Manage TFE from the command line.",
	Long:    `Manage TFE from the command line.`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		if err := setUpLogs(l); err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&l, "log", "l", "", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringP("organization", "o", "", "terraform organization")
	rootCmd.PersistentFlags().StringP("token", "t", "", "terraform token")
}

// SetUpLogs sets the log level.
func setUpLogs(level string) error {
	// Read the log level
	//  1. from the CLI first
	//  2. then the ENV vars
	//  3. then use the default value.
	if level == "" {
		level = os.Getenv("TFE_LOG_LEVEL")
		if level == "" {
			level = logrus.WarnLevel.String()
		}
	}

	// Parse the log level.
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	// Set the log level.
	logrus.SetLevel(lvl)
	return nil
}
