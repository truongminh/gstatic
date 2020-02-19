package main

import (
	"fmt"
	"gstatic/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxResize = 2048
)

func detectContentType(filename string) (contentType string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	buf := make([]byte, 512)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		return
	}
	contentType = http.DetectContentType(buf)
	return
}

func decodeImage(filename string) (img image.Image, err error) {
	contentType, err := detectContentType(filename)
	if err != nil {
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	default:
		err = fmt.Errorf("unknown content type %s", contentType)
	}
	return
}

func resizeImage(out io.Writer, filename string, params string) (err error) {
	img, err := decodeImage(filename)
	if err != nil {
		return
	}
	var width, height uint
	_, err = fmt.Sscanf(params, "%dx%d", &width, &height)
	if err != nil {
		return
	}
	if width > maxResize || height > maxResize {
		err = fmt.Errorf("%dx%d is not in %dx%d", width, height, maxResize, maxResize)
		return
	}
	m := resize.Resize(width, height, img, resize.Lanczos3)
	err = jpeg.Encode(out, m, nil)
	return
}

func (route *FileRoute) transform(w http.ResponseWriter, r *http.Request) {
	ts := r.URL.Query().Get("transform")
	parts := strings.Split(ts, ":")
	if len(parts) < 2 {
		http.Error(w, "transform must be <code>:<params>", http.StatusBadRequest)
		return
	}
	w.Header().Set("X-Transform", ts)
	code := parts[0]
	params := parts[1]
	file := filepath.Join(route.Folder, r.URL.Path)
	switch code {
	case "resize":
		resizeImage(w, file, params)
	}
}
