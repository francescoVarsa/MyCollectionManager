package throttler

import (
	"sync"
	"time"
)

type client struct {
	LastRequest  time.Time
	RequestCount int
}

type ThrottleRequests struct {
	clients  map[string]*client
	limit    int
	duration time.Duration
}

var mu sync.Mutex // Mutex per proteggere l'accesso concorrente alla mappa

func (t *ThrottleRequests) Validate(ip string) bool {
	// Acquisire il lock per accedere in sicurezza alla mappa
	mu.Lock()
	defer mu.Unlock()

	if t.clients == nil {
		t.clients = make(map[string]*client)
	}

	// Ottieni il client o creane uno nuovo se non esiste
	clnt, exists := t.clients[ip]
	if !exists {
		clnt = &client{
			LastRequest:  time.Now(),
			RequestCount: 0,
		}
		t.clients[ip] = clnt
	}

	// Controlla il tempo dalla sua ultima richiesta
	if time.Since(clnt.LastRequest) > t.duration {
		clnt.RequestCount = 0 // Reset del contatore se è passato abbastanza tempo
	}

	clnt.LastRequest = time.Now()
	clnt.RequestCount++

	// Se ha superato il limite, blocca la richiesta
	return clnt.RequestCount < t.limit
}

func (t *ThrottleRequests) SetLimit(limit int) {
	t.limit = limit
}

func (t *ThrottleRequests) SetDuration(duration time.Duration) {
	t.duration = duration
}

func (t *ThrottleRequests) ClientsCleanup(duration time.Duration) {
	for {
		time.Sleep(duration)

		mu.Lock()

		for key, client := range t.clients {
			if time.Since(client.LastRequest) > (duration * 2) {
				delete(t.clients, key)
			}
		}

		mu.Unlock()
	}
}
