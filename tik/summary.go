package tik

import (
	"fmt"
	"strings"
)

func (t *Tik) createSummaryReport(title string, r []*CheckInRecord) string {
	sum := map[string]float32{}
	for i, n := 0, len(r); i < n; i++ {
		sum[r[i].Location] += r[i].Multiplier
	}

	rows := []string{}
	for k, v := range sum {
		rows = append(rows, fmt.Sprintf("*%s* : %0.1f วัน", k, v))
	}

	return title + "\n" + strings.Join(rows, "\n")
}
