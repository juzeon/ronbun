package network

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
	"regexp"
	"ronbun/util"
	"strings"
)

func GetAbstract(doiLink string) (sourceHost string, abstract string, err error) {
	slog.Info("Requesting paper", "doi", doiLink)
	var lastReq *http.Request
	client := clientPool.Get().(*req.Client)
	defer clientPool.Put(client)
	client.SetRedirectPolicy(func(req *http.Request, via []*http.Request) error {
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
	abstract, err = parseAbstract(sourceHost, resp.Body)
	if err != nil {
		return "", "", err
	}
	slog.Info("Successfully collected paper", "doi", doiLink)
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
	abstract = util.StripHTMLTags(abstract)
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
	case "dl.acm.org":
		return ConfigurableAbstractProvider{Selector: "div#abstracts div[role=paragraph]"}, nil
	case "aclanthology.org":
		return ConfigurableAbstractProvider{Selector: ".acl-abstract span"}, nil
	case "proceedings.mlr.press":
		return ConfigurableAbstractProvider{Selector: "div#abstract"}, nil
	case "ojs.aaai.org":
		return ConfigurableAbstractProvider{Selector: "section.abstract"}, nil
	case "www.ijcai.org":
		return ConfigurableAbstractProvider{Selector: "div.col-md-12:first-of-type"}, nil
	case "ebooks.iospress.nl":
		return ConfigurableAbstractProvider{Selector: "div.abstract section"}, nil
	case "proceedings.neurips.cc":
		return ConfigurableAbstractProvider{
			Regexp:      `(?m)<h4>Abstract</h4>([\s\S]*?)</div>`,
			RegexpGroup: 1,
		}, nil
	case "openaccess.thecvf.com":
		return ConfigurableAbstractProvider{Selector: "div#abstract"}, nil
	case "openreview.net":
		return ConfigurableAbstractProvider{Selector: "div.note-content .note-content-value"}, nil
	case "proceedings.kr.org":
		return ConfigurableAbstractProvider{
			Regexp:      `(?m)<h2>Abstract</h2>([\s\S]*)`,
			RegexpGroup: 1,
		}, nil
	default:
		return nil, errors.New("cannot find provider for " + sourceHost)
	}
}

type AbstractProvider interface {
	ParseAbstract(reader io.Reader) (string, error)
}
type ConfigurableAbstractProvider struct {
	Regexp      string
	RegexpGroup int

	Selector string
}

func (o ConfigurableAbstractProvider) ParseAbstract(reader io.Reader) (string, error) {
	if o.Selector != "" {
		return o.parseBySelector(reader)
	}
	if o.Regexp != "" && o.RegexpGroup != 0 {
		return o.parseByRegexp(reader)
	}
	return "", errors.New("misconfigured abstract provider")
}
func (o ConfigurableAbstractProvider) parseBySelector(reader io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(doc.Find(o.Selector).Text())
	if text == "" {
		slog.Error("Abstract is empty: " + o.Selector)
		return "", errors.New("abstract is empty")
	}
	return text, nil
}
func (o ConfigurableAbstractProvider) parseByRegexp(reader io.Reader) (string, error) {
	html, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	arr := regexp.MustCompile(o.Regexp).FindStringSubmatch(string(html))
	if len(arr) <= o.RegexpGroup {
		return "", errors.New("regexp group mismatched: " + o.Regexp)
	}
	text := strings.TrimSpace(bluemonday.StripTagsPolicy().Sanitize(arr[o.RegexpGroup]))
	if text == "" {
		slog.Error("Abstract is empty: " + o.Regexp)
		return "", errors.New("abstract is empty")
	}
	return text, nil
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
		slog.Error("IEEE Abstract is empty")
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
		slog.Error("Springer Abstract is empty")
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
		slog.Error("USENIX Abstract is empty")
		return "", errors.New("usenix abstract is empty")
	}
	return abstract, nil
}
