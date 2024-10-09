package ccf

import (
	"bytes"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"ronbun/storage"
	"strings"
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
