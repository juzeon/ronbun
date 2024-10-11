package network

import (
	"ronbun/storage"
	"strings"
)

func GetGrobidResult(pdfData []byte) (GrobidTEI, error) {
	var empty GrobidTEI
	resp, err := client.Clone().SetTimeout(0).
		R().SetFileBytes("input", "a.pdf", pdfData).
		Post(getGrobidEndpoint() + "/api/processFulltextDocument")
	if err != nil {
		return empty, err
	}
	response, err := NewGrobidTEI(resp.String())
	if err != nil {
		return empty, err
	}
	return response, nil
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
