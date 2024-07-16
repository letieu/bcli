package view

import (
	"fmt"
	"github.com/charmbracelet/glamour"
)

func RenderMarkdown(md string) {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
	)
	out, _ := r.Render(md)
	fmt.Println(out)
}
