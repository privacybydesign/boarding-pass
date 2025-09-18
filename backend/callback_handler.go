package main

import (
	log "boarding-pass/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	irma "github.com/privacybydesign/irmago"
)

type issuanceJSON struct {
	Context     string `json:"@context"`
	CallbackURL string `json:"callbackURL,omitempty"`
	CallbackUrl string `json:"callbackUrl,omitempty"`
	*irma.IssuanceRequest
}

func handleIRMAServerCallback(w http.ResponseWriter, r *http.Request, state *ServerState) {
	ticketID := r.URL.Query().Get("ticketId")
	if ticketID == "" {
		respondWithErr(w, http.StatusBadRequest, "missing ticketId", "callback missing ticketId", fmt.Errorf("missing ticketId"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid body", "failed to read callback body", err)
		return
	}

	//is jwt so decode jwt first before making into json
	Parser := jwt.Parser{SkipClaimsValidation: true}
	parsedJWT, _, err := Parser.ParseUnverified(string(body), jwt.MapClaims{})
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid JWT", "failed to parse callback JWT", err)
		return
	}

	// Extract the claims as a JSON byte slice
	claims, ok := parsedJWT.Claims.(jwt.MapClaims)
	if !ok {
		respondWithErr(w, http.StatusBadRequest, "invalid JWT claims", "failed to extract JWT claims", fmt.Errorf("invalid JWT claims"))
		return
	}

	disclosed, _ := claims["disclosed"].([]interface{})
	inner, _ := disclosed[0].([]interface{})
	attr, _ := inner[0].(map[string]interface{})

	// Get the rawvalue
	disclosedDoc, ok := attr["rawvalue"].(string)
	if !ok {
		respondWithErr(w, http.StatusBadRequest, "missing document_number", "JWT claims missing rawvalue", fmt.Errorf("missing rawvalue"))
		return
	}

	ticket, err := state.ticketStore.Get(ticketID)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to load ticket for callback", err)
		return
	}
	if !strings.EqualFold(ticket.DocumentNumber, disclosedDoc) {
		w.WriteHeader(http.StatusOK)
		return
	}

	credID := irma.NewCredentialTypeIdentifier("irma-demo.demo-airline.boardingpass")
	cred := &irma.CredentialRequest{
		CredentialTypeID: credID,
		Attributes: map[string]string{
			"firstname": ticket.FirstName,
			"lastname":  ticket.LastName,
			"flight":    ticket.Flight,
			"from":      ticket.From,
			"to":        ticket.To,
			"seat":      ticket.Seat,
			"date":      ticket.Date,
			"time":      ticket.Time,
			"gate":      ticket.Gate,
		},
	}
	issuanceReq := irma.NewIssuanceRequest([]*irma.CredentialRequest{cred})
	payload := issuanceJSON{
		Context:         irma.LDContextIssuanceRequest,
		IssuanceRequest: issuanceReq,
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to marshal issuance payload", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bs); err != nil {
		log.Error.Printf("failed to write chained issuance response: %v", err)
	}
}
