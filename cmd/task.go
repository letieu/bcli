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

const taskDetailPrefix = "https://blueprint.cyberlogitec.com.vn/UI_PIM_001_1/"

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
            os.Exit(1)
		}

		markdown, _ := cmd.Flags().GetBool("markdown")
		if markdown {
			view.PrintTaskListInMarkdown(&tasks)
		} else {
			view.PrintTaskList(&tasks)
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
		link := fmt.Sprintf("%s%s", taskDetailPrefix, taskId)

		web, _ := cmd.Flags().GetBool("web")
		if web {
			cmd := exec.Command("xdg-open", link)
			cmd.Run()
			return
		}

		markdown, _ := cmd.Flags().GetBool("markdown")
		if markdown {
			task, err := api.GetTaskDetail(taskId)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			mdText, err := paser.CreateBufText(task.DetailReqVO.ReqTitNm, task.DetailReqVO.ReqCtnt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			mdText = fmt.Sprintf("# #%d\n%s", task.DetailReqVO.SeqNo, mdText)
			view.RenderMarkdown(mdText)
			return
		}

		fmt.Println(link)
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
		currentTask, err := api.GetTaskDetail(taskId)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newTitle, newContent, err := editTask(currentTask.DetailReqVO.ReqTitNm, currentTask.DetailReqVO.ReqCtnt)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Updating task")
		err = api.UpdateTaskTitle(currentTask, newTitle)
		err = api.UpdateTaskContent(currentTask, newContent)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Task updated")
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	Long:  `Create a task`,
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		template, _ := cmd.Flags().GetString("template")

		templateFile, err := os.ReadFile(template)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		payload := paser.CreatePayload(string(templateFile), map[string]string{"title": title})

		createdRes, err := api.CreateTask([]byte(payload))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Task created")
		fmt.Println(createdRes.ReqID)
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(listTaskCmd)
	taskCmd.AddCommand(createTaskCmd)
	taskCmd.AddCommand(viewTaskCmd)
	taskCmd.AddCommand(updateTaskCmd)

	listTaskCmd.Flags().BoolP("markdown", "m", false, "Render task list in markdown")

	viewTaskCmd.Flags().BoolP("web", "w", false, "Open task in web browser")
	viewTaskCmd.Flags().BoolP("markdown", "m", false, "Render task content in markdown")

	createTaskCmd.Flags().StringP("title", "t", "", "Task title")
	createTaskCmd.Flags().StringP("template", "T", "", "Task template")
	createTaskCmd.MarkFlagRequired("title")
	createTaskCmd.MarkPersistentFlagFilename("template")
	createTaskCmd.MarkFlagRequired("template")
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

func editTask(title string, contentHtml string) (string, string, error) {
	bufText, err := paser.CreateBufText(title, contentHtml)
	if err != nil {
		return "", "", err
	}

	newBufText, err := editInEditor(bufText)
	if err != nil {
		return "", "", err
	}

	if newBufText == bufText {
		return "", "", fmt.Errorf("No changes made")
	}

	newTitle, newContentHtml, err := paser.ParseBufText(newBufText)
	if err != nil {
		return "", "", err
	}

	return newTitle, newContentHtml, nil
}
