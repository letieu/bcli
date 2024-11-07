package cmd

import (
	"fmt"
	"github.com/letieu/bcli/api"
	"github.com/letieu/bcli/paser"
	"github.com/letieu/bcli/view"
	"os"
	"os/exec"
	"time"

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

		if markdown, _ := cmd.Flags().GetBool("markdown"); markdown {
			view.PrintTaskListInMarkdown(tasks)
			return
		}

		if simple, _ := cmd.Flags().GetBool("simple"); simple {
			view.PrintSimpleTaskList(tasks)
			return
		}

		view.PrintTaskList(tasks)
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

		taskId, err := getTaskId(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if web, _ := cmd.Flags().GetBool("web"); web {
			link := fmt.Sprintf("%s%s", taskDetailPrefix, taskId)
			fmt.Printf("Opening %s\n in web browser", link)
			cmd := exec.Command("xdg-open", link)
			cmd.Run()
			return
		}

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

		if simple, _ := cmd.Flags().GetBool("simple"); simple {
			fmt.Println(mdText)
		} else {
			view.RenderMarkdown(mdText)
		}
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

		taskId, err := getTaskId(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

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
		content, _ := cmd.Flags().GetString("content")
		template, _ := cmd.Flags().GetString("template")

		templateFile, err := os.ReadFile(template)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// due date is now + 1 week
		dueDate := time.Now().AddDate(0, 0, 7).Format("20060102")

		payload := paser.GetNewTaskPayload(string(templateFile), map[string]string{
			"title": title,
			"date":  dueDate,
		})

		createdRes, err := api.CreateTask([]byte(payload))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Task created")
		fmt.Println(createdRes.ReqID)

        if content != "" {
		    task, err := api.GetTaskDetail(createdRes.ReqID)
            err = api.UpdateTaskContent(task, content)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
	},
}

var updateHeadlessTaskCmd = &cobra.Command{
	Use:   "update-headless",
	Short: "Update a task in headless mode",
	Long:  `Update a task in headless mode without opening an editor`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		taskId, err := getTaskId(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		newTitle, _ := cmd.Flags().GetString("title")
		newContent, _ := cmd.Flags().GetString("content")

		if newTitle == "" && newContent == "" {
			fmt.Println("Please provide a new title or new content")
			return
		}

		currentTask, err := api.GetTaskDetail(taskId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if newTitle != "" {
			err = api.UpdateTaskTitle(currentTask, newTitle)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		if newContent != "" {
			err = api.UpdateTaskContent(currentTask, newContent)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		fmt.Println("Task updated")
	},
}

var branchName = &cobra.Command{
	Use:   "branch",
	Short: "Get branch name",
	Long:  `Get branch name from task title`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		taskId, err := getTaskId(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		task, err := api.GetTaskDetail(taskId)

		taskTitle := task.DetailReqVO.ReqTitNm
		taskNo := task.DetailReqVO.SeqNo

		branchName := paser.GetGitBranchName(taskNo, taskTitle)

		if checkout, _ := cmd.Flags().GetBool("checkout"); checkout {
			cmd := exec.Command("git", "checkout", "-b", branchName)
			cmd.Run()
		}
		fmt.Println(branchName)
	},
}

func init() {
	rootCmd.AddCommand(taskCmd)

	taskCmd.AddCommand(listTaskCmd)
	taskCmd.AddCommand(createTaskCmd)
	taskCmd.AddCommand(viewTaskCmd)
	taskCmd.AddCommand(updateTaskCmd)
	taskCmd.AddCommand(updateHeadlessTaskCmd)
	taskCmd.AddCommand(branchName)

	listTaskCmd.Flags().BoolP("markdown", "m", false, "Render task list in markdown")
	listTaskCmd.Flags().BoolP("simple", "s", false, "Render task list in simple mode")

	viewTaskCmd.Flags().BoolP("markdown", "m", true, "Render task content in markdown")
	viewTaskCmd.Flags().BoolP("simple", "s", false, "Render task content in simple mode")
	viewTaskCmd.Flags().BoolP("web", "w", false, "Open task in web browser")

	createTaskCmd.Flags().StringP("title", "t", "", "Task title")
	createTaskCmd.Flags().StringP("template", "T", "", "Task template")
	createTaskCmd.MarkFlagRequired("title")
	createTaskCmd.MarkPersistentFlagFilename("template")
	createTaskCmd.MarkFlagRequired("template")
    createTaskCmd.Flags().StringP("content", "c", "", "Task content")

	updateHeadlessTaskCmd.Flags().StringP("title", "t", "", "New task title")
	updateHeadlessTaskCmd.Flags().StringP("content", "c", "", "New task content")
	updateHeadlessTaskCmd.MarkFlagRequired("title")
	updateHeadlessTaskCmd.MarkFlagRequired("content")

	branchName.Flags().BoolP("checkout", "c", false, "Checkout branch")
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

func getTaskId(id string) (string, error) {
	if isSeqNo := len(id) < 6; isSeqNo == false {
		return id, nil
	}

	taskId, err := api.SearchTaskByNo(id)
	if err != nil {
		return "", err
	}

	return taskId, nil
}
