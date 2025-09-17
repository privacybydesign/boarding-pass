package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type Ticket struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	DocumentNumber string    `json:"documentNumber"`
	Flight         string    `json:"flight"`
	From           string    `json:"from"`
	To             string    `json:"to"`
	Seat           string    `json:"seat"`
	Date           string    `json:"date"`
	Time           string    `json:"time"`
	Gate           string    `json:"gate"`
	CreatedAt      time.Time `json:"createdAt"`
}

type TicketStore struct {
	mu      sync.RWMutex
	tickets map[string]*Ticket
}

func NewTicketStore() *TicketStore {
	return &TicketStore{tickets: make(map[string]*Ticket)}
}

func (s *TicketStore) Create(ticket *Ticket) *Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()

	cloned := *ticket
	cloned.ID = newID()
	cloned.CreatedAt = time.Now().UTC()
	s.tickets[cloned.ID] = &cloned
	return &cloned
}

func (s *TicketStore) Get(id string) (*Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if ticket, ok := s.tickets[id]; ok {
		cloned := *ticket
		return &cloned, nil
	}
	return nil, errors.New("ticket not found")
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

type SessionTracker struct {
	mu    sync.Mutex
	links map[string]string
}

func NewSessionTracker() *SessionTracker {
	return &SessionTracker{links: make(map[string]string)}
}

func (s *SessionTracker) Link(sessionID, ticketID string) {
	s.mu.Lock()
	s.links[sessionID] = ticketID
	s.mu.Unlock()
}

func (s *SessionTracker) TicketID(sessionID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ticketID, ok := s.links[sessionID]
	return ticketID, ok
}

func (s *SessionTracker) Remove(sessionID string) {
	s.mu.Lock()
	delete(s.links, sessionID)
	s.mu.Unlock()
}
