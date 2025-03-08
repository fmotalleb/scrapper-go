package utils

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseTable extracts table data into []map[string]any
func ParseTable(htmlStr string) ([]map[string]any, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var headers []string
	var rows []map[string]any

	// Select table rows
	doc.Find("tr").Each(func(rowIndex int, row *goquery.Selection) {

		var rowData map[string]any

		// Extract table headers
		if rowIndex == 0 {
			row.Find("th").Each(func(_ int, th *goquery.Selection) {
				headers = append(headers, strings.TrimSpace(th.Text()))
			})
		} else {
			rowData = make(map[string]any)
			row.Find("td").Each(func(colIndex int, td *goquery.Selection) {
				if colIndex < len(headers) { // Ensure it doesn't go out of bounds
					rowData[headers[colIndex]] = strings.TrimSpace(td.Text())
				}
			})
			rows = append(rows, rowData)
		}
	})
	return rows, nil
}
