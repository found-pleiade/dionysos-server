package routes

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

type Stream struct {
	Users      map[uint64]bool
	ClientChan ClientChan
}

type Message struct {
	Event string
	Data  any
}

// New event messages are broadcast to all registered client connection channels
type MessageChan chan Message
type ClientChan map[uint64]MessageChan

func (s *Stream) create() {
	s.Users = make(map[uint64]bool)
	s.ClientChan = make(map[uint64]MessageChan)
}

func (s *Stream) distribute(m Message) {
	for _, clientChan := range s.ClientChan {
		clientChan <- m
	}
}

// addSub adds a user to stream if not already sub.
func (s *Stream) addSub(id uint64) error {
	if s.Users[id] {
		log.Printf("User %v already subscribed to stream", id)
		return errors.New("user already subscribed to stream")
	}

	s.ClientChan[id] = make(MessageChan)
	s.Users[id] = true
	log.Printf("added user %v to stream", id)
	log.Println("The stream users now looks like", s.Users)
	return nil
}

func (s *Stream) delSub(id uint64) error {
	if _, ok := s.Users[id]; ok {
		log.Printf("User %v has not subscribed to stream", id)
		return errors.New("user has not subscribed to stream")
	}
	close(s.ClientChan[id])
	s.Users[id] = false
	log.Printf("deleted user %v to stream", id)
	log.Println("The stream users now looks like", s.Users)
	return nil
}

func HeadersSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	// c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}
