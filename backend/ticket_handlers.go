package main

import (
	log "boarding-pass/logging"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ticketRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DocumentNumber string `json:"documentNumber"`
}

const (
	defaultFlight = "OS123"
	defaultFrom   = "AMS"
	defaultTo     = "BCN"
	defaultSeat   = "12A"
)

func handleCreateTicket(w http.ResponseWriter, r *http.Request, state *ServerState) {
	var req ticketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid request", "failed to decode ticket request", err)
		return
	}

	first := strings.TrimSpace(req.FirstName)
	last := strings.TrimSpace(req.LastName)
	doc := strings.ToUpper(strings.TrimSpace(req.DocumentNumber))

	if first == "" || last == "" || doc == "" {
		respondWithErr(w, http.StatusBadRequest, "missing fields", "ticket request missing required fields", fmt.Errorf("empty values"))
		return
	}

	ticket := &Ticket{
		FirstName:      first,
		LastName:       last,
		DocumentNumber: doc,
		Flight:         defaultFlight,
		From:           defaultFrom,
		To:             defaultTo,
		Seat:           defaultSeat,
	}

	created := state.ticketStore.Create(ticket)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Error.Printf("failed to write ticket response: %v", err)
	}
}

func handleGetTicket(w http.ResponseWriter, r *http.Request, state *ServerState) {
	ticketID := mux.Vars(r)["ticketId"]
	if ticketID == "" {
		respondWithErr(w, http.StatusBadRequest, "missing ticketId", "ticketId path parameter required", fmt.Errorf("missing ticketId"))
		return
	}

	ticket, err := state.ticketStore.Get(ticketID)
	if err != nil {
		respondWithErr(w, http.StatusNotFound, "ticket not found", "failed to find ticket", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ticket); err != nil {
		log.Error.Printf("failed to write ticket response: %v", err)
	}
}
