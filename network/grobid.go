package network

import (
	"ronbun/storage"
	"strings"
	"time"
)

func GetGrobidResult(pdfData []byte) (string, error) {
	resp, err := client.Clone().SetTimeout(5*time.Minute).
		R().SetFileBytes("input", "a.pdf", pdfData).
		Post(getGrobidEndpoint() + "/api/processFulltextDocument")
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

var grobidEndpointChan = make(chan string)

func yieldingGrobidEndpoint() {
	i := 0
	for {
		grobidEndpointChan <- storage.Config.GrobidURLs[i%len(storage.Config.GrobidURLs)]
		i++
	}
}
func getGrobidEndpoint() string {
	return strings.TrimSuffix(<-grobidEndpointChan, "/")
}
