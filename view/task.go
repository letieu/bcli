package view

import (
	"bcli/api"
	"fmt"
)

const taskDetailPrefix = "https://blueprint.cyberlogitec.com.vn/UI_PIM_001_1/"

func PrintTaskInLine(task *api.Task) {
	var seqNo = task.SeqNo
	var taskNm = task.TaskNm

	fmt.Printf(" - #%d %s [%s]\n", seqNo, taskNm, task.ReqID)
}

func PrintTaskList(tasks *api.Tasks) {
	fmt.Println("")
	fmt.Printf("Open tasks (%d):\n", len(tasks.Open))
	for _, task := range tasks.Open {
		PrintTaskInLine(&task)
	}

	fmt.Println("")
	fmt.Printf("In progress tasks (%d):\n", len(tasks.InP))
	for _, task := range tasks.InP {
		PrintTaskInLine(&task)
	}

	fmt.Println("")
	fmt.Printf("Done tasks (%d):\n", len(tasks.Done))
	for _, task := range tasks.Done {
		PrintTaskInLine(&task)
	}
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
