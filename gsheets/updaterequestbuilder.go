package gsheets

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/sheets/v4"
	"strconv"
)

// UpdateRequestBuilder provides an interface to easily and progressively build a range of values
// First get a reference using NewUpdateRequestBuilder and sequentially add rows using UpdateRequestBuilder.AddRow(...string),
// then call UpdateRequestBuilder.Request() to get the request ready for using with BatchUpdate
// See gsheets.DoBatchUpdate()
type UpdateRequestBuilder struct {
	req *sheets.Request
}

// NewUpdateRequestBuilder creates a new UpdateRequestBuilder for a sheet, starting at a certain coordinate
func NewUpdateRequestBuilder(sheetID, x, y int64) *UpdateRequestBuilder {
	return &UpdateRequestBuilder{
		&sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "userEnteredValue",
				Start: &sheets.GridCoordinate{
					SheetId:     sheetID,
					ColumnIndex: x,
					RowIndex:    y,
				},
				Rows: []*sheets.RowData{},
			},
		},
	}
}

// AddRow gets a variable amount of strings, and builds the row off of them
func (ub *UpdateRequestBuilder) AddRow(data ...string) {
	log.Debugf("creating row: %v", data)
	ub.req.UpdateCells.Rows = append(ub.req.UpdateCells.Rows, prepareRowData(data...))
}

// Request returns a *sheets.Request ready to do a BatchUpdate with
func (ub *UpdateRequestBuilder) Request() *sheets.Request {
	return ub.req
}

// prepareRowData takes a slice of strings, and sorts them into cells of a row
// See `getExtValue` for information on cell types
func prepareRowData(cellData ...string) *sheets.RowData {
	return &sheets.RowData{
		Values: func() []*sheets.CellData {
			cd := make([]*sheets.CellData, len(cellData))
			for i, d := range cellData {
				if d == "" {
					continue
				}
				cd[i] = &sheets.CellData{
					UserEnteredValue: getExtValue(d),
				}
			}
			return cd
		}(),
	}
}

// getExtValue sets the kind of data a cell will have.
// Different prefixes have different cell types (bool, number, formula or string)
func getExtValue(data string) *sheets.ExtendedValue {
	prefix := data[:3]
	switch prefix {
	case "#F#":
		return &sheets.ExtendedValue{
			FormulaValue: data[3:],
		}
	case "#B#":
		b, err := strconv.ParseBool(data[3:])
		if err != nil {
			b = false
		}
		return &sheets.ExtendedValue{
			BoolValue: b,
		}
	case "#S#":
		return &sheets.ExtendedValue{
			StringValue: data[3:],
		}
	case "#N#":
		n, err := strconv.Atoi(data[3:])
		if err != nil {
			n = 0
		}
		return &sheets.ExtendedValue{
			NumberValue: float64(n),
		}
	default:
		return &sheets.ExtendedValue{
			StringValue: data,
		}
	}
}
