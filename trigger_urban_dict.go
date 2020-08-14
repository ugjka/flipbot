package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var urbanTrig = regexp.MustCompile(`(?i)^\s*!(?:urban+|ud+)\w*\s+(\S.*)$`)
var urban = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && urbanTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		query := urbanTrig.FindStringSubmatch(m.Content)[1]
		query = strings.ToLower(query)
		defs, err := LookupWordDefinition(query)
		if err != nil {
			log.Warn("urban", "error", err)
			b.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		if len(defs.List) == 0 {
			b.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return
		}
		result := defs.List[0]
		result.Word = strings.ToLower(result.Word)
		if result.Word != query {
			b.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return
		}
		replacer := strings.NewReplacer("[", "", "]", "")
		result.Definition = replacer.Replace(result.Definition)
		result.Definition = strings.TrimSpace(result.Definition)
		b.Reply(m, fmt.Sprintf("%s: %s \n[%s]", m.Name, limit(result.Definition, 1024), result.Permalink))
	},
}

// UrbanResponse response from urban dictionary
type UrbanResponse struct {
	List []struct {
		Definition string
		Permalink  string
		Word       string
	}
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
