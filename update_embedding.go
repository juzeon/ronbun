package main

import (
	"github.com/samber/lo"
	"log/slog"
	"math"
	"ronbun/db"
	"ronbun/network"
	"ronbun/util"
)

func UpdateEmbedding() {
	papers := db.PaperTx.Order("title asc").MustFindMany("abstract!=? and embedding=?", "", "")
	slog.Info("Papers to update", "count", len(papers))
	step := 100
	for i := 0; i < len(papers); i += step {
		ceiling := int(math.Min(float64(i+step), float64(len(papers))))
		slog.Info("Getting embeddings", "start", i, "end", ceiling)
		batch := papers[i:ceiling]
		res := util.AttemptMax(3, func() ([][]float64, error) {
			r, err := network.GetJinaEmbedding(lo.Map(batch, func(paper db.Paper, index int) string {
				return paper.Abstract
			}))
			if err != nil {
				return nil, err
			}
			return r, nil
		})
		for j := range batch {
			batch[j].Embedding = res[j]
			db.PaperTx.MustSave(&batch[j])
		}
		slog.Info("Successfully updating embeddings")
	}
}
