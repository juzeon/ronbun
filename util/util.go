package util

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/microcosm-cc/bluemonday"
	"github.com/samber/lo"
	"log/slog"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func AttemptMax[T any](max int, fun func() (T, error)) T {
	var obj T
	_, err := lo.Attempt(max, func(index int) error {
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
func Attempt[T any](fun func() (T, error)) T {
	return AttemptMax(10, fun)
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
func NormalizeConferenceSlug(slug string) string {
	return regexp.MustCompile(`[\d\-_]`).ReplaceAllString(slug, "")
}
func Sha1(v []byte) string {
	s := sha1.New()
	s.Write(v)
	return hex.EncodeToString(s.Sum(nil))
}

var stripTagsPolicy = bluemonday.StripTagsPolicy()

func StripHTMLTags(html string) string {
	return stripTagsPolicy.Sanitize(html)
}
