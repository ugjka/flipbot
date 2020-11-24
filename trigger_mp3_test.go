package main

import (
	"fmt"
	"testing"
	"time"
)

func TestParseYtDur(t *testing.T) {
	tt := []struct {
		input  string
		output time.Duration
	}{
		{"1", time.Second},
		{"1:12", time.Second*12 + time.Minute},
		{"2:12:1", time.Second + time.Minute*12 + time.Hour*2},
	}
	for _, tc := range tt {
		out, err := ytdlParseDuration(tc.input)
		if err != nil {
			t.Error(tc.input + " shouldn't error")
		}
		if out != tc.output {
			t.Error(fmt.Sprintf("%s is not %s", out, tc.output))
		}
	}
}

func TestYtdlFetch(t *testing.T) {
	yt := ytdlOptions{
		url:           "https://www.youtube.com/watch?v=HJOHoiPGpac",
		directory:     "tmp",
		server:        "mp3.ugjka.net",
		sizeLimit:     "100m",
		durationLimit: youtubeMaxDLDur,
	}
	yt.Fetch()
}

func TestYtdlVideoDuration(t *testing.T) {
	tt := []struct {
		url      string
		duration time.Duration
		err      bool
	}{
		{
			url:      "https://www.youtube.com/watch?v=HJOHoiPGpac",
			duration: time.Minute * 3,
			err:      false,
		},
		{
			url:      "example.org",
			duration: 0,
			err:      true,
		},
		{
			url:      "https://www.instagram.com/p/CH0okSVA1tj/?utm_source=ig_web_button_share_sheet",
			duration: 0,
			err:      false,
		},
		{
			url:      "https://www.youtube.com/watch?v=dfGhUaOBVSs&list=PL5QsCBZgatCDl4_YVlAF5BuGN6WTszWgz",
			duration: time.Hour + time.Minute*54 + time.Second*28,
			err:      false,
		},
	}
	for _, tc := range tt {
		dur, err := ytdlVideoDuration(tc.url)
		if !tc.err && err != nil {
			t.Errorf("unexpected error for %s: %v", tc.url, err)
			continue
		}
		if dur != tc.duration {
			t.Errorf("expected %s, got %s", tc.duration, dur)
		}
	}
}

func TestGetYTDLFilename(t *testing.T) {
	tt := []struct {
		url      string
		filename string
		err      bool
	}{
		{
			url:      "https://www.youtube.com/watch?v=HJOHoiPGpac",
			filename: "Tegan_and_Sara_-_Boyfriend_OFFICIAL_MUSIC_VIDEO-HJOHoiPGpac.mp3",
			err:      false,
		},
		{
			url:      "example.org",
			filename: "",
			err:      true,
		},
		{
			url:      "https://www.instagram.com/p/CH0okSVA1tj/?utm_source=ig_web_button_share_sheet",
			filename: "Video_by_earlewrites-CH0okSVA1tj.mp3",
			err:      false,
		},
		{
			url:      "https://www.youtube.com/watch?v=dfGhUaOBVSs&list=PL5QsCBZgatCDl4_YVlAF5BuGN6WTszWgz",
			filename: "2017_05_27Temecula_dialogue_0-dfGhUaOBVSs.mp3",
			err:      false,
		},
		{
			url:      "https://www.youtube.com/watch?v=ZNC2rm313vY",
			filename: "",
			err:      true,
		},
	}
	for _, tc := range tt {
		f, err := ytdlFilename(tc.url)
		if !tc.err && err != nil {
			t.Errorf("unexpected error for %s: %v", tc.url, err)
			continue
		}
		if f != tc.filename {
			t.Errorf("expected %s, got %s", tc.filename, f)
		}
	}
}
