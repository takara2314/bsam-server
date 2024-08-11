package racehub

import "sync"

type Hub struct {
	AssociationID string
	clients       map[string]*Client
	handler       Handler
	action        Action
	Mu            sync.RWMutex
}

func NewHub(associationID string, handler Handler, action Action) *Hub {
	return &Hub{
		AssociationID: associationID,
		clients:       make(map[string]*Client),
		handler:       handler,
		action:        action,
	}
}
