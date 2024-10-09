package util

import (
	"github.com/samber/lo"
	"log/slog"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
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
func OpenFileWithDefaultProgram(filePath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		panic("unsupported platform")
	}
	lo.Must0(cmd.Start())
}
func FormatShortYear(year int) string {
	return strings.TrimPrefix(strconv.Itoa(year), "20")
}
