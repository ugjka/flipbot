package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	log "gopkg.in/inconshreveable/log15.v2"
)

//OSMmapResult ...
type OSMmapResult struct {
	Lat         string
	Lon         string
	DisplayName string `json:"Display_name"`
}

//OSMmapResults ...
type OSMmapResults []OSMmapResult

//OSMGeocode const
const OSMGeocode = "http://nominatim.openstreetmap.org/search?"

var osmCTR = struct {
	cache map[string][]byte
	*os.File
	sync.RWMutex
}{
	cache: make(map[string][]byte),
}

//OSMGetter Gets OSM DATA
func OSMGetter(url string) (data []byte, err error) {
	data, err = getOSMCache(url)
	switch {
	case err == errNotInCache:
		break
	case err != nil:
		return nil, err
	case err == nil:
		log.Info("osmgetter", "using", "cache")
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s (irc bot) for Freenode %s", ircNick, subreddit))
	get, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer get.Body.Close()
	if get.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status: %d", get.StatusCode)
	}
	data, err = ioutil.ReadAll(get.Body)
	if err != nil {
		return
	}
	if err := setOSMCache(url, data); err != nil {
		log.Warn("setOSMCache", "error", err)
	}
	return
}
