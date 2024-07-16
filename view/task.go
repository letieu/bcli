package view

import (
	"bcli/api"
	"fmt"
)

func PrintTask(task *api.Task) {
	var seqNo = task.SeqNo
	var taskNm = task.TaskNm

	fmt.Printf(" - #%d %s [%s]\n", seqNo, taskNm, task.ReqID)
}
