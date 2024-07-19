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

func PrintTaskList(tasks []api.Task) {
	groupedTasks := groupTaskByStatusName(tasks)

	for status, tasks := range groupedTasks {
		fmt.Printf("\n📂 %s%s%s (%d)\n", ansiBold, status, ansiReset, len(tasks))
		for _, task := range tasks {
			PrintTaskInLine(&task)
		}
	}
}

func PrintTaskListInMarkdown(tasks []api.Task) {
	var markdown = "# Task list\n"

	groupedTasks := groupTaskByStatusName(tasks)
    for status, tasks := range groupedTasks {
        markdown += fmt.Sprintf("\n## %s\n", status)
        for _, task := range tasks {
            markdown += fmt.Sprintf("- #%d %s \n  - [](%s)\n", task.SeqNo, task.TaskNm, taskDetailPrefix+task.ReqID)
        }
    }

	RenderMarkdown(markdown)
}

func PrintSimpleTaskList(tasks []api.Task) {
	for _, task := range tasks {
		PrintTaskInLine(&task)
	}
}

func groupTaskByStatusName(tasks []api.Task) map[string][]api.Task {
	var taskMap = make(map[string][]api.Task)
	for _, task := range tasks {
		taskMap[task.ReqStsNm] = append(taskMap[task.ReqStsNm], task)
	}

	return taskMap
}
