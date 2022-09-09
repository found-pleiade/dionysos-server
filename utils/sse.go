package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Stream got a list of connected users and a channel per user to broadcast events.
type Stream struct {
	// Users is a map of users subscribed to the stream.
	// The key is the user ID and the value is a boolean indicating if the user is subscribed or not.
	Users map[uint64]bool

	// ClientChan is a map of channels to send messages to clients.
	ClientChan ClientChan
}

// Message represents a SSE type message.
type Message struct {
	// Event is the event type.
	Event string
	// Data is the data to send.
	Data any
}

// New event messages are broadcast to all registered client connection channels.
type MessageChan chan Message
type ClientChan map[uint64]MessageChan

// CreateStream creates a new stream and adds it to the list.
func CreateStream(ID uint64, list map[uint64]*Stream) error {
	list[ID] = &Stream{
		Users:      make(map[uint64]bool),
		ClientChan: make(ClientChan),
	}
	return nil
}

// GetStream returns an existing stream or error if it does not exist.
func GetStream(ID uint64, list map[uint64]*Stream) (*Stream, error) {
	stream, ok := list[ID]
	if !ok {
		return nil, errors.New("stream does not exist")
	}
	return stream, nil
}

func (s *Stream) Distribute(m Message) {
	for clientID, clientChan := range s.ClientChan {
		if s.Users[clientID] {
			clientChan <- m
		}
	}
}

// AddSub adds an ID to a stream if not already sub.
func (s *Stream) AddSub(id uint64) error {
	if s.Users[id] {
		return errors.New("user already subscribed to stream")
	}

	s.ClientChan[id] = make(MessageChan)
	s.Users[id] = true

	return nil
}

// DelSub removes an ID from a stream.
func (s *Stream) DelSub(id uint64) error {
	bool, ok := s.Users[id]
	if !ok {
		return errors.New("user has not subscribed to stream")
	} else if !bool {
		return errors.New("user has already been unsubscribed to stream")
	}

	close(s.ClientChan[id])
	s.Users[id] = false

	return nil
}

// HeaderSSE sets the regular headers for SSE at gin level.
func HeadersSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}
