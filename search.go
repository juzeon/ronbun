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

var searchTmplFuncs = template.FuncMap(map[string]any{
	"upper":           strings.ToUpper,
	"formatShortYear": util.FormatShortYear,
	"getRanking":      ccf.GetConferenceRankingBySlug,
})

func Search() {
	keyword := util.PromptInputSearchKeyword()
	conferenceRankings := util.PromptSelectConferenceRankings()
	conferences := ccf.GetConferencesBySubRanking(ccf.GetConferenceSubs(), conferenceRankings)
	papers := db.PaperTx.Order("year desc,conference desc").
		Where("conference in ?", lo.Map(conferences, func(c ccf.Conference, index int) string {
			return c.DBLP
		})).MustFindMany("title like ?", "%"+keyword+"%")
	tmpl := lo.Must(template.New("search_result").Funcs(searchTmplFuncs).Parse(searchResultTmpl))
	out := &bytes.Buffer{}
	lo.Must0(tmpl.Execute(out, SearchResultTmplData{
		Keyword: keyword,
		Papers:  papers,
	}))
	util.OpenFileWithDefaultProgram(storage.WriteTmpFile("Search for "+keyword+".html", out.Bytes()))
}
