package timew

import (
	"strings"
	"time"
)

// TWTime extends time.Time to support ISO timestamps
type TWTime struct {
	time.Time
}

// UnmarshalJSON overrides time.Time.UnmarshalJSON so that it picks up TimeWarrior's time format
func (twt *TWTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		twt.Time = time.Time{}
		return nil
	}
	var err error
	twt.Time, err = time.Parse(timeFormatReference, s)
	return err
}

func (twt *TWTime) String() string {
	return twt.Time.Format("02/01/2006 15:04:05")
}
