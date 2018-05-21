package tik

import (
	"log"
	"os"
	"strings"

	"github.com/chonla/oddsy"
	"github.com/kr/pretty"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"

	firebase "firebase.google.com/go"
	"golang.org/x/net/context"
)

// Tik is tik
type Tik struct {
	ctx    context.Context
	client *firestore.Client
	logger *log.Logger
	token  string
}

// Configuration is for translator
type Configuration struct {
	GcpToken          string
	FirebaseProjectID string
}

// NewTik creates a new tik
func NewTik(conf *Configuration) *Tik {
	t := &Tik{
		token:  conf.GcpToken,
		logger: log.New(os.Stdout, "tik: ", log.Lshortfile|log.LstdFlags),
	}

	envToken := os.Getenv("GCP_TOKEN")
	if envToken != "" {
		t.SetToken(envToken)
	}

	ctx := context.Background()
	fconf := &firebase.Config{
		ProjectID:   conf.FirebaseProjectID,
		DatabaseURL: "https://" + conf.FirebaseProjectID + ".firebaseio.com",
	}

	a, e := firebase.NewApp(ctx, fconf, option.WithAPIKey(t.token))
	if e != nil {
		log.Fatalf("Cannot create app: %v", e)
	}

	c, e := a.Firestore(ctx)
	if e != nil {
		log.Fatalf("Cannot create client: %v", e)
	}

	t.ctx = ctx
	t.client = c

	t.compileCommands()

	return t
}

// SetToken overrides configure token
func (t *Tik) SetToken(s string) {
	t.token = s
	t.logger.Println("GCP token is overwritten by environment variable.")
}

// Release to release client
func (t *Tik) Release() {
	t.client.Close()
}

// Dispatch to work with message
func (t *Tik) Dispatch(o *oddsy.Oddsy, m *oddsy.Message) {
	state, _ := t.FindState(m.From.UID)

	if state != nil {
		switch state.State {
		case "identify":
			pretty.Println(m.From)
			mem := &Member{
				id:        m.From.UID,
				SlackName: m.From.Name,
				Name:      m.Message,
			}
			e := t.NewMember(mem)
			if e != nil {
				o.Send(m.Channel.UID, "ขอโทษทีนะ"+mem.Name+" เหมือนจะเจอปัญหาล่ะ: "+e.Error())
			} else {
				o.Send(m.Channel.UID, "สวัสดีจ้ะ"+mem.Name)
			}
			t.ClearState(m.From.UID)
		case "workplace":
			wp, _ := t.CheckIn(m.From.UID, m.Message)
			o.Send(m.Channel.UID, "ลงชื่อเข้าทำงานที่ "+wp+" เรียบร้อยจ้ะ")
			t.ClearState(m.From.UID)
		}
	} else {
		w, e := t.Find(m.From.UID)
		if e != nil {
			t.SetState(&ConversationState{
				id:    m.From.UID,
				State: "identify",
			})
			o.Send(m.Channel.UID, "ชื่ออะไรเหรอ")
		} else {
			cmd := strings.ToLower(m.Message)

			cmdKey, cmdParams := t.discover(cmd)

			switch cmdKey {
			case "HELP":
				o.Send(m.Channel.UID, t.getHelp())
			case "GREETING":
				o.Send(m.Channel.UID, "สวัสดีจ้ะ"+w.Name)
			case "SUMMARY":
				rep, start, end, _ := t.getSummaryReport(m.From.UID)
				o.Send(m.Channel.UID, t.createSummaryReport(rep, start, end))
			case "CHECKIN":
				l := cmdParams[1]
				if l == "" {
					w, _ := t.FindWorkplace(m.From.UID)
					if w != nil && len(w.Names) == 1 {
						// Auto checkin if workplace is one place
						wp, _ := t.CheckIn(m.From.UID, w.Names[0])
						o.Send(m.Channel.UID, "ลงชื่อเข้าทำงานที่ "+wp+" เรียบร้อยจ้ะ")
					} else {
						if w != nil && len(w.Names) > 1 {
							workplaces := []*oddsy.SelectionOption{}

							for i, n := 0, len(w.Names); i < n; i++ {
								workplaces = append(workplaces, &oddsy.SelectionOption{
									Label: w.Names[i],
									Value: "checkin " + w.Names[i],
								})
							}

							o.SendSelection(m.Channel.UID, "วันนี้เข้าทำงานที่ไหนเหรอ", "เลือกสถานที่ทำงาน", workplaces)
						} else {
							t.SetState(&ConversationState{
								id:    m.From.UID,
								State: "workplace",
							})
							o.Send(m.Channel.UID, "วันนี้เข้าทำงานที่ไหนเหรอ")
						}
					}
				} else {
					wp, _ := t.CheckIn(m.From.UID, l)
					o.Send(m.Channel.UID, "ลงชื่อเข้าทำงานที่ "+wp+" เรียบร้อยจ้ะ")
				}
			default:
				o.Send(m.Channel.UID, "ไม่เข้าใจเลยล่ะ ลองใหม่นะ"+w.Name)
			}
		}
	}
}
