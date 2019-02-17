package hashing

import (
	"fmt"
	"testing"
	"time"
)

func sequentialSearch(arr *[]int64, number int64) bool {
	for _, v := range *arr {
		if v == number {
			return true
		}
	}
	return false
}

func Test_djb2_collision(t *testing.T) {
	yearFrom := 1920
	yearTo := 2319
	// yearTo := 1921
	total := ((yearTo - yearFrom + 1) * 12)
	tests := make([]int64, total)
	var idx int64
	for i := yearFrom; i <= yearTo; i++ {
		for month := time.January; month <= 12; month++ {
			s := fmt.Sprintf("%s %d", month.String(), i)
			h := Djb2(s, true)
			if sequentialSearch(&tests, h) {
				fmt.Printf("collision found: %s (%d)\n", s, h)
				t.Fail()
			}
			if h < 0 {
				fmt.Println("found negative:", h)
				t.Fail()
			}
			// fmt.Printf("(%d/%d) Original: %s | Hashed: %d\n", idx, total, s, h)
			tests[idx] = h
			idx++
		}
	}
}

func Test_craphash(t *testing.T) {
	yearFrom := 1920
	yearTo := 2319
	// yearTo := 1921
	total := ((yearTo - yearFrom + 1) * 12)
	tests := make([]int64, total)
	var idx int64
	for i := yearFrom; i <= yearTo; i++ {
		for month := time.January; month <= 12; month++ {
			s := fmt.Sprintf("%s %d", month.String(), i)
			h := HashDate(s)
			if sequentialSearch(&tests, h) {
				fmt.Printf("collision found: %s (%d)\n", s, h)
				t.Fail()
			}
			if h < 0 {
				fmt.Println("found negative:", h)
				t.Fail()
			}
			// fmt.Printf("(%d/%d) Original: %s | Hashed: %d\n", idx, total, s, h)
			tests[idx] = h
			idx++
		}
	}
}
