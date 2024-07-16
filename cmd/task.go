package cmd

import (
	"bcli/api"
	"bcli/view"
	"fmt"
	"os"

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
            fmt.Println(err)
		}

		fmt.Printf("Open tasks (%d):\n", len(tasks.Open))
		for _, task := range tasks.Open {
            view.PrintTask(&task)
		}

		fmt.Printf("In progress tasks (%d):\n", len(tasks.InP))
		for _, task := range tasks.InP {
            view.PrintTask(&task)
		}

		fmt.Printf("Done tasks (%d):\n", len(tasks.Done))
		for _, task := range tasks.Done {
            view.PrintTask(&task)
		}
	},
}

var viewTaskCmd = &cobra.Command{
	Use:   "view",
	Short: "View a task",
	Long:  `View a task`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		taskId := args[0]
		task, err := api.GetTask(taskId)
		if err != nil {
            fmt.Println(err)
            os.Exit(1)
		}

		fmt.Println(task)
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task",
	Long:  `Update a task`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var updateTaskContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Update task content",
	Long:  `Update task content`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		taskId := args[0]
		err := api.UpdateTaskContent(taskId, "nt..")
		if err != nil {
            fmt.Println(err)
            os.Exit(1)
		}

		fmt.Println("Task content updated")
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	Long:  `Create a task`,
	Run: func(cmd *cobra.Command, args []string) {
		api.CreateTask("template", "title", "content")
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(listTaskCmd)
	taskCmd.AddCommand(createTaskCmd)
	taskCmd.AddCommand(viewTaskCmd)
	taskCmd.AddCommand(updateTaskCmd)

	updateTaskCmd.AddCommand(updateTaskContentCmd)
}
