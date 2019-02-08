package timew

import (
	"fmt"
	"os"
	"testing"
	// "time"
)

const (
	twexport = "timew-export"
)

func Test_parseTimeWarrior(t *testing.T) {
	f, err := os.Open(twexport)
	if err != nil {
		t.FailNow()
	}
	defer f.Close()
	_, summary := ParseTimeWarrior(f)

	for _, s := range summary {
		fmt.Printf("%+v\n", s)
	}
}
