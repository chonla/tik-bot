package tik

import (
	"fmt"
	"regexp"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/chonla/slices"
	"github.com/kr/pretty"
)

// Workplace carries state of conversation
type Workplace struct {
	id    string
	Names []string
}

// CheckInRecord is checkin record
type CheckInRecord struct {
	Stamp                string
	Multiplier           float64
	Location             string
	TransactionTimestamp string
}

// FindWorkplace a workplace for member
func (t *Tik) FindWorkplace(id string) (w *Workplace, e error) {
	dsnap, e := t.client.Collection("workplaces").Doc(id).Get(t.ctx)
	if e != nil {
		return nil, e
	}

	work := Workplace{}
	dsnap.DataTo(&work)
	w = &work
	return
}

// CheckIn to add a new workplace or append to existing workplace
func (t *Tik) CheckIn(id string, name string) (string, error) {
	multiplier := 1.0
	if mul, text, ok := t.detectMultiplier(name); ok {
		name = text
		multiplier = mul
		fmt.Println("name is changed to:" + name)
		fmt.Printf("multiplier is changed to: %0.1f\n", multiplier)
	}

	w, e := t.FindWorkplace(id)
	if w == nil {
		w = &Workplace{
			id:    id,
			Names: []string{name},
		}
	} else {
		if !slices.ContainsString(w.Names, name) {
			w.Names = append(w.Names, name)
		}
	}
	_, e = t.client.Collection("workplaces").Doc(id).Set(t.ctx, w)

	if e != nil {
		return "", e
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	ts := time.Now().In(loc)
	txstamp := ts.Format(time.RFC3339)
	date := ts.Format("20060102")
	month := date[0:6]

	row := &CheckInRecord{
		Stamp:                date,
		Multiplier:           multiplier,
		Location:             name,
		TransactionTimestamp: txstamp,
	}

	k := name + "_" + date

	rec := map[string]*CheckInRecord{}
	rec[k] = row

	_, e = t.client.Collection("checkins").Doc(id).Collection("monthly").Doc(month).Set(t.ctx, rec, firestore.MergeAll)

	return name, e
}

func (t *Tik) detectMultiplier(s string) (mul float64, text string, ok bool) {
	supportedMultiplier := map[string][]string{
		"FULLDAY": []string{"(.+)(เต็มวัน)$", "(.+)( 1 วัน)$", "(.+)(ทั้งวัน)$"},
		"HALFDAY": []string{"(.+)(ครึ่งวัน)$", "(.+)( 0\\.5 วัน)$", "(.+)(เช้า)$", "(.+)(บ่าย)$"},
	}

	multiplier := map[string]float64{
		"FULLDAY": 1.0,
		"HALFDAY": 0.5,
	}

	for k, v := range supportedMultiplier {
		for i, n := 0, len(v); i < n; i++ {
			pretty.Println(v[i])
			re := regexp.MustCompile(v[i])
			result := re.FindStringSubmatch(s)
			pretty.Println(s)
			pretty.Println(result)
			if len(result) > 1 {
				text = result[1]
				mul = multiplier[k]
				ok = true
				return
			}
		}
	}

	return
}
