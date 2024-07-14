package cmd

import (
	"bcli/api"
	"fmt"

	"github.com/spf13/cobra"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task management",
	Long:  `Task management`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listTaskCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := api.ListTasks()

		if err != nil {
			panic(err)
		}

		fmt.Printf("Open tasks (%d):\n", len(tasks.Open))
		for _, task := range tasks.Open {
			fmt.Println(" - " + task.TaskNm)
		}

		fmt.Printf("In progress tasks (%d):\n", len(tasks.InP))
		for _, task := range tasks.InP {
			fmt.Println(" - " + task.TaskNm)
		}

		fmt.Printf("Done tasks (%d):\n", len(tasks.Done))
		for _, task := range tasks.Done {
			fmt.Println(" - " + task.TaskNm)
		}
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(listTaskCmd)
}
