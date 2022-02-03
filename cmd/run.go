package cmd

import (
	"dcamachoj/time-tracker-rest/rest"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run REST Server",
	Run: func(cmd *cobra.Command, args []string) {
		rest.ExecuteServer()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
