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
	Head     TranslationTmplDataSection
	Sections []TranslationTmplDataSection
}
type TranslationTmplDataSection struct {
	ChineseTitle   string
	EnglishTitle   string
	ChineseContent template.HTML
	EnglishContent template.HTML
}

func Translate() {
	pdfFile := util.PromptOpenPDFFile()
	v := lo.Must(os.ReadFile(pdfFile))
	hash := util.Sha1(v)
	translation := db.TranslationTx.MustFindOne("hash=?", hash)
	if translation == nil {
		slog.Info("Getting Grobid result...", "file", pdfFile)
		grobidResult := lo.Must(network.GetGrobidResult(v))
		translation = &db.Translation{
			Hash:       hash,
			GrobidData: grobidResult,
			TextData:   nil,
		}
		db.TranslationTx.MustCreate(translation)
	}
	if len(translation.TextData) == 0 {
		grobidData := util.ParseGrobidXML(translation.GrobidData)
		slog.Info("Translating...", "title", translation.Title)
		transSegArr := translateChinese(grobidData)
		translation.TextData = lo.Map(transSegArr, func(transSeg TransSeg, index int) db.TranslationTextData {
			return db.TranslationTextData{
				Chinese: transSeg.Chinese,
				English: transSeg.English,
			}
		})
		translation.Title = grobidData.Title
		db.TranslationTx.MustSave(translation)
	}
	tmpl := lo.Must(template.New("translation").Funcs(translationFuncsMap).Parse(translationTmpl))
	out := &bytes.Buffer{}
	var tmplData TranslationTmplData
	trimTitleLeft := func(str string) string {
		return strings.TrimLeft(str, "# ")
	}
	getTitleContent := func(text string) (string, template.HTML) {
		arr := strings.Split(text, "\n\n")
		arr = lo.Map(arr, func(line string, index int) string {
			return strings.TrimSpace(line)
		})
		arr = lo.Filter(arr, func(line string, index int) bool {
			return line != ""
		})
		return trimTitleLeft(arr[0]),
			template.HTML("<p>" + strings.Join(arr[1:], "</p><p>") + "</p>")
	}
	for i, textData := range translation.TextData {
		eTitle, eContent := getTitleContent(textData.English)
		cTitle, cContent := getTitleContent(textData.Chinese)
		section := TranslationTmplDataSection{
			ChineseTitle:   cTitle,
			EnglishTitle:   eTitle,
			ChineseContent: cContent,
			EnglishContent: eContent,
		}
		if i == 0 {
			tmplData.Head = section
			continue
		}
		tmplData.Sections = append(tmplData.Sections, section)
	}
	lo.Must0(tmpl.Execute(out, tmplData))
	p := storage.WriteTmpFile("Translation for "+filepath.Base(pdfFile)+".html", out.Bytes())
	savePath, err := zenity.SelectFileSave(
		zenity.ConfirmOverwrite(),
		zenity.Filename(filepath.Join(filepath.Dir(pdfFile), filepath.Base(p))),
	)
	if err == nil {
		lo.Must0(os.WriteFile(savePath, out.Bytes(), 0644))
		util.OpenFileWithDefaultProgram(savePath)
	} else {
		util.OpenFileWithDefaultProgram(p)
	}
}

type TransSeg struct {
	Index   int
	English string
	Chinese string
}

func translateChinese(data util.GrobidData) []TransSeg {
	var transSegArr []TransSeg
	transSegArr = append(transSegArr, TransSeg{
		Index:   0,
		English: "# " + data.Title + "\n\n" + data.Abstract,
		Chinese: "",
	})
	for i, section := range data.Sections {
		transSegArr = append(transSegArr, TransSeg{
			Index:   i + 1,
			English: section.Title + "\n\n" + section.Content,
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
	return resultTransSegArr
}
