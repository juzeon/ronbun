package util

import (
	"github.com/samber/lo"
	"log/slog"
)

func Attempt[T any](fun func() (T, error)) T {
	var obj T
	_, err := lo.Attempt(10, func(index int) error {
		t, err := fun()
		if err != nil {
			slog.Warn("Please retry", "err", err)
			return err
		}
		obj = t
		return nil
	})
	if err != nil {
		panic(err)
	}
	return obj
}
