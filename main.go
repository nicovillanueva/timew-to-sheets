package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/nicovillanueva/timew-to-sheets/gsheets"
	"github.com/nicovillanueva/timew-to-sheets/hashing"
	"github.com/nicovillanueva/timew-to-sheets/timew"
)

const (
	sheetIDEnvVar = "TW_SPREADSHEET_ID"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	spreadsheetID := os.Getenv(sheetIDEnvVar)
	if spreadsheetID == "" {
		log.Fatal("spreadsheet ID not provided; set it via an environment variable named ", sheetIDEnvVar)
	}
	_, summary := timew.ParseTimeWarrior(os.Stdin)
	rp := gsheets.RequestPool{}
	for month, entries := range summary.SplitIntoMonths() {
		rp.AddMany(gsheets.NewStyledSheetRequests(month))
		rp.AddOne(entries.Request(hashing.HashDate(month)))
	}

	err := gsheets.DoBatchUpdate(spreadsheetID, rp)
	if err != nil {
		log.Fatalf("error performing batch update: %v", err)
	}
	log.Infof("all done")
}
