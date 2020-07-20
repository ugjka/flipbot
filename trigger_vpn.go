package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var kickmeTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && m.Content == "!kickme"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Send(fmt.Sprintf("KICK %s %s :why are you kicking yourself", ircChannel, m.Name))
		return false
	},
}

var ipReg = regexp.MustCompile(`^\D*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})$`)
var vpnTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		if m.Command == "JOIN" {
			if !ipReg.MatchString(m.Host) {
				return false
			}
			if m.Name == "klimdaddie" || m.Name == "madk" {
				return false
			}
			if len(m.Params) == 3 && m.Params[1] != "*" {
				return false
			}
		}
		return m.Command == "JOIN"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		arr := ipReg.FindStringSubmatch(m.Host)
		ip := arr[1]
		vpn, err := whoisVPNCheck(ip)
		if err != nil {
			log.Error("whois vpn check", "error", err)
			return false
		}
		if vpn {
			log.Info("whois vpn detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :VPN detected, please identify before joining to bypass this check", ircChannel, m.Name))
			return false
		}
		vpn, err = torrentVPNCheck(ip)
		if err != nil {
			log.Error("torrent vpn check", "error", err)
			return false
		}
		if vpn {
			log.Info("torrent vpn detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :VPN detected, please identify before joining to bypass this check", ircChannel, m.Name))
			return false
		}
		return false
	},
}

var ipRangeReg = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

const torrentThresholdSize = 1024 * 50

func whoisVPNCheck(ip string) (vpn bool, err error) {
	cmd := exec.Command("whois", ip)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	data := string(b)
	if !ipRangeReg.MatchString(data) {
		return false, fmt.Errorf("no range found")
	}
	arr := ipRangeReg.FindStringSubmatch(data)
	start := strings.Split(arr[1], ".")
	end := strings.Split(arr[2], ".")
	for i := 0; i < 2; i++ {
		if start[i] != end[i] {
			return false, nil
		}
	}
	startInt, _ := strconv.ParseInt(start[2], 0, 64)
	endInt, _ := strconv.ParseInt(end[2], 0, 64)
	if endInt-startInt > 1 {
		return false, nil
	}
	return true, nil
}

func torrentVPNCheck(ip string) (vpn bool, err error) {
	res, err := httpClient.Get("https://iknowwhatyoudownload.com/en/peer/?ip=" + ip)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	if len(b) > torrentThresholdSize {
		return true, nil
	}
	return
}
