package app

import (
	"log"
	"net/http"
)

type message struct {
	UserId  string
	Channel string
	Payload string
}

// Internally transfers new client data to main routine
type clientData struct {
	session  *Session
	dataChan chan string
}

// A single Server will be created in this program. It is responsible
// for keeping a list of which clients (browsers) are currently attached
// and broadcasting events (messages) to those clients.
//
type Server struct {
	secret      string
	bind        string
	ctlbind     string
	clients     map[*Session]chan string
	newClients  chan *clientData
	deadClients chan *Session
	messages    chan *message
}

// This Server method starts a new goroutine.  It handles
// the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
//
func (b *Server) Start() {

	// Start a goroutine
	//
	go func() {

		// Loop endlessly
		//
		for {

			// Block until we receive from one of the
			// three following channels.
			select {

			case s := <-b.newClients:
				b.clients[s.session] = s.dataChan
				log.Println("Added new client")

			case s := <-b.deadClients:
				close(b.clients[s])
				delete(b.clients, s)
				log.Println("Removed client")

			case msg := <-b.messages:
				for s, c := range b.clients {
					if s.Channel == msg.Channel && (s.UserId == msg.UserId || msg.UserId == "*") {
						c <- msg.Payload
						log.Printf("Sent message to %s", s.UserId)
					}
				}

			}
		}
	}()

	publicServer := http.NewServeMux()
	publicServer.Handle("/events/", http.HandlerFunc(b.ServeHTTP))

	controlServer := http.NewServeMux()
	controlServer.Handle("/emit/", http.HandlerFunc(b.emitHandler))
	controlServer.Handle("/token/", http.HandlerFunc(b.tokenHandler))

	go func() {
		log.Fatal(http.ListenAndServe(b.ctlbind, controlServer))
	}()

	log.Fatal(http.ListenAndServe(b.bind, publicServer))

}

func NewServer(bind string, ctlbind string, secret string) *Server {
	// Make a new Server instance
	return &Server{
		secret,
		bind,
		ctlbind,
		make(map[*Session]chan string),
		make(chan *clientData),
		make(chan *Session),
		make(chan *message),
	}
}
