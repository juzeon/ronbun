package network

type ConferenceInstance struct {
	Slug    string
	Title   string
	Year    int
	TocLink string
}
type Paper struct {
	Title              string
	DBLPLink           string
	DOILink            string
	ConferenceInstance ConferenceInstance
}

type JinaRequest struct {
	Model         string   `json:"model"`
	Task          string   `json:"task"`
	Dimensions    int      `json:"dimensions"`
	LateChunking  bool     `json:"late_chunking"`
	EmbeddingType string   `json:"embedding_type"`
	Input         []string `json:"input"`
}
type JinaResponse struct {
	Model  string     `json:"model"`
	Object string     `json:"object"`
	Usage  JinaUsage  `json:"usage"`
	Data   []JinaData `json:"data"`
}
type JinaUsage struct {
	TotalTokens  int `json:"total_tokens"`
	PromptTokens int `json:"prompt_tokens"`
}
type JinaData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}
