package storage

import (
	"fmt"
	"github.com/flytam/filenamify"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
	"path/filepath"
)

type config struct {
	CCFPath           string       `yaml:"ccf_path"`
	EmbeddingProvider string       `yaml:"embedding_provider"`
	JinaTokens        []string     `yaml:"jina_tokens"`
	SiliconFlowToken  string       `yaml:"siliconflow_token"`
	Concurrency       int          `yaml:"concurrency"`
	SearchLimit       int          `yaml:"search_limit"`
	GrobidURLs        []string     `yaml:"grobid_urls"`
	OpenAI            openaiConfig `yaml:"openai"`
}
type openaiConfig struct {
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
	Key      string `yaml:"key"`
}

var Config config
var DatabasePath string
var TmpPath string

func init() {
	homePath := lo.Must(os.UserHomeDir())
	storagePath := filepath.Join(homePath, ".ronbun")
	configPath := filepath.Join(storagePath, "config.yml")
	DatabasePath = filepath.Join(storagePath, "ronbun.db")
	TmpPath = filepath.Join(storagePath, "tmp")
	if _, err := os.Stat(storagePath); err != nil {
		lo.Must0(os.MkdirAll(storagePath, 0755))
		lo.Must0(os.WriteFile(configPath, nil, 0644))
		fmt.Println("Please initialize: " + configPath)
		os.Exit(0)
	}
	if _, err := os.Stat(TmpPath); err != nil {
		lo.Must0(os.MkdirAll(TmpPath, 0755))
	}
	v := lo.Must(os.ReadFile(configPath))
	lo.Must0(yaml.Unmarshal(v, &Config))
	Config.JinaTokens = lo.Shuffle(Config.JinaTokens)
	Config.GrobidURLs = lo.Shuffle(Config.GrobidURLs)
	Config.OpenAI.Endpoint = "https://" + lo.Must(url.Parse(Config.OpenAI.Endpoint)).Host + "/v1"
}
func WriteTmpFile(filename string, v []byte) string {
	filename = lo.Must(filenamify.FilenamifyV2(filename))
	path := filepath.Join(TmpPath, filename)
	lo.Must0(os.WriteFile(path, v, 0644))
	return path
}
