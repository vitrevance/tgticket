package ticket

import (
	"sync"
	"time"
)

type Ticket struct {
	Token    string
	ExpireAt time.Time
}

type Store struct {
	tickets map[string]*Ticket
	mu      sync.Mutex
}

func NewStore() *Store {
	return &Store{
		tickets: make(map[string]*Ticket),
	}
}

func (s *Store) Add(ticket *Ticket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tickets[ticket.Token] = ticket
}

func (s *Store) Get(token string) (*Ticket, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tickets[token]
	return t, ok
}

func (s *Store) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tickets, token)
}

// Return all tickets which are still valid
func (s *Store) ActiveTickets() []*Ticket {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	var active []*Ticket
	for _, t := range s.tickets {
		if t.ExpireAt.After(now) {
			active = append(active, t)
		}
	}
	return active
}
