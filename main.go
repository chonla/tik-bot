package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chonla/oddsy"
	"github.com/chonla/tik-bot/tik"
)

var tk *tik.Tik

type configuration struct {
	oddsy.Configuration
	GcpToken          string `json:"gcp-token"`
	FirebaseProjectID string `json:"firebase-project-id"`
}

var conf configuration

func main() {
	loadConfig("./oddsy.json", &conf)

	b := oddsy.NewOddsy(&oddsy.Configuration{
		SlackToken:       conf.SlackToken,
		Debug:            conf.Debug,
		IgnoreBotMessage: conf.IgnoreBotMessage,
		BotInfo:          conf.BotInfo,
	})

	b.DirectMessageReceived(directMessageHandler)
	b.FirstStringTokenReceived("ping", pingMessageHandler)

	defer release()

	b.Start()
}

func release() {
	if tk != nil {
		tk.Release()
	}
}

func directMessageHandler(o *oddsy.Oddsy, m *oddsy.Message) {
	if tk == nil {
		tk = tik.NewTik(&tik.Configuration{
			GcpToken:          conf.GcpToken,
			FirebaseProjectID: conf.FirebaseProjectID,
		})
	}
	tk.Dispatch(o, m)
}

func pingMessageHandler(o *oddsy.Oddsy, m *oddsy.Message) {
	o.Send(m.Channel.UID, "pong :heart:")
}

func loadConfig(filename string, conf *configuration) {
	t, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	json.Unmarshal(t, conf)
}
