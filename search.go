package main

import (
	"bytes"
	_ "embed"
	"github.com/samber/lo"
	"html/template"
	"ronbun/ccf"
	"ronbun/db"
	"ronbun/storage"
	"ronbun/util"
	"strings"
)

//go:embed asset/search_result.html
var searchResultTmpl string

type SearchResultTmplData struct {
	SearchDoc string
	Keyword   string
	Papers    []db.Paper
}

func Search() {
	keyword := util.PromptInputSearchKeyword()
	conferenceSubs := util.PromptSelectConferenceSubs()
	conferenceRankings := util.PromptSelectConferenceRankings()
	conferences := ccf.GetConferencesBySubRanking(conferenceSubs, conferenceRankings)
	papers := db.PaperTx.Order("year desc,conference desc").
		Where("conference in ?", lo.Map(conferences, func(c ccf.Conference, index int) string {
			return c.DBLP
		})).MustFindMany("title like ?", "%"+keyword+"%")
	tmpl := lo.Must(template.New("search_result").Funcs(map[string]any{
		"upper":           strings.ToUpper,
		"formatShortYear": util.FormatShortYear,
	}).Parse(searchResultTmpl))
	out := &bytes.Buffer{}
	lo.Must0(tmpl.Execute(out, SearchResultTmplData{
		Keyword: keyword,
		Papers:  papers,
	}))
	util.OpenFileWithDefaultProgram(storage.WriteTmpFile("Search for "+keyword+".html", out.Bytes()))
}
