package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type db struct {
	records [][]string
	mx      sync.RWMutex
}

func (d *db) auth(username, password string) bool {
	d.mx.RLock()
	defer d.mx.RUnlock()
	for _, rec := range d.records {
		if rec[0] == username && rec[1] == password {
			return true
		}
	}
	return false
}

func main() {
	db := &db{
		records: [][]string{
			{"johndoe", "pass1"},
		},
	}

	addr := os.Getenv("AUTH_ADDR")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("AUTH"))
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		if db.auth(username, password) {
			fmt.Printf("user '%s' provided correct credentials\n", username)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("AUTHORIZED"))
		} else {
			fmt.Printf("user '%s' provided incorrect credentials\n", username)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("NOT AUTHORIZED"))
		}
	})

	log.Printf("Starting server on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
