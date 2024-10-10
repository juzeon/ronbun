package main

import (
	"log/slog"
	"ronbun/crawler"
	"ronbun/db"
	"sync"
)

func UpdatePaper() {
	papers := db.PaperTx.Order("title asc").MustFindMany("source_host=? or abstract=? "+
		"or embedding=?", "", "", "")
	//papers := db.PaperTx.Order("title asc").MustFindMany("source_host=?", "")
	slog.Info("Paper waiting to update", "count", len(papers))
	wg := &sync.WaitGroup{}
	paperChan := make(chan *db.Paper)
	wg.Add(crawler.MaxThread)
	for range crawler.MaxThread {
		go func() {
			for paper := range paperChan {
				if paper.SourceHost == "" || paper.Abstract == "" {
					sourceHost, abstract, err := crawler.GetAbstract(paper.DOILink)
					if err != nil {
						slog.Error("Error getting abstract", "paper", paper, "err", err)
						continue
					}
					paper.SourceHost = sourceHost
					paper.Abstract = abstract
					db.PaperTx.MustSave(paper)
				}
				// TODO
			}
			wg.Done()
		}()
	}
	for _, paper := range papers {
		paper := paper
		paperChan <- &paper
	}
	close(paperChan)
	wg.Wait()
}
