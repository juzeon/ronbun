package db

type Setting struct {
	ID    int
	Key   string `gorm:"not null;uniqueIndex"`
	Value string `gorm:"not null"`
}

type Paper struct {
	ID         int
	Title      string    `gorm:"not null"`
	Conference string    `gorm:"not null"`
	Year       int       `gorm:"not null"`
	DBLPLink   string    `gorm:"column:dblp_link;not null"`
	DOILink    string    `gorm:"column:doi_link;not null;uniqueIndex"`
	SourceHost string    `gorm:"not null;default:''"`
	Abstract   string    `gorm:"not null;default:''"`
	Embedding  []float64 `gorm:"not null;serializer:json;default:''"`
}

type Translation struct {
	ID          int
	Hash        string   `gorm:"uniqueIndex;not null;default:''"`
	Title       string   `gorm:"not null;default:''"`
	GrobidData  string   `gorm:"not null;default:''"`
	ChineseData []string `gorm:"not null;default:'';serializer:json"`
}
