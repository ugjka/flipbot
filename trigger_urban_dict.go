package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

// UD trigger
var urban = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!urban ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		defs, err := LookupWordDefinition(m.Content[7:])
		if err != nil {
			log.Warn("could not get UD definition", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: error: %v", m.Name, err))
			return false
		}
		if len(defs.List) == 0 {
			irc.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return false
		}
		result := defs.List[0]
		replacer := strings.NewReplacer("\r", "", "\n", " ", "[", "", "]", "")
		result.Definition = replacer.Replace(result.Definition)
		if len(result.Definition) > 300 {
			irc.Reply(m, fmt.Sprintf("%s: %s... [%s]", m.Name, result.Definition[:300], result.Permalink))
		} else {
			irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, result.Definition))
		}
		return false
	},
}

// UrbanResponse response from urban dictionary
type UrbanResponse struct {
	List []struct {
		Definition string `json:"definition"`
		Permalink  string `json:"permalink"`
	} `json:"list"`
}

// LookupWordDefinition looks up definition
func LookupWordDefinition(word string) (urbanResponse UrbanResponse, err error) {
	resp, err := httpClient.Get("http://api.urbandictionary.com/v0/define?term=" + url.QueryEscape(word))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&urbanResponse); err != nil {
		return
	}
	return
}
