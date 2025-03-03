package network

import (
	"github.com/samber/lo"
	"log/slog"
	"ronbun/storage"
	"sync"
	"time"
)

type EmbeddingProvider interface {
	GetEmbedding(documents []string) ([][]float64, error)
}
type SiliconFlowEmbeddingProvider struct {
}

func (o SiliconFlowEmbeddingProvider) GetEmbedding(documents []string) ([][]float64, error) {
	var response SiliconFlowResponse
	slog.Info("Firing siliconflow request", "token", storage.Config.SiliconFlowToken)
	_, err := client.Clone().SetTimeout(5*time.Minute).R().SetHeader("Authorization", "Bearer "+
		storage.Config.SiliconFlowToken).
		SetSuccessResult(&response).SetBody(SiliconFlowRequest{
		Model:          "BAAI/bge-m3",
		Input:          documents,
		EncodingFormat: "float",
	}).Post("https://api.siliconflow.cn/v1/embeddings")
	if err != nil {
		return nil, err
	}
	return lo.Map(response.Data, func(item SiliconFlowData, index int) []float64 {
		return item.Embedding
	}), nil
}

type JinaEmbeddingProvider struct {
}

func (o JinaEmbeddingProvider) GetEmbedding(documents []string) ([][]float64, error) {
	startYieldingJinaToken()
	var response JinaResponse
	token := getJinaToken()
	slog.Info("Firing jina request", "token", token)
	_, err := client.Clone().SetTimeout(5*time.Minute).R().SetHeader("Authorization", "Bearer "+token).
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
var startYieldingJinaToken = sync.OnceFunc(func() {
	go yieldingJinaToken()
})

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
