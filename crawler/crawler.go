package crawler

import (
	"errors"
	"github.com/imroc/req/v3"
	"time"
)

const MaxThread = 20

var client *req.Client

func init() {
	client = req.NewClient().ImpersonateFirefox().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.IsErrorState() {
				return errors.New("http status: " + resp.GetStatus() + " (response: " + resp.String() + ")")
			}
			return nil
		}).SetTimeout(15 * time.Second).SetCommonRetryCount(5).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.IsErrorState()
		})
	go yieldingDBLPDomain()
}
