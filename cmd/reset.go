package cmd

import (
	"dcamachoj/time-tracker-rest/dbx"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset DB Schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbx.ExecuteReset()
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
