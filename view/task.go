package view

import (
	"bcli/api"
	"fmt"
	"github.com/charmbracelet/glamour"
)

func PrintTaskInLine(task *api.Task) {
	var seqNo = task.SeqNo
	var taskNm = task.TaskNm

	fmt.Printf(" - #%d %s [%s]\n", seqNo, taskNm, task.ReqID)
}

func RenderMarkdown(md string) {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	out, _ := r.Render(md)
	fmt.Println(out)
}
