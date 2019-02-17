package hashing

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var reverseTable map[string]int

func init() {
	reverseTable = make(map[string]int)
	for i := time.January; i <= time.December; i++ {
		reverseTable[i.String()] = int(i)
	}
}

// HashDate is an abomination
// Returns an int64 based off a "{month} {date}"
// "February 2017" returns "22017"
func HashDate(s string) int64 {
	s1 := strings.Split(s, " ")
	m := reverseTable[s1[0]]
	y := s1[1]
	h, _ := strconv.Atoi(fmt.Sprintf("%d%s", m, y))
	return int64(h)
}
