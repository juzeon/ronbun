package crawler

import (
	"github.com/imroc/req/v3"
	"time"
)

var client *req.Client

func init() {
	client = req.NewClient().SetTimeout(15 * time.Second).SetCommonRetryCount(3).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.IsErrorState()
		})
}
