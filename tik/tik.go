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
			t.CheckIn(m.From.UID, m.Message)
			o.Send(m.Channel.UID, "ลงชื่อเข้าทำงานที่ "+m.Message+" เรียบร้อยจ้ะ")
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
			switch {
			case t.isCheckIn(cmd):
				w, _ := t.FindWorkplace(m.From.UID)
				if w != nil && len(w.Names) == 1 {
					// Auto checkin if workplace is one place
					e = t.CheckIn(m.From.UID, w.Names[0])
					pretty.Println(e)
					o.Send(m.Channel.UID, "ลงชื่อเข้าทำงานที่ "+w.Names[0]+" เรียบร้อยจ้ะ")
				} else {
					if w != nil && len(w.Names) > 1 {
						o.Send(m.Channel.UID, "ทำงานมากกว่า 1 ที่")
					} else {
						t.SetState(&ConversationState{
							id:    m.From.UID,
							State: "workplace",
						})
						o.Send(m.Channel.UID, "วันนี้เข้าทำงานที่ไหนเหรอ")
					}
				}
			case t.isGreeting(cmd):
				o.Send(m.Channel.UID, "สวัสดีจ้ะ"+w.Name)
			default:
				o.Send(m.Channel.UID, "ไม่เข้าใจเลยล่ะ ลองใหม่นะ"+w.Name)
			}
		}
	}
}

func (t *Tik) isCheckIn(s string) bool {
	return s == "checkin" || s == "check-in" || s == "มาแล้ว" || s == "ทำงาน"
}

func (t *Tik) isGreeting(s string) bool {
	return s == "สวัสดี" || s == "hi" || s == "hello"
}

func (t *Tik) tokenize(s string, n int) []string {
	return strings.SplitN(s, " ", n)
}
