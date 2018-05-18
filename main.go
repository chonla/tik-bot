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
	SlackToken        string `json:"slack-token"`
	Debug             bool   `json:"debug"`
	IgnoreBotMessage  bool   `json:"ignore-bot-message"`
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
	})

	// b.MessageReceived(messageHandler)
	b.DirectMessageReceived(directMessageHandler)
	b.FirstStringTokenReceived("help", helpMessageHandler)
	b.FirstStringTokenReceived("ping", pingMessageHandler)
	// b.FirstStringTokenReceived("tik", tikMessageHandler)

	defer release()

	b.Start()
}

func release() {
	if tk != nil {
		tk.Release()
	}
}

func messageHandler(o *oddsy.Oddsy, m *oddsy.Message) {
	if m.Mentioned {
		o.Send(m.Channel.UID, "<@"+m.From.UID+"> เรียกเค้าทำไมจ๊ะ คิดถึงล่ะสิ")
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

func helpMessageHandler(o *oddsy.Oddsy, m *oddsy.Message) {
	o.Send(m.Channel.UID, "ข้อความที่"+o.Name+" เข้าใจนะจ๊ะ\n```"+`
ping - ทดสอบ ping/pong
help - ข้อความนี้แหละจ้ะ
tik - ลง worksheet`+"```")
}

func tikMessageHandler(o *oddsy.Oddsy, m *oddsy.Message) {
	if tk == nil {
		tk = tik.NewTik(&tik.Configuration{
			GcpToken:          conf.GcpToken,
			FirebaseProjectID: conf.FirebaseProjectID,
		})
	}

	tk.Dispatch(o, m)
}

func loadConfig(filename string, conf *configuration) {
	t, e := ioutil.ReadFile(filename)
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}

	json.Unmarshal(t, conf)
}
