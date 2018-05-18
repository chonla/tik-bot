package tik

import (
	"time"

	"cloud.google.com/go/firestore"
)

// Workplace carries state of conversation
type Workplace struct {
	id    string
	Names []string
}

// CheckInRecord is checkin record
type CheckInRecord struct {
	Multiplier float32
	Location   string
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
		w.Names = append(w.Names, name)
	}
	_, e = t.client.Collection("workplaces").Doc(id).Set(t.ctx, w)

	if e != nil {
		return e
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	date := time.Now().In(loc).Format("20060102")
	month := date[0:6]

	row := &CheckInRecord{
		Multiplier: 1.0,
		Location:   name,
	}

	rec := map[string]*CheckInRecord{}
	rec[date] = row

	recset := map[string]map[string]*CheckInRecord{}
	recset[month] = rec

	_, e = t.client.Collection("checkins").Doc(id).Set(t.ctx, recset, firestore.MergeAll)

	return e
}
