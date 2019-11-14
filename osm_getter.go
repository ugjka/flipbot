package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/boltdb/bolt"
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

func setOSMCache(url string, data []byte) error {
	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("osmcache"))
		return b.Put([]byte(url), data)
	})
}

func getOSMCache(url string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("osmcache"))
		data := b.Get([]byte(url))
		if data == nil {
			return errNotInCache
		}
		_, err := buf.Write(data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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
