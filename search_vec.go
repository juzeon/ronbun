package main

import (
	"bytes"
	"github.com/samber/lo"
	"html/template"
	"log/slog"
	"math"
	"os"
	"ronbun/crawler"
	"ronbun/db"
	"ronbun/storage"
	"ronbun/util"
	"slices"
	"strings"
	"time"
)

func SearchVec() {
	file := storage.WriteTmpFile("search-vec-input.txt", nil)
	util.OpenFileWithDefaultProgram(file)
	util.PromptConfirmation("Please input the content for search, save the file and click to confirm.")
	v := lo.Must(os.ReadFile(file))
	slog.Info("Getting the embedding of the query...")
	queryDocText := string(v)
	arr := lo.Must(crawler.GetJinaEmbedding([]string{queryDocText}))
	query := arr[0]
	slog.Info("Fetching papers from the database...")
	papers := db.PaperTx.Select("id,embedding").MustFindMany("embedding!=?", "")
	type PaperWithDistance struct {
		ID       int
		Distance float64
	}
	slog.Info("Calculating cosine similarity...")
	res := util.ComputeCosine(query, lo.Map(papers, func(paper db.Paper, index int) []float64 {
		return paper.Embedding
	}))
	slog.Info("Sorting...")
	var papersWithDistance []PaperWithDistance
	for i, dis := range res {
		papersWithDistance = append(papersWithDistance, PaperWithDistance{
			ID:       papers[i].ID,
			Distance: dis,
		})
	}
	slices.SortFunc(papersWithDistance, func(a, b PaperWithDistance) int {
		if a.Distance > b.Distance {
			return -1
		}
		if a.Distance < b.Distance {
			return 1
		}
		return 0
	})
	slog.Info("Generating result...")
	ceiling := int(math.Min(float64(20), float64(len(papersWithDistance))))
	papers = db.PaperTx.MustFindMany("id in ?",
		lo.Map(papersWithDistance, func(paper PaperWithDistance, index int) int {
			return paper.ID
		})[0:ceiling])
	papersMap := lo.SliceToMap(papers, func(paper db.Paper) (int, db.Paper) {
		return paper.ID, paper
	})
	papers = make([]db.Paper, 0)
	for i, paperWithDistance := range papersWithDistance {
		if i >= ceiling {
			break
		}
		papers = append(papers, papersMap[paperWithDistance.ID])
	}
	tmpl := lo.Must(template.New("search_result").Funcs(map[string]any{
		"upper":           strings.ToUpper,
		"formatShortYear": util.FormatShortYear,
	}).Parse(searchResultTmpl))
	out := &bytes.Buffer{}
	lo.Must0(tmpl.Execute(out, SearchResultTmplData{
		SearchDoc: queryDocText,
		Papers:    papers,
	}))
	util.OpenFileWithDefaultProgram(storage.WriteTmpFile("Search by document "+
		time.Now().Format(time.DateTime)+".html", out.Bytes()))
}