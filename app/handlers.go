package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (b *Server) emitHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/emit/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b.messages <- &message{
		r.FormValue("userId"),
		r.FormValue("channel"),
		r.FormValue("payload"),
	}

	w.Write([]byte("OK"))
}

// This Server method handles and HTTP request at the "/events/" URL.
//
func (b *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Make sure that the writer supports flushing.
	//
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	auth := r.Header.Get("Authorization")
	if len(auth) <= 6 || strings.ToLower(auth[:6]) != "bearer" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	session, err := getSession(auth[7:], b.secret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	dataChan := make(chan string)
	b.newClients <- &clientData{session: session, dataChan: dataChan}

	// Listen to the closing of the http connection via the CloseNotifier
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		// Remove this client from the map of attached clients
		// when `EventHandler` exits.
		b.deadClients <- session
		log.Println("HTTP connection just closed.")
	}()

	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Don't close the connection, instead loop endlessly.
	for {
		msg, open := <-dataChan
		if !open {
			break
		}
		fmt.Fprintln(w, msg)
		f.Flush()
	}

	// Done.
	log.Println("Finished HTTP request at ", r.URL.Path)
}

func (b *Server) tokenHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	channel := r.FormValue("channel")
	expires, err := strconv.ParseInt(r.FormValue("exp"), 10, 64)
	if err != nil {
		http.Error(w, "Incorrect expires", 500)
		return
	}
	claims := &Session{
		UserId:  userId,
		Channel: channel,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + expires,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(b.secret))
	w.Write([]byte(ss))
}
