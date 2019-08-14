package main

import (
	"net/http"
)

func static(c Config) {
	for _, route := range c.Static {
		dir := http.Dir(route.Folder)
		s := http.FileServer(dir)
		headers := route.Headers
		handler := func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			s.ServeHTTP(w, r)
		}
		http.Handle(route.Route+"/", http.StripPrefix(route.Route, http.HandlerFunc(handler)))
	}
}
