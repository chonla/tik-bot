package tik

import "github.com/chonla/oddsy"

func (t *Tik) createWizardSelection() []*oddsy.SelectionOption {
	menus := []*oddsy.SelectionOption{
		&oddsy.SelectionOption{
			Label: "ลงชื่อเข้าทำงาน",
			Value: "check-in",
		},
		&oddsy.SelectionOption{
			Label: "สรุปวันรอบเงินเดือนเดือนนี้",
			Value: "sum",
		},
	}

	return menus
}
