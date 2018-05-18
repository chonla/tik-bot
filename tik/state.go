package tik

import (
	"time"
)

// ConversationState carries state of conversation
type ConversationState struct {
	id         string
	State      string
	Expiration int64
}

const ttl = "2m"

// FindState to find a state
func (t *Tik) FindState(id string) (c *ConversationState, e error) {
	dsnap, e := t.client.Collection("states").Doc(id).Get(t.ctx)
	if e != nil {
		return nil, nil
	}

	cs := ConversationState{}
	dsnap.DataTo(&cs)

	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc).Unix()

	if cs.Expiration < now {
		return nil, nil
	}
	c = &cs

	return
}

// SetState to create a new state
func (t *Tik) SetState(c *ConversationState) error {
	loc, _ := time.LoadLocation("UTC")
	dur, _ := time.ParseDuration(ttl)
	c.Expiration = time.Now().In(loc).Add(dur).Unix()
	_, e := t.client.Collection("states").Doc(c.id).Set(t.ctx, c)
	return e
}

// ClearState to remove state
func (t *Tik) ClearState(id string) error {
	_, e := t.client.Collection("states").Doc(id).Delete(t.ctx)
	return e
}
