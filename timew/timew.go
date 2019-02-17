package timew

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	gsheets "github.com/nicovillanueva/timew-to-sheets/gsheets"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

const timeFormatReference = "20060102T150405Z"

// Config represents the read config properties from TimeWarrior
type Config map[string]string

// Summary lists all of the entries reported from TimeWarrior
type Summary []struct {
	Start TWTime   `json:"start"`
	End   TWTime   `json:"end"`
	Tags  []string `json:"tags"`
}

// Request takes all of the Summary information and builds a single `*sheets.Request`
// with a header and all entries.
func (s Summary) Request(sheetID int64) *sheets.Request {
	rb := gsheets.NewUpdateRequestBuilder(sheetID, 0, 0)
	rb.AddRow("Start", "End", "Delta", "Tags")
	for i, entry := range s {
		deltaFormula := fmt.Sprintf("#F#=(B%d-A%d)", i+2, i+2)
		rb.AddRow(entry.Start.String(), entry.End.String(), deltaFormula, strings.Join(entry.Tags, ", "))
	}
	return rb.Request()
}

func (s Summary) affectedMonth(idx int) string {
	year, month, _ := s[idx].Start.Date()
	dateTag := fmt.Sprintf("%s %d", month.String(), year)
	return dateTag
}

// SplitIntoMonths returns the summary data in a map keyed by "{month} {year}"
func (s Summary) SplitIntoMonths() map[string]Summary {
	m := make(map[string]Summary)
	for i, entry := range s {
		dateTag := s.affectedMonth(i)
		if _, ok := m[dateTag]; !ok {
			m[dateTag] = Summary{entry}
		} else {
			m[dateTag] = append(m[dateTag], entry)
		}
	}
	return m
}

// Equal implements a deep equality comparison between two Summary's
func (s Summary) Equal(another Summary) bool {
	for idx, entry := range another {
		if s[idx].Start.Time != entry.Start.Time {
			return false
		} else if s[idx].End.Time != entry.End.Time {
			return false
		}
		for tID, tag := range s[idx].Tags {
			if entry.Tags[tID] != tag {
				return false
			}
		}
	}
	return true
}

// ParseTimeWarrior reads a os.File (stdin for common usage; a file for tests) and parses
// the current configuration and the reported summary
func ParseTimeWarrior(in *os.File) (Config, Summary) {
	sc := bufio.NewScanner(in)
	config := Config{}
	var configDone = false
	var rawSummary string
	var summary Summary

	for sc.Scan() {
		input := sc.Text()
		if input == "" {
			configDone = true
			continue
		}
		if !configDone {
			line := strings.Split(input, " ")
			key := line[0][:len(line[0])-1]
			var value string
			if len(line) < 2 {
				value = ""
			} else {
				value = line[1]
			}
			config[key] = value
		} else {
			rawSummary += input
		}
	}
	if sc.Err() != nil {
		log.Fatalf("cannot read input from TW: %v", sc.Err())
	}
	if err := json.Unmarshal([]byte(rawSummary), &summary); err != nil {
		log.Fatalf("cannot unmarshal TW summary data: %v", err)
	}

	log.Debugf("read %d summary entries and %d config entries from TW", len(summary), len(config))
	return config, summary
}
