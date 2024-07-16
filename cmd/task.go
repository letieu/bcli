package cmd

import (
	"bcli/api"
	"bcli/paser"
	"bcli/view"
	"fmt"
	"os"
	"os/exec"

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
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		taskId := args[0]
		currentTask, err := api.GetTask(taskId)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskTitle := currentTask.DetailReqVO.ReqTitNm
		fmt.Println(taskTitle)

		updatedContent, updated := editTask(currentTask.DetailReqVO.ReqCtnt)
        if !updated {
            fmt.Println("No changes made")
            return
        }

		fmt.Println("Updating task content")
		err = api.UpdateTaskContent(taskId, updatedContent)
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

	updateTaskCmd.Flags().StringP("title", "t", "", "Task title")
	updateTaskCmd.Flags().StringP("content", "c", "", "Task content")
	updateTaskCmd.Flags().BoolP("edit", "e", false, "Edit task in editor")
	updateTaskCmd.MarkFlagsOneRequired("title", "content", "edit")
}

func editInEditor(content string) (string, error) {
	tmpfile, err := os.CreateTemp("", "task-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		return "", err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	updatedContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return string(updatedContent), nil
}

func editTask(contentHtml string) (string, bool) {
	taskContentMd, err := paser.HtmlToMd(contentHtml)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updatedContent, err := editInEditor(taskContentMd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updatedContentHtml, err := paser.MdToHtml(updatedContent)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if updatedContent == taskContentMd {
		return "", false
	}

    return updatedContentHtml, true
}
