package main

import (
	"net/http"
)

func static(c Config) {
	for _, route := range c.Static {
		route.register()
	}
}

func (route *FileRoute) register() {
	dir := http.Dir(route.Folder)
	fs := http.FileServer(dir)
	headers := route.Headers
	handler := func(w http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		transform := r.URL.Query().Get("transform")
		if len(transform) > 0 {
			route.transform(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	}
	http.Handle(route.Route+"/", http.StripPrefix(route.Route, http.HandlerFunc(handler)))
}
