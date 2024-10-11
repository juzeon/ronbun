package main

import (
	"bytes"
	_ "embed"
	"github.com/ncruces/zenity"
	"github.com/samber/lo"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"ronbun/db"
	"ronbun/network"
	"ronbun/storage"
	"ronbun/util"
	"slices"
	"strings"
	"sync"
)

//go:embed asset/translation.html
var translationTmpl string
var translationFuncsMap = template.FuncMap(map[string]any{})

type TranslationTmplData struct {
	Title    string
	Abstract string
	Section  []TranslationTmplDataSection
}
type TranslationTmplDataSection struct {
	Title   string
	Content template.HTML
}

func Translate() {
	pdfFile := util.PromptOpenPDFFile()
	v := lo.Must(os.ReadFile(pdfFile))
	hash := util.Sha1(v)
	translation := db.TranslationTx.MustFindOne("hash=?", hash)
	if translation == nil {
		slog.Info("Getting Grobid result...", "file", pdfFile)
		tei := lo.Must(network.GetGrobidResult(v))
		translation = &db.Translation{
			Hash:        hash,
			Title:       tei.Header.FileDesc.TitleStmt.Title.Value,
			GrobidData:  tei.String(),
			ChineseData: nil,
		}
		db.TranslationTx.MustCreate(translation)
	}
	if len(translation.ChineseData) == 0 {
		chineseData := translateChinese(lo.Must(network.NewGrobidTEI(translation.GrobidData)))
		translation.ChineseData = chineseData
		db.TranslationTx.MustSave(translation)
	}
	tmpl := lo.Must(template.New("translation").Funcs(translationFuncsMap).Parse(translationTmpl))
	out := &bytes.Buffer{}
	var tmplData TranslationTmplData
	for i, chinese := range translation.ChineseData {
		arr := strings.Split(chinese, "\n")
		arr = lo.Map(arr, func(line string, index int) string {
			return strings.TrimSpace(line)
		})
		arr = lo.Filter(arr, func(line string, index int) bool {
			return line != ""
		})
		if i == 0 {
			tmplData.Title = arr[0]
			tmplData.Abstract = arr[1]
			continue
		}
		tmplData.Section = append(tmplData.Section, TranslationTmplDataSection{
			Title:   arr[0],
			Content: template.HTML("<p>" + strings.Join(arr[1:], "</p><p>") + "</p>"),
		})
	}
	lo.Must0(tmpl.Execute(out, tmplData))
	p := storage.WriteTmpFile("Translation for "+translation.Title+".html", out.Bytes())
	savePath, err := zenity.SelectFileSave(zenity.Filename(filepath.Base(p)), zenity.ConfirmOverwrite())
	if err == nil {
		lo.Must0(os.WriteFile(savePath, out.Bytes(), 0644))
		util.OpenFileWithDefaultProgram(savePath)
	} else {
		util.OpenFileWithDefaultProgram(p)
	}
}

func translateChinese(tei network.GrobidTEI) []string {
	type TransSeg struct {
		Index   int
		English string
		Chinese string
	}
	title := tei.Header.FileDesc.TitleStmt.Title.Value
	abstract := strings.Join(lo.Map(tei.Header.ProfileDesc.Abstract.Div.Paragraphs, func(para network.Paragraph, index int) string {
		return para.Content
	}), "\n\n")
	var transSegArr []TransSeg
	transSegArr = append(transSegArr, TransSeg{
		Index:   0,
		English: "# " + title + "\n\n" + abstract,
		Chinese: "",
	})
	for i, div := range tei.Text.Body.Divs {
		content := strings.Join(lo.Map(div.Paragraphs, func(para network.Paragraph, index int) string {
			return para.Content
		}), "\n\n")
		heading := div.Head.N + " " + div.Head.Value
		transSegArr = append(transSegArr, TransSeg{
			Index:   i + 1,
			English: heading + "\n\n" + content,
			Chinese: "",
		})
	}
	slog.Info("Total segments", "count", len(transSegArr))
	wg := &sync.WaitGroup{}
	wg.Add(storage.Config.Concurrency)
	transSegChan := make(chan TransSeg)
	var resultTransSegArr []TransSeg
	resultLock := &sync.Mutex{}
	for range storage.Config.Concurrency {
		go func() {
			defer wg.Done()
			for transSeg := range transSegChan {
				slog.Info("Translating segment", "i", transSeg.Index)
				transSeg.Chinese = network.GetOpenAITranslation(transSeg.English)
				resultLock.Lock()
				resultTransSegArr = append(resultTransSegArr, transSeg)
				resultLock.Unlock()
				slog.Info("Successfully translated segment", "i", transSeg.Index)
			}
		}()
	}
	for _, transSeg := range transSegArr {
		transSegChan <- transSeg
	}
	close(transSegChan)
	wg.Wait()
	slices.SortFunc(resultTransSegArr, func(a, b TransSeg) int {
		return a.Index - b.Index
	})
	return lo.Map(resultTransSegArr, func(seg TransSeg, index int) string {
		return seg.Chinese
	})
}
