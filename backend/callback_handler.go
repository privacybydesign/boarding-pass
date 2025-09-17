package main

import (
	log "boarding-pass/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	var result sessionResultPayload
	if err := json.Unmarshal(body, &result); err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid result", "failed to parse callback JSON", err)
		return
	}

	if result.ProofStatus != string(irma.ProofStatusValid) || len(result.Disclosed) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	expectedAttr := fmt.Sprintf("pbdf-staging.%s.%s.%s", state.credentialConfig.IssuerId, state.credentialConfig.Credential, state.credentialConfig.Attribute)
	disclosedDoc, ok := extractDocumentNumber(&result, expectedAttr)
	if !ok || disclosedDoc == "" {
		w.WriteHeader(http.StatusOK)
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
