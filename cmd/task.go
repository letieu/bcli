package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/letieu/bcli/api"
	"github.com/letieu/bcli/paser"
	"github.com/letieu/bcli/view"

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

var addPointCmd = &cobra.Command{
	Use:   "add-point",
	Short: "Add a point to a task",
	Long:  `Add a point to a task`,
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

		templatePath, err := cmd.Flags().GetString("template")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		volume, err := cmd.Flags().GetInt("volume")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		currentTask, err := api.GetTaskDetail(taskId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data := map[string]any{
			"ReqId":      currentTask.DetailReqVO.ReqID,
			"Volume":     volume,
			"TotalPoint": int(volume)*50 + 30,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		payload := buf.Bytes()
		err = api.UpdateTaskPoint(currentTask, payload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Point added")
	},
}

var addTimeWorkCmd = &cobra.Command{
	Use:   "add-time",
	Short: "Add time work to a task",
	Long:  `Add time work to a task`,
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

		templatePath, err := cmd.Flags().GetString("template")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		hour, err := cmd.Flags().GetInt("hour")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		date, err := cmd.Flags().GetString("date")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		currentTask, err := api.GetTaskDetail(taskId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		parsedDate, err := time.Parse("20060102", date)
		if err != nil {
			fmt.Println("Invalid date format. Please use YYYYMMDD.")
			os.Exit(1)
		}
		formattedDate := parsedDate.Format("Jan 02, 2006")

		data := map[string]any{
			// Req_001
			"ReqId": currentTask.DetailReqVO.ReqID,
			// 20241107
			"WrkDt": date,
			// 2 Hour
			"WrkTm": fmt.Sprintf("%d Hour", hour),
			//Nov 07, 2024
			"Dt": formattedDate,
			// 120
			"ActEfrtMnt": hour * 60,
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		payload := buf.Bytes()
		err = api.AddTimeWork(currentTask, payload)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Time work added")
	},
}

var addFileCmd = &cobra.Command{
	Use:   "add-file",
	Short: "Add file to a task",
	Long:  `Add file to a task`,
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

		filePath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		currentTask, err := api.GetTaskDetail(taskId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		monthStr := time.Now().Format("200601")
		bizFolder := fmt.Sprintf("/PIM_REQ/%s/PRQ", monthStr)
		childFolder := "PRQ"

		uploadRes, err := api.UploadFile(filePath, bizFolder, childFolder)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fileStat, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fileSize := fileStat.Size()
		fileSizeStr := fmt.Sprintf("%.1f KB", float64(fileSize)/1024)

		err = api.AddFileToTask(
			currentTask,
			uploadRes.LstFlNm[0],
			fileSizeStr,
			fmt.Sprintf("%s/%s", uploadRes.BizFolder, uploadRes.LstFlNm[0]),
			uploadRes.BizFolder,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("File added")
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

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a task",
	Long:  `Submit a task`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a task ID")
			return
		}

		comment, _ := cmd.Flags().GetString("comment")
		if comment == "" {
			fmt.Println("Please provide a comment")
			return
		}

		taskId, err := getTaskId(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		taskDetail, err := api.GetTaskDetail(taskId)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cmtKey, err := api.GenereateCommentKey()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = api.SubmitTask(taskDetail, cmtKey, comment)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        fmt.Println("Task submitted")
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
	taskCmd.AddCommand(addPointCmd)
	taskCmd.AddCommand(addTimeWorkCmd)
	taskCmd.AddCommand(addFileCmd)
	taskCmd.AddCommand(submitCmd)

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

	addPointCmd.Flags().IntP("volume", "v", 0, "Point volume")
	addPointCmd.Flags().StringP("template", "T", "", "Point template")
	addPointCmd.MarkFlagRequired("volume")
	addPointCmd.MarkFlagRequired("template")
	addPointCmd.MarkPersistentFlagFilename("template")

	addTimeWorkCmd.Flags().IntP("hour", "H", 0, "Hour")
	addTimeWorkCmd.Flags().StringP("date", "d", "", "Date")
	addTimeWorkCmd.Flags().StringP("template", "T", "", "Time work template")
	addTimeWorkCmd.MarkFlagRequired("hour")
	addTimeWorkCmd.MarkFlagRequired("date")
	addTimeWorkCmd.MarkFlagRequired("template")
	addTimeWorkCmd.MarkPersistentFlagFilename("template")

	addFileCmd.Flags().StringP("file", "f", "", "File path")
	addFileCmd.MarkFlagRequired("file")

	branchName.Flags().BoolP("checkout", "c", false, "Checkout branch")

    submitCmd.Flags().StringP("comment", "c", "", "Comment")
    submitCmd.MarkFlagRequired("comment")
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
