package mysqlexplainparser

import (
	"regexp"

	"github.com/moznion/go-text-visual-width"
)

type stringArray []string

func (rows *stringArray) removeUnnecessaryLinesOnTop() {
	_rows := *rows

	// Skip the blank line(s) that are on the top
	for _, row := range _rows {
		if row != "" {
			break
		}
		_rows = _rows[1:]
	}
	*rows = _rows[1:] // Skip the top of outline
}

func (rows *stringArray) removeUnnecessaryLinesOnBottom() {
	_rows := *rows

	regexToRemoveQueryResult := regexp.MustCompile("^[0-9]+[ \t]+rows") // e.g. "2 rows in set, 1 warning (0.00 sec)"
	regexToDetectLineOfTable := regexp.MustCompile("^\\+-+")

	// Skip bottom unnecessary line(s)
	for {
		lastPos := len(_rows) - 1
		lastLine := _rows[lastPos]
		_rows = _rows[:lastPos]

		if (!regexToRemoveQueryResult.MatchString(lastLine) && lastLine != "") || _rows == nil {
			if !regexToDetectLineOfTable.MatchString(lastLine) {
				_rows = _rows[:lastPos+1]
			}
			break
		}
	}

	*rows = _rows
}

func (rows *stringArray) removeUnnecessaryLinesOnBothEnds() {
	rows.removeUnnecessaryLinesOnTop()
	rows.removeUnnecessaryLinesOnBottom()
}

// Parse returns the result of parsed EXPLAIN as array of map
func Parse(explain string) []map[string]string {
	rows := stringArray(regexp.MustCompile("\r?\n").Split(explain, -1))
	rows.removeUnnecessaryLinesOnBothEnds()

	var indexRow string
	indexRow, rows = rows[0], rows[2:]
	//                              ~~ to skip separator between header and body

	regexpToStrip := regexp.MustCompile("^[ \t]*(.*?)[ \t]*$") // " hoge fuga  " => "hoge fuga" ($1)

	indexes := regexp.MustCompile("\\|").Split(indexRow, -1)
	lengths := make([]int, len(indexes))

	var i int
	for _, index := range indexes {
		if index == "" {
			continue
		}
		lengths[i] = len(index)
		indexes[i] = regexpToStrip.ReplaceAllString(index, "$1")
		i++
	}

	parsed := make([]map[string]string, len(rows))
	for i, row := range rows {
		parsedRow := make(map[string]string)
		for j, index := range indexes {
			var column string
			row = row[1:]
			if row == "" {
				break
			}
			column, row = visualwidth.Separate(row, lengths[j])
			column = regexpToStrip.ReplaceAllString(column, "$1")
			parsedRow[index] = column
		}
		parsed[i] = parsedRow
	}

	return parsed
}

// ParseVertical returns the result of parsed EXPLAIN as vertical as array of map
func ParseVertical(explain string) []map[string]string {
	rows := stringArray(regexp.MustCompile("\r?\n").Split(explain, -1))
	rows.removeUnnecessaryLinesOnBothEnds()

	var parsed []map[string]string
	parsedRow := make(map[string]string)
	regexToDetectRowSeparator := regexp.MustCompile("^\\*+")
	regexToGetKeyValue := regexp.MustCompile("^[ \t]*(.+)?:([ \t]+(.+))?")

	for _, row := range rows {
		if regexToDetectRowSeparator.MatchString(row) {
			parsed = append(parsed, parsedRow)
			parsedRow = make(map[string]string)
			continue
		}
		matched := regexToGetKeyValue.FindStringSubmatch(row)
		key, value := matched[1], matched[3]
		parsedRow[key] = value
	}
	parsed = append(parsed, parsedRow)

	return parsed
}
