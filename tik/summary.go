package tik

import (
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/iterator"
)

type monthlyCheckInRecord map[string][]*CheckInRecord

func (t *Tik) createSummaryReport(r []*CheckInRecord, startDate, endDate string) string {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	start, _ := time.ParseInLocation("20060102", startDate, loc)
	from := start.Format("2 Jan 2006")
	end, _ := time.ParseInLocation("20060102", endDate, loc)
	to := end.Format("2 Jan 2006")

	sum := map[string]float64{}
	for i, n := 0, len(r); i < n; i++ {
		sum[r[i].Location] += r[i].Multiplier
	}

	rows := []string{}
	for k, v := range sum {
		rows = append(rows, fmt.Sprintf("*%s* : %0.1f วัน", k, v))
	}

	return fmt.Sprintf("สรุปรอบเงินเดือน %s - %s", from, to) + "\n" + strings.Join(rows, "\n")
}

func (t *Tik) getSummaryReport(uid string) (cir []*CheckInRecord, start string, end string, e error) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	ts := time.Now().In(loc)
	thisMonth := ts.Format("200601")

	firstOfThisMonth, _ := time.Parse("20060102", thisMonth+"01")
	lastMonthTS := firstOfThisMonth.Add(-24 * time.Hour)
	lastMonth := lastMonthTS.Format("200601")

	start = fmt.Sprintf("%s%02d", lastMonth, 26)
	end = fmt.Sprintf("%s%02d", thisMonth, 25)

	lastCir, e := t.getMonthlyReport(uid, lastMonth, 26, 31)
	if e == nil {
		cir = append(cir, lastCir...)
	}

	thisCir, e := t.getMonthlyReport(uid, thisMonth, 1, 25)
	if e == nil {
		cir = append(cir, thisCir...)
	}

	return
}

func (t *Tik) getMonthlyReport(uid, month string, start int, end int) (cir []*CheckInRecord, e error) {
	rangeStart := fmt.Sprintf("%s%02d", month, start)
	rangeEnd := fmt.Sprintf("%s%02d", month, end)
	iter := t.client.Collection("checkins").Doc(uid).Collection("monthly").Documents(t.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		recSet := map[string]*CheckInRecord{}

		doc.DataTo(&recSet)

		for _, v := range recSet {
			if rangeStart <= v.Stamp && v.Stamp <= rangeEnd {
				cir = append(cir, v)
			}
		}
	}

	return
}
