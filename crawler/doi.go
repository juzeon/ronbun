package crawler

import (
	"log/slog"
	"net/http"
)

func GetAbstract(doiLink string) (sourceHost string, abstract string, err error) {
	slog.Info("Requesting paper", "doi", doiLink)
	var lastReq *http.Request
	client := client.Clone().SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
		lastReq = req
		return nil
	})
	resp, err := client.R().Get(doiLink)
	if err != nil {
		return "", "", err
	}
	if lastReq == nil {
		lastReq = resp.Request.RawRequest
	}
	sourceHost = lastReq.URL.Hostname()
	return sourceHost, abstract, nil
}
