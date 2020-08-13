package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var kickmeTrigger = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && m.Content == "!kickme"
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		irc.Send(fmt.Sprintf("KICK %s %s :why are you kicking yourself", ircChannel, m.Name))
	},
}

var ipReg = regexp.MustCompile(`^\D*(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})$`)
var vpnTrigger = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		if m.Command == "JOIN" {
			if !ipReg.MatchString(m.Host) {
				return false
			}
			if m.Name == "klimdaddie" || m.Name == "yousei" || m.Name == ircNick {
				return false
			}
			if len(m.Params) == 3 && m.Params[1] != "*" {
				return false
			}
		}
		return m.Command == "JOIN"
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		const warning = "VPN/Proxy/Datacenter IP addresses are banned, please identify with freenode before joining to bypass this check"
		arr := ipReg.FindStringSubmatch(m.Host)
		ip := arr[1]
		vpn, err := denyListVPNCheck(ip)
		if err != nil {
			log.Error("denylist vpn check", "error", err)
			return
		}
		if vpn {
			log.Info("denylist vpn detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :%s", ircChannel, m.Name, warning))
			return
		}
		vpn, err = subnetVPNCheck(ip)
		if err != nil {
			log.Error("subnet vpn check", "error", err)
			return
		}
		if vpn {
			log.Info("subnet vpn detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :%s", ircChannel, m.Name, warning))
			return
		}
		vpn, err = providerVPNCheck(ip)
		if err != nil {
			log.Error("provider vpn check", "error", err)
			return
		}
		if vpn {
			log.Info("provider vpn detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :%s", ircChannel, m.Name, warning))
			return
		}
	},
}

func subnetVPNCheck(ip string) (vpn bool, err error) {
	netmaskValid := false
	rangeValid := true
	var routeReg = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/(\d{1,2})`)
	var ipRangeReg = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	cmd := exec.Command("whois", ip)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	data := string(b)
	if routeReg.MatchString(data) {
		arr := routeReg.FindStringSubmatch(data)
		netmask, err := strconv.Atoi(arr[1])
		if err != nil {
			return false, err
		}
		if netmask < 22 {
			netmaskValid = true
		}
	}
	if !ipRangeReg.MatchString(data) {
		return false, fmt.Errorf("no range found")
	}
	res := ipRangeReg.FindAllStringSubmatch(data, -1)
	for _, arr := range res {
		start := strings.Split(arr[1], ".")
		end := strings.Split(arr[2], ".")
		for i := 0; i < 2; i++ {
			if start[i] != end[i] {
				return !(rangeValid || netmaskValid), nil
			}
		}
		startInt, _ := strconv.ParseInt(start[2], 0, 64)
		endInt, _ := strconv.ParseInt(end[2], 0, 64)
		if endInt-startInt > 1 {
			return !(rangeValid || netmaskValid), nil
		}
	}
	rangeValid = false
	return !(rangeValid || netmaskValid), nil
}

func providerVPNCheck(ip string) (vpn bool, err error) {
	var providerDenylist = []string{
		"abuse@m247.ro",
		"abuse@cdn77.com",
		"abuse@creanova.org",
		"abuse@estoxy.com",
		"abuse@panq.nl",
	}
	cmd := exec.Command("whois", ip)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	data := string(b)
	for _, v := range providerDenylist {
		if strings.Contains(data, v) {
			return true, nil
		}
	}
	return false, nil
}

var denyList = []string{}
var denyListOnce = &sync.Once{}

func denyListVPNCheck(ip string) (vpn bool, err error) {
	var denyLists = []string{
		"https://raw.githubusercontent.com/ejrv/VPNs/master/vpn-ipv4.txt",
		"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_level1.netset",
		"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_level2.netset",
		"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_abusers_30d.netset",
		"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/firehol_proxies.netset",
	}
	denyListOnce.Do(func() {
		for _, denyListURL := range denyLists {
			var res = &http.Response{}
			res, err = httpClient.Get(denyListURL)
			if err != nil {
				return
			}
			defer res.Body.Close()
			var data = []byte{}
			data, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}
			for _, v := range strings.Split(string(data), "\n") {
				v = strings.TrimSpace(v)
				if strings.HasPrefix(v, "#") || v == "" {
					continue
				}
				denyList = append(denyList, v)
			}
		}
		log.Info("denylists", "status", "loaded")
	})
	if err != nil {
		denyListOnce = &sync.Once{}
		return
	}
	for _, v := range denyList {
		if strings.Contains(v, "/") {
			_, subnet, err := net.ParseCIDR(v)
			if err != nil {
				return false, err
			}
			ipNet := net.ParseIP(ip)
			if subnet.Contains(ipNet) {
				return true, nil
			}
		} else {
			if ip == v {
				return true, nil
			}
		}
	}
	return
}

var denyBETrigger = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		if m.Command == "JOIN" {
			if m.Name == "klimdaddie" || m.Name == "yousei" || m.Name == ircNick {
				return false
			}
			if len(m.Params) == 3 && m.Params[1] != "*" {
				return false
			}
		}
		return m.Command == "JOIN"
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		const warning = "BELGIUM is banned, please identify with freenode before joining to bypass this check"
		ip := ""
		if ipReg.MatchString(m.Host) {
			arr := ipReg.FindStringSubmatch(m.Host)
			ip = arr[1]
		} else {
			ipRAW, err := net.ResolveIPAddr("ip4", m.Host)
			if err != nil {
				log.Error("deny BE", "can't resolve host", m.Host)
				return
			}
			ip = ipRAW.String()
		}
		be, err := denyListBECheck(ip)
		if err != nil {
			log.Error("denylist Belgium check", "error", err)
			return
		}
		if be {
			log.Info("denylist Belgium detected", "kicking", fmt.Sprintf("%s!%s@%s", m.Name, m.User, m.Host))
			irc.Send(fmt.Sprintf("REMOVE %s %s :%s", ircChannel, m.Name, warning))
			return
		}
	},
}

var denyBEList = []string{}
var denyBEListOnce = &sync.Once{}

func denyListBECheck(ip string) (BE bool, err error) {
	var denyLists = []string{
		// Ban Belgium
		"https://raw.githubusercontent.com/firehol/blocklist-ipsets/master/ipdeny_country/id_country_be.netset",
	}
	denyBEListOnce.Do(func() {
		for _, denyListURL := range denyLists {
			var res = &http.Response{}
			res, err = httpClient.Get(denyListURL)
			if err != nil {
				return
			}
			defer res.Body.Close()
			var data = []byte{}
			data, err = ioutil.ReadAll(res.Body)
			if err != nil {
				return
			}
			for _, v := range strings.Split(string(data), "\n") {
				v = strings.TrimSpace(v)
				if strings.HasPrefix(v, "#") || v == "" {
					continue
				}
				denyBEList = append(denyBEList, v)
			}
		}
		log.Info("denylist BE", "status", "loaded")
	})
	if err != nil {
		denyBEListOnce = &sync.Once{}
		return
	}
	for _, v := range denyBEList {
		if strings.Contains(v, "/") {
			_, subnet, err := net.ParseCIDR(v)
			if err != nil {
				return false, err
			}
			ipNet := net.ParseIP(ip)
			if subnet.Contains(ipNet) {
				return true, nil
			}
		} else {
			if ip == v {
				return true, nil
			}
		}
	}
	return
}
