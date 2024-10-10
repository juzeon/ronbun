package crawler

import (
	"github.com/samber/lo"
	"log/slog"
	"net/http"
	"net/url"
)

func GetAbstract(doiLink string) (sourceHost string, abstract string, err error) {
	slog.Info("Requesting paper", "doi", doiLink)
	doiURL := lo.Must(url.Parse(doiLink))
	if doiURL.Hostname() == "dl.acm.org" {
		return "dl.acm.org", "", nil // TODO
	}
	var lastReq *http.Request
	client := client.Clone().SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
		lastReq = req
		if req.URL.Hostname() == "dl.acm.org" {
			return http.ErrUseLastResponse
		}
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
	if sourceHost == "dl.acm.org" {
		return
	}
	
	return sourceHost, abstract, nil
}
