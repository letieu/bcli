package paser

import (
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
