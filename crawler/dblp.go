package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"log/slog"
	"strconv"
	"strings"
)

func GetPapersByConferenceInstance(ins ConferenceInstance) ([]Paper, error) {
	slog.Info("Requesting papers of conference instance",
		"slug", ins.Slug, "year", ins.Year, "url", ins.TocLink)
	resp, err := client.R().Get(ins.TocLink)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var papers []Paper
	doc.Find("ul.publ-list li.inproceedings").Each(func(i int, paperSelection *goquery.Selection) {
		title := strings.TrimSuffix(paperSelection.Find("cite span.title[itemprop=name]").Text(), ".")
		if title == "" {
			slog.Warn("Title of paper is empty", "paper", lo.Must(paperSelection.Html()))
			return
		}
		button := paperSelection.Find("nav.publ li").First()
		doiLink := button.Find("div.head a").AttrOr("href", "")
		if doiLink == "" {
			slog.Warn("DOILink of paper is empty", "paper", lo.Must(paperSelection.Html()))
			return
		}
		dblpLink := button.Find("div.body ul").First().Find("li.details a").
			AttrOr("href", "")
		if dblpLink == "" {
			slog.Warn("DBLPLink of paper is empty", "paper", lo.Must(paperSelection.Html()))
			return
		}
		papers = append(papers, Paper{
			Title:              title,
			DBLPLink:           dblpLink,
			DOILink:            doiLink,
			ConferenceInstance: ins,
		})
	})
	return papers, nil
}
func GetConferenceInstancesBySlug(slug string) ([]ConferenceInstance, error) {
	url := "https://dblp.org/db/conf/" + slug
	slog.Info("Requesting conference", "url", url)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var conferenceInstances []ConferenceInstance
	doc.Find("ul.publ-list").Each(func(i int, pubList *goquery.Selection) {
		ins := pubList.Find("li cite").First()
		title := strings.TrimSuffix(ins.Find("span.title[itemprop=name]").Text(), ".")
		if title == "" {
			slog.Warn("Title of conference instance is empty", "ins", lo.Must(ins.Html()))
			return
		}
		year, err := strconv.Atoi(ins.Find("span[itemprop=datePublished]").Text())
		if err != nil {
			slog.Warn("Error parsing datePublished", "ins", lo.Must(ins.Html()))
			return
		}
		tocLink := ins.Find("a.toc-link").AttrOr("href", "")
		if tocLink == "" {
			slog.Warn("TocLink of conference instance is empty", "ins", lo.Must(ins.Html()))
			return
		}
		conferenceInstances = append(conferenceInstances, ConferenceInstance{
			Slug:    slug,
			Title:   title,
			Year:    year,
			TocLink: tocLink,
		})
	})
	return conferenceInstances, nil
}