package util

import (
	"github.com/ncruces/zenity"
	"github.com/samber/lo"
	"ronbun/ccf"
	"ronbun/db"
	"strconv"
)

func PromptSelectConferenceSubs() []ccf.ConferenceSub {
	conferenceSubs := ccf.GetConferenceSubs()
	lastConferenceSubs := db.GetSettingObj[[]string]("last_conference_subs")
	conferenceSubNames := lo.Map(conferenceSubs,
		func(item ccf.ConferenceSub, index int) string {
			return item.Name
		})
	selectedSubNames := Attempt(func() ([]string, error) {
		arr, err := zenity.ListMultiple("Select conference subs:", conferenceSubNames,
			zenity.DefaultItems(lastConferenceSubs...))
		if err != nil {
			return nil, err
		}
		return arr, nil
	})
	db.SetSettingObj[[]string]("last_conference_subs", selectedSubNames)
	var result []ccf.ConferenceSub
	for _, name := range selectedSubNames {
		if t, ok := lo.Find(conferenceSubs, func(item ccf.ConferenceSub) bool {
			return item.Name == name
		}); ok {
			result = append(result, t)
		}
	}
	return result
}
func PromptSelectConferenceRankings() []string {
	lastRankings := db.GetSettingObj[[]string]("last_conference_rankings")
	selectedRankings := Attempt(func() ([]string, error) {
		arr, err := zenity.ListMultiple("Select conference rankings:", []string{"A", "B", "C"},
			zenity.DefaultItems(lastRankings...))
		if err != nil {
			return nil, err
		}
		return arr, nil
	})
	db.SetSettingObj[[]string]("last_conference_rankings", selectedRankings)
	return selectedRankings
}
func PromptInputStartYear() int {
	lastStartYear := db.GetSettingObj[int]("last_conference_start_year")
	year := Attempt(func() (int, error) {
		text, err := zenity.Entry("Input a conference start year:",
			zenity.EntryText(strconv.Itoa(lastStartYear)))
		if err != nil {
			return 0, err
		}
		year, err := strconv.Atoi(text)
		if err != nil {
			return 0, err
		}
		return year, nil
	})
	db.SetSettingObj("last_conference_start_year", year)
	return year
}
func PromptInputSearchKeyword() string {
	lastSearchKeyword := db.GetSetting("last_search_keyword")
	keyword := Attempt(func() (string, error) {
		text, err := zenity.Entry("Input a search keyword:",
			zenity.EntryText(lastSearchKeyword))
		if err != nil {
			return "", err
		}
		return text, nil
	})
	db.SetSetting("last_search_keyword", keyword)
	return keyword
}
func PromptConfirmation(text string) {
	Attempt(func() (string, error) {
		if err := zenity.Question(text); err != nil {
			return "", err
		}
		return "", nil
	})
}
