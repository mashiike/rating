package rating_test

import (
	"testing"

	"github.com/mashiike/golicko/rating"
)

type ParseTest struct {
	name     string
	format   string
	value    string
	expected rating.Rating
}

var parseTests = []ParseTest{
	{"StrengthOnly", rating.StrengthOnlyFormat, "1700.082", rating.New(1700.082, 350.0, 0.06)},
	{"IntStrengthOnly", rating.StrengthOnlyFormat, "1320", rating.New(1320.00, 350.0, 0.06)},
	{"WithRange", rating.WithRangeFormat, "1320.00 (1170.0-1470.0)", rating.New(1320.00, 75.0, 0.06)},
	{"CSV", rating.CSVFormat, "1320,75,0.06", rating.New(1320.00, 75.0, 0.06)},
	{"Detail", rating.DetailFormat, "1320 (1170-1470 v=0.25)", rating.New(1320.00, 75.0, 0.25)},
	{"PlusMinus", rating.PlusMinusFormat, "1320p-150.0", rating.New(1320.00, 75.0, 0.06)},
	{"Default", rating.DefaultFormat, "1320p-150.0v=0.25", rating.New(1320.00, 75.0, 0.25)},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		got, err := rating.Parse(test.format, test.value)
		if err != nil {
			t.Errorf("%s error: %v", test.name, err)
		}
		if got != test.expected {
			t.Errorf("%s unexpected: got is %+v", test.name, got)
		}
	}
}

var formatTests = []FormatTest{
	{"StrengthOnly", rating.StrengthOnlyFormat, "1500.0"},
	{"WithRange", rating.WithRangeFormat, "1500.0 (800.0-2200.0)"},
	{"CSV", rating.CSVFormat, "1500.0,350.0,0.06"},
	{"Detail", rating.DetailFormat, "1500.0 (800.0-2200.0 v=0.06)"},
	{"PlusMinus", rating.PlusMinusFormat, "1500.0p-700.0"},
	{"Default", rating.DefaultFormat, "1500.0p-700.0v=0.06"},
}

type FormatTest struct {
	name   string
	format string
	result string
}

func TestFormat(t *testing.T) {
	player := rating.Default(0.06)
	for _, test := range formatTests {
		result := player.Format(test.format)
		if result != test.result {
			t.Errorf("%s expected %q got %q", test.name, test.result, result)
		}
	}
}
