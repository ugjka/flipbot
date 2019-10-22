package main

import (
	"encoding/json"
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
	osmCTR.Lock()
	defer osmCTR.Unlock()
	if v, ok := osmCTR.cache[url]; ok {
		return v, nil
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
	osmCTR.cache[url] = data
	tmp, err := json.Marshal(osmCTR.cache)
	if err == nil {
		err := osmCTR.Truncate(0)
		if err != nil {
			log.Crit("could not truncate the osmCacheFile", "error", err)
			return nil, err
		}
		if _, err := osmCTR.WriteAt(tmp, 0); err != nil {
			log.Crit("Could not write to osmCacheFile", "error", err)
			return nil, err
		}
	}
	return
}
