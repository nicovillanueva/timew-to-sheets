package gsheets

import (
	"google.golang.org/api/sheets/v4"
	"strconv"
)

// GetFormattingRequests returns an array of requests to apply the sheet's format
func GetFormattingRequests() []*sheets.Request {
	return []*sheets.Request{
		&sheets.Request{
			UpdateSpreadsheetProperties: &sheets.UpdateSpreadsheetPropertiesRequest{
				Properties: &sheets.SpreadsheetProperties{
					Locale: "es_AR",
				},
				Fields: "locale",
			},
		},
		&sheets.Request{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					GridProperties: &sheets.GridProperties{
						FrozenRowCount: 1,
					},
				},
				Fields: "gridProperties.frozenRowCount",
			},
		},
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
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
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
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
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
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
		&sheets.Request{
			RepeatCell: &sheets.RepeatCellRequest{
				Range: &sheets.GridRange{
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

		&sheets.Request{
			AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
				Dimensions: &sheets.DimensionRange{
					Dimension:  "COLUMNS",
					StartIndex: 0,
					EndIndex:   1,
				},
			},
		},
	}
}

// GetTotalsRequests returns an array of requests that make up the Total counter
func GetTotalsRequests() []*sheets.Request {
	return []*sheets.Request{
		&sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Fields: "userEnteredValue",
				Range: &sheets.GridRange{
					StartColumnIndex: 5,
					StartRowIndex:    1,
				},
				Rows: []*sheets.RowData{
					&sheets.RowData{
						Values: []*sheets.CellData{
							&sheets.CellData{
								UserEnteredValue: &sheets.ExtendedValue{
									StringValue: "Total",
								},
							},
							&sheets.CellData{
								UserEnteredValue: &sheets.ExtendedValue{
									FormulaValue: "=SUM(C2:C1000)",
								},
							},
						},
					},
				},
			},
		},
	}
}

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

func PrepareRowData(cellData ...string) *sheets.RowData {
	return &sheets.RowData{
		Values: func() []*sheets.CellData {
			cd := make([]*sheets.CellData, len(cellData))
			for i, d := range cellData {
				cd[i] = &sheets.CellData{
					UserEnteredValue: getExtValue(d),
				}
			}
			return cd
		}(),
	}
}
