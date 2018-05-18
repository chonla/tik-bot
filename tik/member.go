package tik

// Member is team member
type Member struct {
	id        string
	Name      string
	SlackName string
}

// Find a member
func (t *Tik) Find(uid string) (m *Member, e error) {
	dsnap, e := t.client.Collection("members").Doc(uid).Get(t.ctx)
	if e != nil {
		return nil, e
	}

	mem := Member{}
	dsnap.DataTo(&mem)
	m = &mem
	return
}

// NewMember to create a member
func (t *Tik) NewMember(m *Member) error {
	_, e := t.client.Collection("members").Doc(m.id).Set(t.ctx, m)
	return e
}
