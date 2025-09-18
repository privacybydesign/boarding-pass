package main

import (
	log "boarding-pass/logging"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	irma "github.com/privacybydesign/irmago"
)

type startVerificationRequest struct {
	TicketID       string `json:"ticketId"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DocumentNumber string `json:"documentNumber"`
}

type verificationResponse struct {
	SessionResult json.RawMessage `json:"sessionResult"`
	Verified      bool            `json:"verified,omitempty"`
	Message       string          `json:"message,omitempty"`
	Issuance      *sessionPointer `json:"issuance,omitempty"`
}

type sessionPointer struct {
	SessionPtr json.RawMessage `json:"sessionPtr"`
	SessionID  string          `json:"sessionId"`
}

// sessionResultPayload models the subset of the IRMA server result we use
type sessionResultPayload struct {
	Status      string                 `json:"status"`
	ProofStatus string                 `json:"proofStatus"`
	Disclosed   [][]disclosedAttribute `json:"disclosed"`
	Err         *struct {
		Message string `json:"message"`
	} `json:"err,omitempty"`
}

type disclosedAttribute struct {
	ID string `json:"id"`
	// IRMA uses "rawvalue" (all lowercase); be tolerant to alternate casing
	RawValue  string      `json:"rawvalue"`
	RawValue2 string      `json:"rawValue"`
	Value     interface{} `json:"value"`
}

func extractDocumentNumber(res *sessionResultPayload, expectedAttr string) (string, bool) {
	for _, group := range res.Disclosed {
		for _, attr := range group {
			if strings.EqualFold(attr.ID, expectedAttr) {
				v := attr.RawValue
				if v == "" {
					v = attr.RawValue2
				}
				if v == "" {
					if s, ok := attr.Value.(string); ok {
						v = s
					}
				}
				v = strings.ToUpper(strings.TrimSpace(v))
				if v != "" {
					return v, true
				}
			}
		}
	}
	return "", false
}

func handleStartIRMASession(w http.ResponseWriter, r *http.Request, state *ServerState) {
	var req startVerificationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid request", "failed to decode start request", err)
		return
	}

	ticketID := strings.TrimSpace(req.TicketID)
	if ticketID == "" {
		respondWithErr(w, http.StatusBadRequest, "missing ticketId", "ticketId is required", fmt.Errorf("empty ticketId"))
		return
	}

	first := strings.TrimSpace(req.FirstName)
	last := strings.TrimSpace(req.LastName)
	doc := strings.ToUpper(strings.TrimSpace(req.DocumentNumber))
	if first == "" || last == "" || doc == "" {
		respondWithErr(w, http.StatusBadRequest, "missing fields", "start request missing required fields", fmt.Errorf("empty values"))
		return
	}

	ticket, err := state.ticketStore.Get(ticketID)
	if err != nil {
		respondWithErr(w, http.StatusNotFound, "ticket not found", "failed to load ticket", err)
		return
	}

	if !strings.EqualFold(ticket.FirstName, first) || !strings.EqualFold(ticket.LastName, last) || !strings.EqualFold(ticket.DocumentNumber, doc) {
		respondWithErr(w, http.StatusBadRequest, "ticket mismatch", "provided passenger data does not match ticket", fmt.Errorf("ticket mismatch"))
		return
	}

	disclosureReq := makeDisclosureRequest(&state.credentialConfig)
	cbURL := buildCallbackURL(r, ticket.ID)
	ext := irma.ServiceProviderRequest{
		RequestorBaseRequest: irma.RequestorBaseRequest{
			NextSession: &irma.NextSessionData{URL: cbURL},
		},
		Request: disclosureReq,
	}

	keyBytes, err := os.ReadFile(state.credentialConfig.PrivateKeyPath)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to read private key", err)
		return
	}

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to parse private key", err)
		return
	}

	sp, err := startChainedSession(state.irmaServerURL, ext, state.credentialConfig.RequestorId, priv)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to start IRMA session", err)
		return
	}

	sessionID, err := extractSessionIDFromPtr(sp.SessionPtr)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to extract session ID", err)
		return
	}

	if err := state.tokenStorage.StoreToken(sessionID, sp.Token); err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to store session token", err)
		return
	}
	state.sessionTracker.Link(sessionID, ticket.ID)
	log.Info.Printf("started IRMA disclosure (chained) session for ticket %s", ticket.ID)

	type resp struct {
		SessionPtr json.RawMessage `json:"sessionPtr"`
		SessionID  string          `json:"sessionId"`
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp{SessionPtr: sp.SessionPtr, SessionID: sessionID}); err != nil {
		log.Error.Printf("failed to write response: %v", err)
	}
}

func handleResultIRMASession(w http.ResponseWriter, r *http.Request, state *ServerState) {
	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		respondWithErr(w, http.StatusBadRequest, "missing sessionID", "sessionID query parameter is required", fmt.Errorf("missing sessionID"))
		return
	}

	token, err := state.tokenStorage.RetrieveToken(sessionID)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid sessionID", "failed to retrieve token", err)
		return
	}

	requestorResultURL := fmt.Sprintf("%s/session/%s/result", state.irmaServerURL, token)
	getResultReq, err := http.NewRequest(http.MethodGet, requestorResultURL, nil)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to create get request to IRMA server", err)
		return
	}
	getResultReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(getResultReq)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to send get request to IRMA server", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to read response body from IRMA server", err)
		return
	}

	response := verificationResponse{SessionResult: json.RawMessage(body)}

	if ticketID, ok := state.sessionTracker.TicketID(sessionID); ok {
		state.sessionTracker.Remove(sessionID)
		var result sessionResultPayload
		if err := json.Unmarshal(body, &result); err != nil {
			respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to parse IRMA result", err)
			return
		}

		if result.Status == string(irma.ServerStatusDone) && result.ProofStatus == string(irma.ProofStatusValid) && result.Err == nil {
			expectedAttr := fmt.Sprintf("pbdf-staging.%s.%s.%s", state.credentialConfig.IssuerId, state.credentialConfig.Credential, state.credentialConfig.Attribute)
			disclosedDoc, ok := extractDocumentNumber(&result, expectedAttr)
			if !ok {
				response.Message = "required attribute not disclosed"
			} else {
				ticket, err := state.ticketStore.Get(ticketID)
				if err != nil {
					respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to load ticket for verification", err)
					return
				}

				if !strings.EqualFold(ticket.DocumentNumber, disclosedDoc) {
					response.Message = "passport data does not match ticket"
				} else {
					response.Verified = true
				}
			}
		} else if result.Err != nil {
			response.Message = result.Err.Message
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error.Printf("failed to write result response: %v", err)
	}
}
