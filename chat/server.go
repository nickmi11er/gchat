package chat

import (
	"../documents"
	"golang.org/x/net/websocket"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type Server struct {
	pattern   string
	messages  []*documents.Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *documents.Message
	doneCh    chan bool
	errCh     chan error
}

func NewServer(pattern string) *Server {
	var messages []*documents.Message
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *documents.Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *documents.Message) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *documents.Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

func (s *Server) Listen() {
	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		token := ws.Request().FormValue("session_token")
		var user documents.UserDocument
		err := documents.Connection().UsersCollection.Find(bson.M{"session_token": token}).One(&user)
		if err != nil {
			log.Println("Invalid session token")
			ws.Write([]byte(`{"error":"invalid token"}`))
			return
		}

		client := NewClient(ws, s, user.Username)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

			// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

			// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
