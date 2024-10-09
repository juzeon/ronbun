package main

import (
	"github.com/samber/lo"
	"log/slog"
	"ronbun/ccf"
	"ronbun/util"
	"slices"
)

func UpdateList() {
	conferenceSubs := util.PromptSelectConferenceSubs()
	conferenceRankings := util.PromptSelectConferenceRankings()
	//startYear := util.PromptInputStartYear()
	var conferences []ccf.Conference
	for _, sub := range conferenceSubs {
		conferences = append(conferences, ccf.GetConferencesBySub(sub.Sub)...)
	}
	conferences = lo.Filter(conferences, func(c ccf.Conference, index int) bool {
		return slices.Contains(conferenceRankings, c.Rank.CCF)
	})
	slog.Info("v", "v", conferences)
}
