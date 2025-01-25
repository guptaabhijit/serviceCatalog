package core

import (
	"log"
	"strings"
)

// Add these utility functions
func formatSQL(sql string) string {
	// Remove extra spaces
	sql = strings.Join(strings.Fields(sql), " ")

	// Keywords to break line before
	keywords := []string{"SELECT", "FROM", "WHERE", "LEFT JOIN", "INNER JOIN", "GROUP BY", "ORDER BY", "LIMIT", "OFFSET"}

	// Add newlines before keywords
	for _, keyword := range keywords {
		sql = strings.ReplaceAll(sql, keyword, "\n"+keyword)
	}

	// Add indentation
	lines := strings.Split(sql, "\n")
	var formattedLines []string
	indent := ""

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// Adjust indentation based on keywords
		if strings.HasPrefix(trimmedLine, "SELECT") {
			indent = "  "
		} else if strings.HasPrefix(trimmedLine, "FROM") {
			indent = "  "
		} else if strings.HasPrefix(trimmedLine, "LEFT JOIN") ||
			strings.HasPrefix(trimmedLine, "INNER JOIN") {
			indent = "    "
		} else if strings.HasPrefix(trimmedLine, "WHERE") ||
			strings.HasPrefix(trimmedLine, "GROUP BY") ||
			strings.HasPrefix(trimmedLine, "ORDER BY") {
			indent = "  "
		}

		formattedLines = append(formattedLines, indent+trimmedLine)
	}

	return "\n" + strings.Join(formattedLines, "\n")
}

func LogQuery(prefix string, sql string) {
	formattedSQL := formatSQL(sql)
	log.Printf("\n%s: %s\n", prefix, formattedSQL)
}
