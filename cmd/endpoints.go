package main

import (
	"log"
	"net/http"
)

func generalHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("General handler hit")
	w.Write([]byte("<html><body><p>Hello</p></body></html>"))
}

func healthHandler(run *bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health endpoint hit")
		if *run {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Wait"))
		}
	}
}

func metricsHandler(metrics *Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(metrics.Print()))
		return
	}
}
