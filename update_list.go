package main

import (
	"github.com/samber/lo"
	"log/slog"
	"ronbun/ccf"
	"ronbun/crawler"
	"ronbun/db"
	"ronbun/util"
	"slices"
	"sync"
)

func UpdateList() {
	conferenceSubs := util.PromptSelectConferenceSubs()
	conferenceRankings := util.PromptSelectConferenceRankings()
	startYear := util.PromptInputStartYear()
	var conferences []ccf.Conference
	for _, sub := range conferenceSubs {
		conferences = append(conferences, ccf.GetConferencesBySub(sub.Sub)...)
	}
	conferences = lo.Filter(conferences, func(c ccf.Conference, index int) bool {
		return slices.Contains(conferenceRankings, c.Rank.CCF)
	})
	getPapersFromDBLP(lo.Map(conferences, func(conference ccf.Conference, index int) string {
		return conference.DBLP
	}), startYear)
}
func getPapersFromDBLP(slugs []string, startYear int) {
	wg := &sync.WaitGroup{}
	slugChan := make(chan string)
	wg.Add(crawler.MaxThread)
	for range crawler.MaxThread {
		go func() {
			for slug := range slugChan {
				insArr, err := crawler.GetConferenceInstancesBySlug(slug)
				if err != nil {
					slog.Error("Failed to get conference instance", "slug", slug, "err", err)
					continue
				}
				insArr = lo.Filter(insArr, func(ins crawler.ConferenceInstance, index int) bool {
					return ins.Year >= startYear
				})
				for _, ins := range insArr {
					papers, err := crawler.GetPapersByConferenceInstance(ins)
					if err != nil {
						slog.Error("Failed to get papers", "slug", slug, "year", ins.Year, "err", err)
						continue
					}
					duplicatePapers := db.PaperTx.MustFindMany("doi_link in ?",
						lo.Map(papers, func(paper crawler.Paper, index int) string {
							return paper.DOILink
						}))
					duplicateDOILinkMap := lo.SliceToMap(duplicatePapers, func(paper db.Paper) (string, struct{}) {
						return paper.DOILink, struct{}{}
					})
					papers = lo.Filter(papers, func(paper crawler.Paper, index int) bool {
						_, exist := duplicateDOILinkMap[paper.DOILink]
						return !exist
					})
					slog.Info("Collected deduplicate papers",
						"count", len(papers), "slug", slug, "year", ins.Year)
					db.PaperTx.MustCreateMany(lo.Map(papers, func(paper crawler.Paper, index int) db.Paper {
						return db.Paper{
							Title:      paper.Title,
							Conference: paper.ConferenceInstance.Slug,
							Year:       paper.ConferenceInstance.Year,
							DBLPLink:   paper.DBLPLink,
							DOILink:    paper.DOILink,
						}
					}))
				}
			}
			wg.Done()
		}()
	}
	for _, slug := range slugs {
		slugChan <- slug
	}
	close(slugChan)
	wg.Wait()
}
