package network

import (
	"github.com/samber/lo"
	"log/slog"
	"ronbun/storage"
)

func GetJinaEmbedding(documents []string) ([][]float64, error) {
	var response JinaResponse
	token := getJinaToken()
	slog.Info("Firing jina request", "token", token)
	_, err := client.Clone().SetTimeout(0).R().SetHeader("Authorization", "Bearer "+token).
		SetSuccessResult(&response).SetBody(JinaRequest{
		Model:         "jina-embeddings-v3",
		Task:          "text-matching",
		Dimensions:    1024,
		LateChunking:  false,
		EmbeddingType: "float",
		Input:         documents,
	}).Post("https://api.jina.ai/v1/embeddings")
	if err != nil {
		return nil, err
	}
	return lo.Map(response.Data, func(item JinaData, index int) []float64 {
		return item.Embedding
	}), nil
}

var jinaTokenChan = make(chan string)

func yieldingJinaToken() {
	i := 0
	for {
		jinaTokenChan <- storage.Config.JinaTokens[i%len(storage.Config.JinaTokens)]
		i++
	}
}
func getJinaToken() string {
	return <-jinaTokenChan
}
