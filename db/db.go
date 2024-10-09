package db

import (
	"github.com/glebarez/sqlite"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"ronbun/storage"
)

var DB *gorm.DB
var SettingTx *TxWrapper[Setting]
var PaperTx *TxWrapper[Paper]

func init() {
	db := lo.Must(gorm.Open(sqlite.Open(storage.DatabasePath), &gorm.Config{}))
	DB = db

	lo.Must0(db.AutoMigrate(&Setting{}))
	SettingTx = NewTxWrapper[Setting](db)

	lo.Must0(db.AutoMigrate(&Paper{}))
	PaperTx = NewTxWrapper[Paper](db)
}
