package racehub

import "sync"

type Hub struct {
	AssocID string
	clients map[string]*Client
	handler Handler
	action  Action
	Mu      sync.RWMutex
}

func NewHub(assocID string, handler Handler, action Action) *Hub {
	return &Hub{
		AssocID: assocID,
		clients: make(map[string]*Client),
		handler: handler,
		action:  action,
	}
}
