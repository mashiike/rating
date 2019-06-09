package rating

import (
	"errors"
	"strconv"
	"unicode"
)

const (
	_ = iota
	stdStrength
	stdLower
	stdUpper
	stdDeviation
	stdError
	stdVolatility
)

//Rating Format examples
const (
	StrengthOnlyFormat = "1500.0"
	WithRangeFormat    = "1500.0 (800.0-2200.0)"
	CSVFormat          = "1500.0,350.0,0.06"
	DetailFormat       = "1500.0 (800.0-2200.0 v=0.06)"
	PlusMinusFormat    = "1500.0p-700.0"
	DefaultFormat      = "1500.0p-700.0v=0.06"
)

var errBad = errors.New("bad value for field")

// ParseError describes a problem parsing a time string.
type ParseError struct {
	Layout     string
	Value      string
	LayoutElem string
	ValueElem  string
	Message    string
}

func quote(s string) string {
	return "\"" + s + "\""
}

// Error returns the string representation of a ParseError.
func (e *ParseError) Error() string {
	if e.Message == "" {
		return "parsing rating " +
			quote(e.Value) + " as " +
			quote(e.Layout) + ": cannot parse " +
			quote(e.ValueElem) + " as " +
			quote(e.LayoutElem)
	}
	return "parsing rating " +
		quote(e.Value) + e.Message
}

//Parse parses a formatted string and returns the rating value it represents.
//if not include volatility in layout, volatility set 0.06
func Parse(layout, value string) (Rating, error) {
	return parse(layout, value, 0.06)
}

//ParseWithVolatility is parse a formatted string with default volatility.
func ParseWithVolatility(layout, value string, volatility float64) (Rating, error) {
	return parse(layout, value, volatility)
}

func parse(layout, value string, defaultVolatility float64) (Rating, error) {
	alayout, avalue := layout, value
	var (
		strength   = float64(centerValue)
		deviation  float64
		volatility = defaultVolatility
		lower      float64
		upper      float64
	)

	for {
		var err error
		prefix, std, suffix := nextStdChunk(layout)
		value, err = skip(value, prefix)
		if err != nil {
			return Rating{}, &ParseError{alayout, avalue, prefix, value, ""}
		}
		if std == 0 {
			if len(value) != 0 {
				return Rating{}, &ParseError{alayout, avalue, "", value, ": extra text: " + value}
			}
			break
		}
		layout = suffix
		var fval float64
		fval, value, err = extractFloat(value)
		if err != nil {
			return Rating{}, &ParseError{alayout, avalue, prefix, value, ""}
		}
		switch std {
		case stdStrength:
			strength = fval
		case stdLower:
			lower = fval
		case stdUpper:
			upper = fval
		case stdDeviation:
			deviation = fval
		case stdError:
			deviation = fval / 2.0
		case stdVolatility:
			volatility = fval
		}
	}
	if deviation == 0.0 {
		deviation = (upper - lower) / 4.0
		if deviation == 0.0 {
			deviation = InitialDeviation
		}
	}
	return New(strength, deviation, volatility), nil
}

type formatElem struct {
	face     string
	stdChunk int
}

func (e formatElem) first() int {
	return int(e.face[0])
}

func (e formatElem) isMatch(i int, layout string) bool {
	return len(layout) >= i+len(e.face) && layout[i:i+len(e.face)] == e.face
}

func (e formatElem) separate(i int, layout string) (prefix string, std int, suffix string) {
	return layout[0:i], e.stdChunk, layout[i+len(e.face):]
}

var elems = []formatElem{
	{"0.06", stdVolatility},
	{"1500.0", stdStrength},
	{"2200.0", stdUpper},
	{"350.0", stdDeviation},
	{"700.0", stdError},
	{"800.0", stdLower},
}

func nextStdChunk(layout string) (prefix string, std int, suffix string) {
	for i := 0; i < len(layout); i++ {
		for _, elem := range elems {
			if c := int(layout[i]); c == elem.first() {
				if elem.isMatch(i, layout) {
					return elem.separate(i, layout)
				}
				continue
			}
		}
	}
	return layout, 0, ""
}

func cutspace(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	return s
}

func skip(value, prefix string) (string, error) {
	for len(prefix) > 0 {
		if prefix[0] == ' ' {
			if len(value) > 0 && value[0] != ' ' {
				return value, errBad
			}
			prefix = cutspace(prefix)
			value = cutspace(value)
			continue
		}
		if len(value) == 0 || value[0] != prefix[0] {
			return value, errBad
		}
		prefix = prefix[1:]
		value = value[1:]
	}
	return value, nil
}

func extractFloat(value string) (float64, string, error) {
	isPointed := false
	i := 0
	for ; i < len(value); i++ {
		if isPointed == false && value[i] == '.' {
			isPointed = true
			continue
		}
		if !unicode.IsDigit(rune(value[i])) {
			break
		}
	}
	fval, err := strconv.ParseFloat(value[0:i], 64)
	return fval, value[i:], err
}

// AppendFormat is like Format but appends the textual
// as same as time.Time
func (r Rating) AppendFormat(b []byte, layout string) []byte {
	lower, upper := r.Interval()
	for layout != "" {
		prefix, std, suffix := nextStdChunk(layout)
		if prefix != "" {
			b = append(b, prefix...)
		}
		if std == 0 {
			break
		}
		layout = suffix
		var value string
		switch std {
		case stdStrength:
			value = strconv.FormatFloat(r.Strength(), 'f', 1, 64)
		case stdLower:
			value = strconv.FormatFloat(lower, 'f', 1, 64)
		case stdUpper:
			value = strconv.FormatFloat(upper, 'f', 1, 64)
		case stdDeviation:
			value = strconv.FormatFloat(r.Deviation(), 'f', 1, 64)
		case stdError:
			value = strconv.FormatFloat(r.Deviation()*2.0, 'f', 1, 64)
		case stdVolatility:
			value = strconv.FormatFloat(r.Volatility(), 'f', -1, 64)
		}
		b = append(b, value...)
	}
	return b
}

// Format returns a textual representation of the time value formatted
// as same as time.Time
func (r Rating) Format(layout string) string {
	const bufSize = 64
	var b []byte
	max := len(layout) + 10
	if max < bufSize {
		var buf [bufSize]byte
		b = buf[:0]
	} else {
		b = make([]byte, 0, max)
	}
	b = r.AppendFormat(b, layout)
	return string(b)
}

//String is for dump. fmt.Stringer interface implements
//format is DetailFormat as 1500.0 (800.0-2200.0 v=0.06)
func (r Rating) String() string {
	return r.Format(DetailFormat)
}
