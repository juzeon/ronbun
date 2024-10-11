package crawler

import (
	"errors"
	"github.com/imroc/req/v3"
	"ronbun/storage"
	"sync"
	"time"
)

var client *req.Client
var clientPool = &sync.Pool{}

func init() {
	client = req.NewClient().ImpersonateChrome().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.IsErrorState() {
				return errors.New("http status: " + resp.GetStatus() + " (response: " + resp.String() + ")")
			}
			return nil
		}).SetTimeout(15 * time.Second).SetCommonRetryCount(5).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.IsErrorState()
		})
	for range storage.Config.Concurrency {
		clientPool.Put(client.Clone())
	}
	clientPool.New = func() any {
		return client.Clone()
	}
	go yieldingDBLPDomain()
	go yieldingJinaToken()
}
