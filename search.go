package main

import (
	"bytes"
	"github.com/samber/lo"
	"ronbun/ccf"
	"ronbun/db"
	"ronbun/storage"
	"ronbun/util"
	"strings"
)

func Search() {
	keyword := util.PromptInputSearchKeyword()
	conferenceSubs := util.PromptSelectConferenceSubs()
	conferenceRankings := util.PromptSelectConferenceRankings()
	conferences := ccf.GetConferencesBySubRanking(conferenceSubs, conferenceRankings)
	papers := db.PaperTx.Order("year desc,conference desc").
		Where("conference in ?", lo.Map(conferences, func(c ccf.Conference, index int) string {
			return c.DBLP
		})).MustFindMany("title like ?", "%"+keyword+"%")
	var out bytes.Buffer
	out.WriteString("# Search for `" + keyword + "`\n\n")
	for _, paper := range papers {
		out.WriteString("\\[" + strings.ToUpper(paper.Conference) + " '" + util.FormatShortYear(paper.Year) + "] " +
			paper.Title + " \\[[dblp](" + paper.DBLPLink + ")] \\[[doi](" + paper.DOILink + ")]\n\n")
	}
	util.OpenFileWithDefaultProgram(storage.WriteTmpFile("Search for "+keyword+".md", out.Bytes()))
}
