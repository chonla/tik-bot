package tik

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Command list
var cmdList = map[string][]string{
	"SUMMARY":  []string{"สรุป", "sum", "summary"},
	"CHECKIN":  []string{"checkin", "check-in", "เข้าทำงาน", "ลงชื่อ"},
	"GREETING": []string{"สวัสดี", "hello", "hi"},
	"HELP":     []string{"?", "งง", "help"},
}

type cmdKV struct {
	source  string
	useReg  bool
	matcher *regexp.Regexp
}

// Compiled command list
var cmdListRegex map[string][]*cmdKV

// englishCmdRegex
var englishCmdRegex *regexp.Regexp

func (t *Tik) isEnglishText(s string) bool {
	if englishCmdRegex == nil {
		englishCmdRegex = regexp.MustCompile("^[a-z\\-0-9_]+$")
	}
	return englishCmdRegex.MatchString(s)
}

func (t *Tik) compileCommands() {
	cmdListRegex = map[string][]*cmdKV{}
	for k, v := range cmdList {
		cmdListRegex[k] = []*cmdKV{}
		for i, n := 0, len(v); i < n; i++ {
			r := regexp.MustCompile(fmt.Sprintf("^(%s)\\s*(.*)\\s*$", v[i]))
			cmdListRegex[k] = append(cmdListRegex[k], &cmdKV{
				source:  v[i],
				useReg:  !t.isEnglishText(v[i]),
				matcher: r,
			})
		}
	}
}

func (t *Tik) tryParse(s, k string) ([]string, error) {
	if cmdListRegex == nil {
		return []string{}, errors.New("command not found")
	}
	if set, ok := cmdListRegex[k]; ok {
		for i, n := 0, len(set); i < n; i++ {
			if set[i].useReg {
				result := set[i].matcher.FindStringSubmatch(s)
				if len(result) > 0 {
					return result[1:], nil
				}
			} else {
				cmdTokens := t.tokenize(s, 2)
				if cmdTokens[0] == set[i].source {
					return cmdTokens, nil
				}
			}
		}
	}
	return []string{}, errors.New("command not found")
}

func (t *Tik) discover(s string) (string, []string) {
	for k := range cmdList {
		r, e := t.tryParse(s, k)
		if e == nil {
			return k, r
		}
	}
	return "UNKNOWN", []string{}
}

func (t *Tik) tokenize(s string, n int) []string {
	if strings.Contains(s, " ") {
		return strings.SplitN(s, " ", n)
	}
	return []string{s, ""}
}
