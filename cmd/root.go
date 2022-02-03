package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "time-tracker-rest",
	Short:   "Time Tracker REST Server",
	Version: "0.1.0",
}

func init() {
	cobra.OnInitialize(initConfig)
}
func initConfig() {

}

// Execute app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
