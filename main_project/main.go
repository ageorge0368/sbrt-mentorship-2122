package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	/// inport "strconv" and "github.com/gorilla/websocket"
	"strconv"

	"github.com/gorilla/websocket"
)

func main() {
	http.Handle(`/`, http.FileServer(http.Dir(`.`)))
	http.HandleFunc(`/login`, loginHandler)
	http.HandleFunc(`/chat`, chatHandler)
	go Broadcast()
	log.Println(`Server Initialized`)
	log.Fatal(http.ListenAndServeTLS(`:443`, `server.crt`, `server.key`, nil))

}

var (
	sessions   = make(map[string]*User)
	sessionMtx = sync.RWMutex{}
)

type User struct {
	sessionID int
	ws        SafeSocket
}

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
	sessionMtx.RLock()
	defer sessionMtx.RUnlock()
	if sessions[username] == nil {
		return 0
	}
	return sessions[username].sessionID
}

func AddSession(username string) int {
	sessionMtx.RLock()
	defer sessionMtx.Unlock()
	sessions[username] = &User{sessionID: rand.Intn(1_000_000) + 1}
	return sessions[username].sessionID
}

func RemoveSession(username string) {
	sessionMtx.Lock()
	defer sessionMtx.Unlock()
	sessions[username].ws.Close()
	delete(sessions, username)
}

type SafeSocket struct {
	*websocket.Conn
	*sync.RWMutex
}

type Message struct {
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

var (
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
	}
)

var broadcast = make(chan Message)

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if e := r.ParseForm(); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	name := strings.TrimSpace(r.FormValue(`username`))
	sessionID, e := strconv.Atoi(strings.TrimSpace(r.FormValue(`sessionID`)))
	if e != nil || name == `` {
		http.Error(w, `error: the sessionID does not match the inputted username`, http.StatusForbidden)
		return
	}
	if GetSession(name) != sessionID {
		http.Error(w, `error: sessionID does not match the given username`, http.StatusForbidden)
		return
	}
	ws, e := wsUpgrader.Upgrade(w, r, nil)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	sessionMtx.Lock()
	sessions[name].ws = SafeSocket{ws, &sync.RWMutex{}}
	sessionMtx.Unlock()
	defer RemoveSession(name)
	var message Message
	for {
		if e := ws.ReadJSON(&message); e != nil {
			log.Println(e)
			return
		}
		message.Sender = name
		broadcast <- message
	}
}

func Broadcast() {
	for message := range broadcast {
		sessionMtx.Lock()
		for _, user := range sessions {
			user.ws.Lock()
			user.ws.WriteJSON(message)
			user.ws.Unlock()
		}
		sessionMtx.Unlock()
	}
}
