package main

import (
	kitty "bootybot/kittybot"
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"
)

var youtubedl = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && !kitty.Self(m) && xurls.Relaxed().MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		url := xurls.Relaxed().FindStringSubmatch(m.Content)[0]
		res, err := http.Get(url)
		if err == nil {
			content := res.Header.Get("Content-Type")
			if !strings.Contains(content, "text/html") {
				res.Body.Close()
				return
			}
			res.Body.Close()
		}
		video := ytdlOptions{
			url:           url,
			directory:     mp3Dir,
			server:        mp3Server,
			sizeLimit:     "100m",
			durationLimit: youtubeMaxDLDur,
		}
		link, err := video.Fetch()
		if err != nil {
			ytErrLog.Lock()
			ytErrLog.WriteString(time.Now().Format(time.RFC3339) + " | " + err.Error())
			ytErrLog.Unlock()
			return
		}
		desc := strings.TrimPrefix(link, fmt.Sprintf("https://%s/", mp3Server))
		if len(desc) > 45 {
			desc = desc[:45] + "..."
		}
		bot.ReplyRich(m, kitty.Rich{
			URL:         link,
			Title:       fmt.Sprintf("MP3 download of %s's link", m.Name),
			Description: desc,
			IconURL:     "https://raw.githubusercontent.com/ugjka/flipbot/bootybot/icon.png",
		})
	},
}

func ytdlParseDuration(format string) (time.Duration, error) {
	if strings.TrimSpace(format) == "" {
		return 0, nil
	}
	var units = []string{"s", "m", "h"}
	values := strings.Split(format, ":")
	format = ""
	for i, v := range values {
		format = v + units[len(values)-1-i] + format
	}
	return time.ParseDuration(format)
}

type ytdlOptions struct {
	url           string
	directory     string
	server        string
	sizeLimit     string
	durationLimit time.Duration
}

func (yt *ytdlOptions) Fetch() (string, error) {
	tmpfile := "/tmp/" + fmt.Sprintf("%x", sha1.Sum([]byte(yt.url)))
	if _, err := os.Stat(tmpfile); err == nil {
		return "", fmt.Errorf("video already downloading")
	}
	ioutil.WriteFile(tmpfile, []byte{}, 0777)
	defer os.Remove(tmpfile)
	options := []string{
		"--no-part",
		"--embed-thumbnail",
		"--add-metadata",
		"-x",
		"--audio-format=mp3",
		"--audio-quality=3",
		"--restrict-filenames",
		"--playlist-items=1",
		"--quiet",
		"--no-warnings",
		"--no-progress",
		"--match-filter=!is_live",
	}
	filename, err := ytdlFilename(yt.url)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(yt.directory + "/" + filename); err == nil {
		return "https://" + yt.server + "/" + filename, nil
	}
	dur, err := ytdlVideoDuration(yt.url)
	if err != nil {
		return "", err
	}
	if dur > yt.durationLimit {
		return "", fmt.Errorf("ytdl: video too long or short")
	}
	if dur == 0 {
		options = append(options, "--max-filesize="+yt.sizeLimit)
	}
	cmd := exec.Command("youtube-dl", append(options, yt.url)...)
	cmd.Dir = yt.directory + "/"
	errout := bytes.NewBuffer(nil)
	cmd.Stderr = errout
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("yt fetch: %v", errout.String())
	}
	return "https://" + yt.server + "/" + filename, nil
}

func ytdlVideoDuration(url string) (time.Duration, error) {
	options := []string{
		"--playlist-items=1",
		"--no-playlist",
		"--get-duration",
		"--quiet",
		"--no-warnings",
		"--no-progress",
		"--match-filter=!is_live",
	}
	cmd := exec.Command("youtube-dl", append(options, url)...)
	stdout := bytes.NewBuffer(nil)
	errout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = errout
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("ytdl duration: %s", errout.String())
	}
	return ytdlParseDuration(strings.TrimSpace(stdout.String()))
}

func ytdlFilename(url string) (string, error) {
	options := []string{
		"--restrict-filenames",
		"--playlist-items=1",
		"--no-playlist",
		"--get-filename",
		"--quiet",
		"--no-warnings",
		"--no-progress",
		"--match-filter=!is_live",
	}
	cmd := exec.Command("youtube-dl", append(options, url)...)
	stdout := bytes.NewBuffer(nil)
	errout := bytes.NewBuffer(nil)
	cmd.Stdout = stdout
	cmd.Stderr = errout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ytdl filename: %s", errout.String())
	}
	if strings.TrimSpace(stdout.String()) == "" {
		return "", fmt.Errorf("ytdl: no filename, live stream?")
	}
	dots := strings.Split(stdout.String(), ".")
	dots[len(dots)-1] = "mp3"
	return strings.Join(dots, "."), nil
}
