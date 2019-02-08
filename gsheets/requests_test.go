package gsheets

import (
	"google.golang.org/api/sheets/v4"
	"testing"
)

func Test_getExtValue(t *testing.T) {
	td := []struct {
		input  string
		output *sheets.ExtendedValue
	}{
		{
			"#F#=SUM(2,2)",
			&sheets.ExtendedValue{
				FormulaValue: "=SUM(2,2)",
			},
		},
		{
			"#S#somestring",
			&sheets.ExtendedValue{
				StringValue: "somestring",
			},
		},
		{
			"bare",
			&sheets.ExtendedValue{
				StringValue: "bare",
			},
		},
		{
			"#B#false",
			&sheets.ExtendedValue{
				BoolValue: false,
			},
		},
		{
			"#B#True",
			&sheets.ExtendedValue{
				BoolValue: true,
			},
		},
		{
			"#B#1",
			&sheets.ExtendedValue{
				BoolValue: true,
			},
		},
		{
			"#N#42",
			&sheets.ExtendedValue{
				NumberValue: 42,
			},
		},
	}
	for _, testCase := range td {
		r := getExtValue(testCase.input)
		switch testCase.input[1] {
		case 'F':
			if r.FormulaValue != testCase.output.FormulaValue {
				t.Fail()
			}
		case 'S':
			if r.StringValue != testCase.output.StringValue {
				t.Fail()
			}
		case 'B':
			if r.BoolValue != testCase.output.BoolValue {
				t.Fail()
			}
		case 'N':
			if r.NumberValue != testCase.output.NumberValue {
				t.Fail()
			}
		}
	}
}
