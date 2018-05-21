package tik

import (
	"time"

	"cloud.google.com/go/firestore"
	"github.com/chonla/slices"
)

// Workplace carries state of conversation
type Workplace struct {
	id    string
	Names []string
}

// CheckInRecord is checkin record
type CheckInRecord struct {
	Stamp                string
	Multiplier           float32
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
func (t *Tik) CheckIn(id string, name string) error {
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
		return e
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	ts := time.Now().In(loc)
	txstamp := ts.Format(time.RFC3339)
	date := ts.Format("20060102")
	month := date[0:6]

	row := &CheckInRecord{
		Stamp:                date,
		Multiplier:           1.0,
		Location:             name,
		TransactionTimestamp: txstamp,
	}

	k := name + "_" + date

	rec := map[string]*CheckInRecord{}
	rec[k] = row

	// recset := map[string]map[string]*CheckInRecord{}
	// recset[month] = rec

	_, e = t.client.Collection("checkins").Doc(id).Collection("monthly").Doc(month).Set(t.ctx, rec, firestore.MergeAll)

	return e
}
