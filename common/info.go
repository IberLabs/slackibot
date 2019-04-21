package common

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of iBot",
	Long:  `All software has versions. This is iBot`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iBot v0.1.3")
	},
}

var rootCmd = &cobra.Command{
	Use:   "iBot",
	Short: "iBot automatic alerts for slack",
	Long: `A Fast and Flexible slack bot written in Go.
                It trigger alerts automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No arguments detected. Running out-of-the-box.")
	},
}

func Execute() bool {
	continueFlag := true

	if err := versionCmd.Execute(); err != nil {
		continueFlag = false
	}

	return continueFlag
}