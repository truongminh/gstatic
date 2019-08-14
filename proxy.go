package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/logger"
)

func proxy(c Config) {
	folder := c.Proxy.Folder
	header := "X-GSTATIC-PROXY"
	linkToName := func(link string) string {
		ext := filepath.Ext(link)
		buf := md5.Sum([]byte(link))
		sum := hex.EncodeToString(buf[:])
		dir := filepath.Join(folder, sum[:2], sum[2:4])
		return filepath.Join(dir, sum+ext)
	}
	setHeaders := func(w http.ResponseWriter) {
		for k, v := range c.Proxy.Headers {
			w.Header().Set(k, v)
		}
	}
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(header) != "" {
			http.Error(w, "loopback is not allowed", http.StatusBadRequest)
			return
		}
		link := r.URL.Query().Get("link")
		if link == "" {
			http.Error(w, "missing link", http.StatusBadRequest)
			return
		}

		setHeaders(w)
		name := linkToName(link)
		_, err := os.Stat(name)
		if err == nil {
			http.ServeFile(w, r, name)
			return
		}
		if !os.IsNotExist(err) {
			http.Error(w, "read file error "+err.Error(), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest(http.MethodGet, link, nil)
		if err != nil {
			http.Error(w, "invalida link "+err.Error(), http.StatusBadRequest)
			return
		}
		req.Header.Set(header, "ON")
		// Download the content
		res, err := http.DefaultClient.Do(req.WithContext(r.Context()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			http.Error(w, "upstream status: "+res.Status, http.StatusInternalServerError)
			return
		}
		if err := os.MkdirAll(filepath.Dir(name), os.ModePerm); err != nil {
			logger.Errorf("mkdir %s error: %s", filepath.Dir(name), err.Error())
			http.Error(w, "dir error", http.StatusInternalServerError)
			return
		}
		file, err := ioutil.TempFile("/tmp", "proxy")
		if err != nil {
			logger.Errorf("create temp file %s error %s", link, err.Error())
			http.Error(w, "create temp file failed", http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(file, res.Body)
		if err != nil {
			logger.Errorf("download file %s error %s", link, err.Error())
			http.Error(w, "download file error", http.StatusInternalServerError)
			return
		}
		err = file.Close()
		if err != nil {
			logger.Errorf("close temp file %s error %s", link, err.Error())
			http.Error(w, "close temp file error", http.StatusInternalServerError)
			return
		}
		err = os.Rename(file.Name(), name)
		if err != nil {
			logger.Errorf("remove cache %s error %s", link, err.Error())
			http.Error(w, "local cache rename error ", http.StatusInternalServerError)
			return
		}
		logger.Infof("cache %s to %s", link, name)
		http.ServeFile(w, r, name)
	})
}
