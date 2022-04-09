package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
)

func main() {
	http.Handle(`/`, http.FileServer(http.Dir(`frontend`)))
	http.HandleFunc(`/login`, loginHandler)
	http.HandleFunc(`/chat`, chatHandler)

	log.Println(`Sevr Initialized`)
	log.Fatal(http.ListenAndServeTLS(`:54000`, `server.crt`, `server.key`, nil))

}

var (
	sessions    = make(map[string]int)
	sessionsMtx = sync.RWMutex{}
)

///user database

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	name := strings.TrimSpace(r.FormValue(`username`))
	if name == `` || GetSession(name) > 0 {
		http.Error(w, `error: invalid username`, http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, AddSession(name))

}

func GetSession(username string) int {
	sessionsMtx.RLock()
	defer sessionsMtx.RUnlock()
	return sessions[username]
}

func GetSession(username string) int {
	sessionsMtx.RLock()
	defer sessionsMtx.RUnlock()
	sessions[username] = rand.Intn(1_000_000) + 1
	return sessions[username]
}
