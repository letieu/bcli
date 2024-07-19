package paser

import (
	"fmt"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/russross/blackfriday/v2"
)

var converter = md.NewConverter("", true, nil)

func HtmlToMd(html string) (string, error) {
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func MdToHtml(md string) (string, error) {
	html := blackfriday.Run([]byte(md))
	return string(html), nil
}

func CreateBufText(title string, content string) (string, error) {
	contentMd, err := HtmlToMd(content)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("# %s\n\n%s", title, contentMd), nil
}

func ParseBufText(bufText string) (string, string, error) {
	lines := strings.Split(bufText, "\n")
	title := strings.TrimPrefix(lines[0], "# ")
	content := strings.Join(lines[2:], "\n")
	contentHtml, err := MdToHtml(content)
	if err != nil {
		return "", "", err
	}

	return title, contentHtml, nil
}

func GetNewTaskPayload(templateString string, data map[string]string) string {
	for key, value := range data {
		templateString = strings.ReplaceAll(templateString, fmt.Sprintf("{{%s}}", key), value)
	}
	return templateString
}

func GetGitBranchName(taskNo int, taskTitle string) string {
    name := fmt.Sprintf("%d-%s", taskNo, strings.ToLower(strings.ReplaceAll(taskTitle, " ", "-")))

    re := regexp.MustCompile(`\[[^\]]*\]`)
    name = re.ReplaceAllString(name, "")

    // replace all non-alphanumeric characters with a dash
    re = regexp.MustCompile(`[^a-zA-Z0-9-]`)
    name = re.ReplaceAllString(name, "-")

    return name
}
