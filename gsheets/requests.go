package gsheets

import (
	"github.com/nicovillanueva/timew-to-sheets/hashing"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
)

// RequestPool groups requests
type RequestPool []*sheets.Request

// AddMany takes many requests to add to the pool
func (rp *RequestPool) AddMany(reqs []*sheets.Request) {
	log.Debugf("adding %d requests to pool", len(reqs))
	*rp = append(*rp, reqs...)
}

// AddOne takes a single request to add to the pool
func (rp *RequestPool) AddOne(r *sheets.Request) {
	log.Debugf("Adding request to pool: %+v", r)
	*rp = append(*rp, r)
}

// getFormattingRequests returns an array of requests to apply styling to a sheet
func getFormattingRequests(sheetID int64) []*sheets.Request {
	return []*sheets.Request{
		{
			// Forced locale - TODO: Remove?
			UpdateSpreadsheetProperties: &sheets.UpdateSpreadsheetPropertiesRequest{
				Properties: &sheets.SpreadsheetProperties{
					Locale: "es_AR",
				},
				Fields: "locale",
			},
		},
		{
			// Frozen count in all sheets (new and old)
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId: sheetID,
					GridProperties: &sheets.GridProperties{
						FrozenRowCount: 1,
					},
				},
				Fields: "gridProperties.frozenRowCount",
			},
		},
		{
			// Header bold
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:        sheetID,
					EndRowIndex:    1,
					EndColumnIndex: 4,
				},
				Fields: "userEnteredFormat.textFormat.bold",
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						TextFormat: &sheets.TextFormat{
							Bold: true,
						},
					},
				},
			},
		},
		{
			// Totals bold
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    1,
					EndRowIndex:      2,
					StartColumnIndex: 5,
					EndColumnIndex:   6,
				},
				Fields: "userEnteredFormat.textFormat.bold",
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						TextFormat: &sheets.TextFormat{
							Bold: true,
						},
					},
				},
			},
		},
		{
			// Entries date format
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartColumnIndex: 0,
					EndColumnIndex:   2,
					EndRowIndex:      10,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						NumberFormat: &sheets.NumberFormat{
							Type:    "DATE_TIME",
							Pattern: "dd/mm/yyyy hh:mm",
						},
					},
				},
				Fields: "userEnteredFormat.numberFormat",
			},
		},
		{
			// Delta durations format
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    1,
					StartColumnIndex: 2,
					EndColumnIndex:   3,
				},
				Cell: &sheets.CellData{
					UserEnteredFormat: &sheets.CellFormat{
						NumberFormat: &sheets.NumberFormat{
							Type:    "TIME",
							Pattern: "[h]:mm:ss",
						},
					},
				},
				Fields: "userEnteredFormat.numberFormat",
			},
		},
		{
			// Resize entries columns - not working though?
			AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
				Dimensions: &sheets.DimensionRange{
					SheetId:    sheetID,
					Dimension:  "COLUMNS",
					StartIndex: 0,
					EndIndex:   4,
				},
			},
		},
	}
}

// getTotalsRequest returns a request that display the totals in a sheet; styling is applied by `getFormattingRequests()`
func getTotalsRequest(sheetID int64) *sheets.Request {
	b := NewUpdateRequestBuilder(sheetID, 5, 1)
	b.AddRow("Total", "#F#=SUM(C2:C1000)")
	return b.Request()
}

// NewStyledSheetRequests returns a slice of requests that create and style a new sheet
func NewStyledSheetRequests(name string) []*sheets.Request {
	hashedName := hashing.HashDate(name)
	requests := make([]*sheets.Request, 1)
	requests[0] = &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title:   name,
				SheetId: hashedName,
				GridProperties: &sheets.GridProperties{
					FrozenRowCount: 1,
				},
			},
		},
	}
	requests = append(requests, getFormattingRequests(hashedName)...)
	requests = append(requests, getTotalsRequest(hashedName))
	return requests
}
