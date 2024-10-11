package main

import (
	"github.com/samber/lo"
	"log/slog"
	"ronbun/ccf"
	"ronbun/crawler"
	"ronbun/db"
	"ronbun/storage"
	"ronbun/util"
	"sync"
)

func UpdateList() {
	conferenceSubs := util.PromptSelectConferenceSubs()
	conferenceRankings := util.PromptSelectConferenceRankings()
	startYear := util.PromptInputStartYear()
	conferences := ccf.GetConferencesBySubRanking(conferenceSubs, conferenceRankings)
	getPapersFromDBLP(lo.Map(conferences, func(conference ccf.Conference, index int) string {
		return conference.DBLP
	}), startYear)
}
func getPapersFromDBLP(slugs []string, startYear int) {
	wg := &sync.WaitGroup{}
	slugChan := make(chan string)
	wg.Add(storage.Config.Concurrency)
	for range storage.Config.Concurrency {
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
					duplicatePapers := db.PaperTx.MustFindMany("dblp_link in ?",
						lo.Map(papers, func(paper crawler.Paper, index int) string {
							return paper.DBLPLink
						}))
					duplicateDBLPLinkMap := lo.SliceToMap(duplicatePapers, func(paper db.Paper) (string, struct{}) {
						return paper.DBLPLink, struct{}{}
					})
					papers = lo.Filter(papers, func(paper crawler.Paper, index int) bool {
						_, exist := duplicateDBLPLinkMap[paper.DBLPLink]
						return !exist
					})
					slog.Info("Collected deduplicate papers",
						"count", len(papers), "slug", slug, "year", ins.Year)
					arr := lo.Map(papers, func(paper crawler.Paper, index int) db.Paper {
						return db.Paper{
							Title:      paper.Title,
							Conference: paper.ConferenceInstance.Slug,
							Year:       paper.ConferenceInstance.Year,
							DBLPLink:   paper.DBLPLink,
							DOILink:    paper.DOILink,
						}
					})
					if len(arr) != 0 {
						db.PaperTx.MustCreateMany(arr)
					}
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
