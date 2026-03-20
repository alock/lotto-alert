package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/alock/lotto-alert/config"
)

var testMode bool

var rootCmd = &cobra.Command{
	Use:   "lotto-alert",
	Short: "PA Pick 3 Evening lottery alert system",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.LoadConfigs()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&testMode, "test", false, "run commands with fake test data")
}
