package ai

import (
	"fmt"
	"sort"
	"strings"

	"github.com/emmanuel326/gopher_intel/internal/fetcher"
)

// LocalDigest produces a structured digest from raw messages without any AI API
func LocalDigest(data map[string][]fetcher.Message) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(" GOPHER INTEL — LOCAL DIGEST (AI unavailable)\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	for source, messages := range data {
		sb.WriteString(fmt.Sprintf("━━━ %s ━━━\n", strings.ToUpper(source)))
		sb.WriteString(fmt.Sprintf("total messages: %d\n\n", len(messages)))

		// --- activity by author ---
		authorCount := map[string]int{}
		for _, m := range messages {
			authorCount[m.Author]++
		}
		type authorStat struct {
			name  string
			count int
		}
		var authors []authorStat
		for name, count := range authorCount {
			authors = append(authors, authorStat{name, count})
		}
		sort.Slice(authors, func(i, j int) bool {
			return authors[i].count > authors[j].count
		})

		sb.WriteString("TOP CONTRIBUTORS:\n")
		limit := 5
		if len(authors) < limit {
			limit = len(authors)
		}
		for _, a := range authors[:limit] {
			sb.WriteString(fmt.Sprintf("  %-30s %d messages\n", a.name, a.count))
		}
		sb.WriteString("\n")

		// --- classify messages ---
		var patches, rfcs, bugs, replies, series []fetcher.Message
		for _, m := range messages {
			subj := strings.ToLower(m.Subject)
			switch {
			case strings.HasPrefix(subj, "re:"):
				replies = append(replies, m)
			case strings.Contains(subj, "rfc"):
				rfcs = append(rfcs, m)
			case strings.Contains(subj, "bug") || strings.Contains(subj, "fix") || strings.Contains(subj, "regression"):
				bugs = append(bugs, m)
			case strings.Contains(subj, "[patch 0/") || strings.Contains(subj, "[patch v"):
				series = append(series, m)
			case strings.Contains(subj, "[patch]") || strings.Contains(subj, "[patch "):
				patches = append(patches, m)
			}
		}

		sb.WriteString(fmt.Sprintf("MESSAGE BREAKDOWN:\n"))
		sb.WriteString(fmt.Sprintf("  patches:      %d\n", len(patches)))
		sb.WriteString(fmt.Sprintf("  patch series: %d\n", len(series)))
		sb.WriteString(fmt.Sprintf("  rfcs:         %d\n", len(rfcs)))
		sb.WriteString(fmt.Sprintf("  bug/fix:      %d\n", len(bugs)))
		sb.WriteString(fmt.Sprintf("  replies:      %d\n", len(replies)))
		sb.WriteString("\n")

		// --- notable items ---
		if len(rfcs) > 0 {
			sb.WriteString("RFCs:\n")
			for _, m := range rfcs {
				sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", m.Date.Format("Jan 02"), m.Author, m.Subject))
			}
			sb.WriteString("\n")
		}

		if len(bugs) > 0 {
			sb.WriteString("BUG/FIX SIGNALS:\n")
			for _, m := range bugs {
				sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", m.Date.Format("Jan 02"), m.Author, m.Subject))
			}
			sb.WriteString("\n")
		}

		if len(series) > 0 {
			sb.WriteString("PATCH SERIES:\n")
			for _, m := range series {
				sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", m.Date.Format("Jan 02"), m.Author, m.Subject))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString("tip: run again later when Gemini quota resets for full AI analysis\n")

	return sb.String()
}
