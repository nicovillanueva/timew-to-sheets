package timew

import (
	"os"
	"testing"
	"time"
)

const (
	twexport = "timew-export"
)

var parsedExport = Summary{
	{
		Start: TWTime{Time: time.Date(2018, time.December, 12, 13, 30, 00, 00, getUtcLocation())},
		End:   TWTime{Time: time.Date(2018, time.December, 12, 14, 45, 00, 00, getUtcLocation())},
		Tags:  []string{"work"},
	},
	{
		Start: TWTime{Time: time.Date(2018, time.December, 13, 13, 30, 00, 00, getUtcLocation())},
		End:   TWTime{Time: time.Date(2018, time.December, 14, 13, 30, 10, 00, getUtcLocation())},
		Tags:  []string{"work", "more work"},
	},
	{
		Start: TWTime{Time: time.Date(2018, time.December, 15, 13, 30, 00, 00, getUtcLocation())},
		End:   TWTime{Time: time.Date(2018, time.December, 17, 13, 30, 00, 00, getUtcLocation())},
		Tags:  []string{"getting high"},
	},
}

func getUtcLocation() *time.Location {
	location, _ := time.LoadLocation("")
	return location
}

func Test_parseTimeWarrior(t *testing.T) {
	f, err := os.Open(twexport)
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	_, summary := ParseTimeWarrior(f)

	if !summary.Equal(parsedExport) {
		t.Fail()
	}
}

func Test_SplitIntoMonths(t *testing.T) {
	location := getUtcLocation()
	s := Summary{
		{
			Start: TWTime{time.Date(2019, time.January, 01, 00, 00, 00, 00, location)},
			End:   TWTime{Time: time.Date(2019, time.January, 01, 02, 00, 00, 00, location)},
			Tags:  []string{"hello", "world", "jan"},
		},
		{
			Start: TWTime{Time: time.Date(2019, time.February, 01, 00, 00, 00, 00, location)},
			End:   TWTime{Time: time.Date(2019, time.February, 01, 02, 00, 00, 00, location)},
			Tags:  []string{"hello", "world", "feb"},
		},
	}
	split := s.SplitIntoMonths()
	if _, ok := split["January 2019"]; !ok {
		t.Fail()
	}
	if _, ok := split["February 2019"]; !ok {
		t.Fail()
	}
}
