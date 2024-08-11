package racehub

import "sync"

type Hub struct {
	AssocID string
	clients map[string]*Client
	handler Handler
	Mu      sync.RWMutex
}

func NewHub(assocID string, handler Handler) *Hub {
	return &Hub{
		AssocID: assocID,
		clients: make(map[string]*Client),
		handler: handler,
	}
}
