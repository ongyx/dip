package main

import (
	"log"
	"net/http"
)

type LogHandler struct {
	log     *log.Logger
	handler http.Handler
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.handler.ServeHTTP(w, r)

	l.log.Printf(`%s "%s %s %s"`, r.RemoteAddr, r.Method, r.URL.Path, r.Proto)
}
