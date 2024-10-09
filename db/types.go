package db

type Setting struct {
	ID    int
	Key   string `gorm:"not null;uniqueIndex"`
	Value string `gorm:"not null"`
}
