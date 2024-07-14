package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bcli",
	Short: "Bcli is command line interface for Blueprint.",
	Long:  `Bcli is command line interface for Blueprint.`,
	Run: func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.Flags().GetBool("version"); version {
			fmt.Println("Blueprint CLI v0.1.0")
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
