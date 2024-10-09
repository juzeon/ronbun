package db

import (
	"encoding/json"
	"github.com/samber/lo"
)

func SetSetting(key string, value string) {
	obj := SettingTx.MustFindOne("key=?", key)
	if obj == nil {
		SettingTx.MustCreate(&Setting{
			Key:   key,
			Value: value,
		})
	} else {
		SettingTx.Where("key=?", key).MustUpdate("value", value)
	}
}
func SetSettingObj[T any](key string, value T) {
	v := lo.Must(json.Marshal(&value))
	SetSetting(key, string(v))
}
func GetSetting(key string) string {
	obj := SettingTx.MustFindOne("key=?", key)
	if obj == nil {
		return ""
	}
	return obj.Value
}
func GetSettingObj[T any](key string) T {
	var obj T
	value := GetSetting(key)
	if value == "" {
		return obj
	}
	lo.Must0(json.Unmarshal([]byte(value), &obj))
	return obj
}
