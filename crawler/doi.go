package crawler

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func GetAbstract(doiLink string) (sourceHost string, abstract string, err error) {
	slog.Info("Requesting paper", "doi", doiLink)
	doiURL := lo.Must(url.Parse(doiLink))
	if doiURL.Hostname() == "dl.acm.org" {
		return "dl.acm.org", "", nil // TODO
	}
	var lastReq *http.Request
	client := clientPool.Get().(*req.Client)
	defer clientPool.Put(client)
	client.SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
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
	abstract, err = parseAbstract(sourceHost, resp.Body)
	if err != nil {
		return "", "", err
	}
	return sourceHost, abstract, nil
}

func parseAbstract(sourceHost string, reader io.Reader) (string, error) {
	abstractProvider, err := getAbstractProvider(sourceHost)
	if err != nil {
		return "", err
	}
	abstract, err := abstractProvider.ParseAbstract(reader)
	if err != nil {
		return "", err
	}
	abstract = stripTagsPolicy.Sanitize(abstract)
	abstract = regexp.MustCompile(`(?m)\s+`).ReplaceAllString(abstract, " ")
	return abstract, nil
}
func getAbstractProvider(sourceHost string) (AbstractProvider, error) {
	switch sourceHost {
	case "ieeexplore.ieee.org":
		return IEEEProvider{}, nil
	case "link.springer.com":
		return SpringerProvider{}, nil
	case "www.usenix.org":
		return USENIXProvider{}, nil
	default:
		return nil, errors.New("cannot find provider for " + sourceHost)
	}
}

var stripTagsPolicy = bluemonday.StripTagsPolicy()

type AbstractProvider interface {
	ParseAbstract(reader io.Reader) (string, error)
}
type IEEEProvider struct {
}

func (I IEEEProvider) ParseAbstract(reader io.Reader) (string, error) {
	html := string(lo.Must(io.ReadAll(reader)))
	re := regexp.MustCompile(`(?m)xplGlobal\.document\.metadata=(.*);$`)
	arr := re.FindStringSubmatch(html)
	if len(arr) == 0 {
		slog.Error("Cannot find ieee abstract", "html", html)
		return "", errors.New("cannot find ieee abstract")
	}
	res := gjson.Parse(arr[1])
	abstract := res.Get("abstract").String()
	if abstract == "" {
		slog.Error("IEEE Abstract is empty", "html", html)
		return "", errors.New("ieee abstract is empty")
	}
	return abstract, nil
}

type SpringerProvider struct {
}

func (s SpringerProvider) ParseAbstract(reader io.Reader) (string, error) {
	html := string(lo.Must(io.ReadAll(reader)))
	re := regexp.MustCompile(`<script type="application/ld\+json">(.*)</script>`)
	arr := re.FindStringSubmatch(html)
	if len(arr) == 0 {
		slog.Error("Cannot find springer abstract", "html", html)
		return "", errors.New("cannot find springer abstract")
	}
	res := gjson.Parse(arr[1])
	abstract := res.Get("description").String()
	if abstract == "" {
		slog.Error("Springer Abstract is empty", "html", html)
		return "", errors.New("springer abstract is empty")
	}
	return abstract, nil
}

type USENIXProvider struct {
}

func (U USENIXProvider) ParseAbstract(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}
	abstract := ""
	doc.Find("div.block-content div.content div.field").Each(func(i int, field *goquery.Selection) {
		if abstract != "" {
			return
		}
		if !strings.Contains(field.Find("div.field-label").Text(), "Abstract:") {
			return
		}
		abstract = field.Find("div.field-items").Text()
	})
	if abstract == "" {
		slog.Error("USENIX Abstract is empty", "html", lo.Must(doc.Html()))
		return "", errors.New("usenix abstract is empty")
	}
	return abstract, nil
}
