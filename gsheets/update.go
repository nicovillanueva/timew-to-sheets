package gsheets

import (
	"google.golang.org/api/sheets/v4"
)

type UpdateBuilder struct {
	req *sheets.Request
}

func NewUpdateBuilder(x, y int64) *UpdateBuilder {
	return &UpdateBuilder{
		&sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "userEnteredValue",
				Start: &sheets.GridCoordinate{
					ColumnIndex: x,
					RowIndex:    y,
				},
				Rows: []*sheets.RowData{},
			},
		},
	}
}

func (ub *UpdateBuilder) AddRow(data ...string) {
	ub.req.UpdateCells.Rows = append(ub.req.UpdateCells.Rows, PrepareRowData(data...))
}

func (ub *UpdateBuilder) Request() *sheets.Request {
	return ub.req
}
