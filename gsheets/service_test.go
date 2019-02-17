package gsheets

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
	"testing"
)

func Test_cancelPreExistingSheets(t *testing.T) {
	existing := []*sheets.Sheet{
		{
			Properties: &sheets.SheetProperties{
				SheetId: 1,
			},
		},
		{
			Properties: &sheets.SheetProperties{
				SheetId: 2,
			},
		},
		{
			Properties: &sheets.SheetProperties{
				SheetId: 3,
			},
		},
		{
			Properties: &sheets.SheetProperties{
				SheetId: 4,
			},
		},
		{
			Properties: &sheets.SheetProperties{
				SheetId: 5,
			},
		},
	}
	candidates := []*sheets.Request{
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 6,
				},
			},
		},
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 4,
				},
			},
		},
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 2,
				},
			},
		},
	}
	expected := []*sheets.Request{
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
		{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					SheetId: 6,
				},
			},
		},
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
		{
			UpdateCells: &sheets.UpdateCellsRequest{},
		},
	}

	// fmt.Printf("pre:  %+v\n", candidates)
	candidates = cancelPreExistingSheets(candidates, existing)
	// fmt.Printf("post: %+v\n", candidates)
	if len(candidates) != len(expected) {
		fmt.Printf("different lengths between candidates and expected (have %d, want %d)\n", len(candidates), len(expected))
		t.Fail()
	}

	for _, candidate := range candidates {
		if candidate.AddSheet != nil && candidate.AddSheet.Properties.SheetId != 6 {
			fmt.Printf("found non-matching SheetId in left over candidate (want %d, have %d)\n", 6, candidate.AddSheet.Properties.SheetId)
			t.Fail()
		}
	}
}
