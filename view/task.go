package view

import (
	"bcli/api"
	"fmt"
)

const taskDetailPrefix = "https://blueprint.cyberlogitec.com.vn/UI_PIM_001_1/"
const ansiBold = "\033[1m"
const ansiReset = "\033[0m"

func PrintTaskInLine(task *api.Task) {
	var seqNo = task.SeqNo
	var taskNm = task.TaskNm

	fmt.Printf(" ⭐ %s%d%s %s\n", ansiBold, seqNo, ansiReset, taskNm)
}

func PrintTaskList(tasks *api.Tasks) {
	fmt.Printf("\n📂 %sOpen%s (%d)\n", ansiBold, ansiReset, len(tasks.Open))
	for _, task := range tasks.Open {
		PrintTaskInLine(&task)
	}

	fmt.Printf("\n⏳ %sIn progress%s (%d)\n", ansiBold, ansiReset, len(tasks.InP))
	for _, task := range tasks.InP {
		PrintTaskInLine(&task)
	}

	fmt.Printf("\n✅ %sDone%s (%d)\n", ansiBold, ansiReset, len(tasks.Done))
	for _, task := range tasks.Done {
		PrintTaskInLine(&task)
	}

	fmt.Println()
}

func PrintTaskListInMarkdown(tasks *api.Tasks) {
	var markdown = "# Task list\n"

	markdown += "## 📂 Open\n"

	for _, task := range tasks.Open {
		markdown += fmt.Sprintf("- #%d %s \n  - [](%s)\n", task.SeqNo, task.TaskNm, taskDetailPrefix+task.ReqID)
	}

	markdown += "\n## ⏳ In process\n"
	for _, task := range tasks.InP {
		markdown += fmt.Sprintf("- #%d %s \n  - [](%s)\n", task.SeqNo, task.TaskNm, taskDetailPrefix+task.ReqID)
	}

	markdown += "\n## ✅ Done\n"
	for _, task := range tasks.Done {
		markdown += fmt.Sprintf("- #%d %s \n  - [](%s)\n", task.SeqNo, task.TaskNm, taskDetailPrefix+task.ReqID)
	}

	RenderMarkdown(markdown)
}

// print in one list, don't care about status
func PrintSimpleTaskList(tasks *api.Tasks) {
	for _, task := range tasks.Open {
		PrintTaskInLine(&task)
	}

	for _, task := range tasks.InP {
		PrintTaskInLine(&task)
	}

	for _, task := range tasks.Done {
		PrintTaskInLine(&task)
	}
}
