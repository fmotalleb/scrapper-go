package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseTable extracts table data, supporting both key-value and column-row formats
func ParseTable(htmlStr string) ([]map[string]any, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var headers []string
	var rows []map[string]any

	doc.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {
		columns := row.Find("td, th")

		// Handle key-value pair tables
		if columns.Length() == 2 && len(headers) == 0 {
			key := strings.TrimSpace(columns.Eq(0).Text())
			value := strings.TrimSpace(columns.Eq(1).Text())
			rows = append(rows, map[string]any{key: value})
			return
		}

		// Handle column-row formatted tables
		if rowIndex == 0 {
			headers = make([]string, 0, columns.Length())
			columns.Each(func(_ int, cell *goquery.Selection) {
				headers = append(headers, strings.TrimSpace(cell.Text()))
			})
		} else {
			rowData := make(map[string]any)
			columns.Each(func(colIndex int, cell *goquery.Selection) {
				if colIndex < len(headers) {
					rowData[headers[colIndex]] = strings.TrimSpace(cell.Text())
				}
			})
			rows = append(rows, rowData)
		}
	})

	return rows, nil
}
