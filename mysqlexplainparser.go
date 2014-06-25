package mysqlexplainparser

import (
	"regexp"

	"github.com/moznion/go-text-visual-width"
)

// Parse returns the result of parsed EXPLAIN as array of map
func Parse(explain string) []map[string]string {
	var rows = regexp.MustCompile("\r?\n").Split(explain, -1)

	// Skip the blank line(s) that are on the top
	for _, row := range rows {
		if row != "" {
			break
		}
		rows = rows[1:]
	}
	rows = rows[1:] // Skip the top of outline

	// Skip bottom unnecessary line(s)
	for {
		lastPos := len(rows) - 1
		lastLine := rows[lastPos]
		rows = rows[:lastPos]

		regexToRemoveQueryResult := regexp.MustCompile("^[0-9]+[ \t]+rows") // e.g. "2 rows in set, 1 warning (0.00 sec)"
		if (!regexToRemoveQueryResult.MatchString(lastLine) && lastLine != "") || rows == nil {
			break
		}
	}

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
		var parsedRow = make(map[string]string)
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
