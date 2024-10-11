package util

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"strings"
)

type GrobidData struct {
	Title    string
	Abstract string
	Sections []GrobidDataSection
}
type GrobidDataSection struct {
	Title   string
	Content string
}

func ParseGrobidXML(xml string) GrobidData {
	xml = strings.ReplaceAll(xml, "<head", "<my-head")
	xml = strings.ReplaceAll(xml, "</head>", "</my-head>")
	xml = strings.ReplaceAll(xml, "<body", "<my-body")
	xml = strings.ReplaceAll(xml, "</body>", "</my-body>")
	doc := lo.Must(goquery.NewDocumentFromReader(strings.NewReader(xml)))
	title := doc.Find("titleStmt title").Text()
	abstract := strings.TrimSpace(doc.Find("profileDesc abstract").Text())
	var sections []GrobidDataSection
	doc.Find("text my-body div").Each(func(i int, div *goquery.Selection) {
		head := div.Find("my-head")
		headIndex := head.AttrOr("n", "")
		var content bytes.Buffer
		div.Contents().Each(func(i int, child *goquery.Selection) {
			if child.IsSelection(head) {
				return
			}
			content.WriteString(strings.TrimSpace(child.Text()) + "\n\n")
		})
		sections = append(sections, GrobidDataSection{
			Title:   lo.Ternary(headIndex == "", "", headIndex+" ") + head.Text(),
			Content: strings.TrimSpace(content.String()),
		})
	})
	return GrobidData{
		Title:    title,
		Abstract: abstract,
		Sections: sections,
	}
}
