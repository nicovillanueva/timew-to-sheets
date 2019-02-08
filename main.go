package main

import (
	"log"
	"os"

	gsheets "github.com/nicovillanueva/timew-to-sheets/gsheets"
	timew "github.com/nicovillanueva/timew-to-sheets/timew"
)

const (
	spreadsheetID = "1c1kBa8M_DaliMQqXaoUetP70Q7SGnPVxzUYes2IWtAo"
)

func main() {
	_, summary := timew.ParseTimeWarrior(os.Stdin) // TODO: Try closing if need to get input again
	r := append(gsheets.GetFormattingRequests(), gsheets.GetTotalsRequests()...)
	r = append(r, summary.Request())

	err := gsheets.DoBatchUpdate(spreadsheetID, r)
	if err != nil {
		log.Fatalf("error performing batch update: %v", err)
	}
}
