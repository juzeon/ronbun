package ccf

import (
	"bytes"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"ronbun/storage"
	"slices"
	"strings"
	"sync"
)

type ConferenceSub struct {
	Name string `yaml:"name"`
	Sub  string `yaml:"sub"`
}

func GetConferenceSubs() []ConferenceSub {
	yamlFile := filepath.Join(storage.Config.CCFPath, "conference", "types.yml")
	v := lo.Must(os.ReadFile(yamlFile))
	var conferenceSubs []ConferenceSub
	lo.Must0(yaml.Unmarshal(v, &conferenceSubs))
	return conferenceSubs
}
func GetConferencesBySub(sub string) []Conference {
	subDir := filepath.Join(storage.Config.CCFPath, "conference", sub)
	fileList := lo.Must(os.ReadDir(subDir))
	var str bytes.Buffer
	for _, conferenceFile := range fileList {
		if !strings.HasSuffix(conferenceFile.Name(), ".yml") {
			continue
		}
		str.Write(lo.Must(os.ReadFile(filepath.Join(subDir, conferenceFile.Name()))))
		str.WriteString("\n")
	}
	var conferences []Conference
	lo.Must0(yaml.Unmarshal(str.Bytes(), &conferences))
	return conferences
}
func GetConferencesBySubRanking(subs []ConferenceSub, rankings []string) []Conference {
	var conferences []Conference
	for _, sub := range subs {
		conferences = append(conferences, GetConferencesBySub(sub.Sub)...)
	}
	conferences = lo.Filter(conferences, func(c Conference, index int) bool {
		return slices.Contains(rankings, c.Rank.CCF)
	})
	return conferences
}

var conferenceSlugMap = map[string]Conference{}
var initConferenceSlugMapOnce = sync.OnceFunc(func() {
	for _, sub := range GetConferenceSubs() {
		for _, conference := range GetConferencesBySub(sub.Sub) {
			conference := conference
			conferenceSlugMap[conference.DBLP] = conference
		}
	}
})

func GetConferenceBySlug(slug string) Conference {
	initConferenceSlugMapOnce()
	return conferenceSlugMap[slug]
}
func GetConferenceRankingBySlug(slug string) string {
	return GetConferenceBySlug(slug).Rank.CCF
}
