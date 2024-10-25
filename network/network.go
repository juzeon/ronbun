package network

import (
	"errors"
	"github.com/imroc/req/v3"
	"github.com/sashabaranov/go-openai"
	"log/slog"
	"ronbun/storage"
	"sync"
	"time"
)

var client *req.Client
var clientPool = &sync.Pool{}
var openaiClient *openai.Client

func init() {
	client = req.NewClient().ImpersonateChrome().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.IsErrorState() {
				return errors.New("http status: " + resp.GetStatus() + " (response: " + resp.String() + ")")
			}
			return nil
		}).SetTimeout(15 * time.Second).SetCommonRetryCount(10).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			if resp.GetStatusCode() == 404 {
				return false
			}
			return resp.IsErrorState()
		}).
		AddCommonRetryHook(func(resp *req.Response, err error) {
			url := ""
			if resp.Request != nil {
				url = resp.Request.RawURL
			}
			slog.Warn("Retry triggered", "url", url, "status", resp.GetStatusCode(), "err", err)
		})
	for range storage.Config.Concurrency {
		clientPool.Put(client.Clone())
	}
	clientPool.New = func() any {
		return client.Clone()
	}
	openaiConfig := openai.DefaultConfig(storage.Config.OpenAI.Key)
	openaiConfig.BaseURL = storage.Config.OpenAI.Endpoint
	openaiClient = openai.NewClientWithConfig(openaiConfig)
	go yieldingDBLPDomain()
	go yieldingJinaToken()
	go yieldingGrobidEndpoint()
}
