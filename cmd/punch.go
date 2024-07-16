package cmd

import (
	"bcli/api"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var punch = &cobra.Command{
	Use:   "punch",
	Short: "Punch in/out",
	Long:  `Punch in/out`,
	Run: func(cmd *cobra.Command, args []string) {
        err := api.Punch()

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
	},
}

func init() {
	rootCmd.AddCommand(punch)
}
