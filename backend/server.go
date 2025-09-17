package main

import (
	log "boarding-pass/logging"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	irma "github.com/privacybydesign/irmago"
)

const ErrorInternal = "error:internal"

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Server struct {
	server *http.Server
	config *ServerConfig
}

type SpaHandler struct {
	staticPath string
	indexPath  string
}
type ServerState struct {
	irmaServerURL    string
	apiToken         string
	tokenStorage     TokenStorage
	credentialConfig CredentialConfig
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Join internally call path.Clean to prevent directory traversal
	path := filepath.Join(h.staticPath, r.URL.Path)
	// check whether a file exists or is a directory at the given path
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// file does not exist or path is a directory, serve index.html
		// Ensure index is never cached to avoid stale SPA after rebuilds
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func NewServer(state *ServerState, config *ServerConfig) *Server {
	router := mux.NewRouter()

	// Resolve the frontend dist path robustly regardless of working directory
	spa := SpaHandler{staticPath: "../frontend/dist", indexPath: "index.html"}
	// Register API routes before the SPA catch-all
	router.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		handleStartIRMASession(w, r, state)
	})
	router.HandleFunc("/api/result", func(w http.ResponseWriter, r *http.Request) {
		handleResultIRMASession(w, r, state)
	})

	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         config.Host + ":" + strconv.Itoa(config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return &Server{
		server: srv,
		config: config,
	}

}

func makeDisclosureRequest(credentialConfig *CredentialConfig) *irma.DisclosureRequest {

	irma.NewRequestorIdentifier("boarding-pass")

	disclosureRequest := irma.NewDisclosureRequest()

	// Request specific attributes
	disclosureRequest.Disclose = irma.AttributeConDisCon{
		irma.AttributeDisCon{
			irma.AttributeCon{irma.NewAttributeRequest("pbdf-staging." + credentialConfig.IssuerId + "." + credentialConfig.Credential + "." + credentialConfig.Attribute)},
		},
	}

	return disclosureRequest
}

// handlers -----------------------------

func handleStartIRMASession(w http.ResponseWriter, r *http.Request, state *ServerState) {

	if r.Method != http.MethodGet {
		respondWithErr(w, http.StatusMethodNotAllowed, "method not allowed", "invalid method", fmt.Errorf(r.Method))
		return
	}

	// Build a minimal disclosure request like the verifier tool
	disclosureReq := makeDisclosureRequest(&state.credentialConfig)

	disclosurePayload := disclosureJSON{
		Context:           "https://irma.app/ld/request/disclosure/v2",
		DisclosureRequest: disclosureReq,
	}

	sp, err := startSession(state.irmaServerURL, state.apiToken, disclosurePayload)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to start IRMA session", err)
		return
	}
	// Extract session ID from session pointer
	sessionID, err := extractSessionIDFromPtr(sp.SessionPtr)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to extract session ID from session pointer", err)
		return
	}

	// Store the session token for later verification
	if err := state.tokenStorage.StoreToken(sessionID, sp.Token); err != nil {
		respondWithErr(w, http.StatusInternalServerError, ErrorInternal, "failed to store session token", err)
		return
	}
	fmt.Println("stored token for sessionID " + sessionID)

	// Return the session pointer
	w.Header().Set("Content-Type", "application/json")
	type resp struct {
		SessionPtr json.RawMessage `json:"sessionPtr"`
		SessionID  string          `json:"sessionId"`
	}

	out, _ := json.Marshal(resp{SessionPtr: sp.SessionPtr})
	log.Info.Printf("started IRMA session with token %s", sp.SessionPtr)
	if _, err := w.Write(out); err != nil {
		log.Error.Printf("failed to write response: %v", err)
	}
}

func handleResultIRMASession(w http.ResponseWriter, r *http.Request, state *ServerState) {

	if r.Method != http.MethodGet {
		respondWithErr(w, http.StatusMethodNotAllowed, "method not allowed", "invalid method", fmt.Errorf(r.Method))
		return
	}

	// Get sessionID from query parameter
	sessionID := r.URL.Query().Get("sessionID")
	if sessionID == "" {
		respondWithErr(w, http.StatusBadRequest, "missing sessionID", "sessionID query parameter is required", fmt.Errorf("missing sessionID"))
		return
	}

	log.Info.Printf("Received sessionID: %s", sessionID)

	token, err := state.tokenStorage.RetrieveToken(sessionID)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "invalid sessionID", "failed to retrieve token for sessionID "+sessionID, err)
		return
	}

	log.Info.Printf("Retrieved token: %s for sessionID: %s", token, sessionID)

	requestorResultURL := fmt.Sprintf("%s/session/%s/result", state.irmaServerURL, token)

	// Send GET request to IRMA server
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

	var respBody json.RawMessage

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to read response body from IRMA server", err)
		return
	}
	if err := json.Unmarshal(b, &respBody); err != nil {
		respondWithErr(w, http.StatusBadGateway, ErrorInternal, "failed to unmarshal response body from IRMA server", err)
		return
	}

	// Forward the IRMA server response to the spa
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBody); err != nil {
		log.Error.Printf("failed to write response: %v", err)
	}

}

// helpers -----------------------------
func respondWithErr(w http.ResponseWriter, code int, responseBody string, logMsg string, e error) {
	m := fmt.Sprintf("%v: %v", logMsg, e)
	log.Error.Printf("%s\n -> returning statuscode %d with message %v", m, code, responseBody)
	w.WriteHeader(code)
	if _, err := w.Write([]byte(responseBody)); err != nil {
		log.Error.Printf("failed to write body to http response: %v", err)
	}
}

// --- IRMA/Yivi session helpers ---

type SessionPackage struct {
	Token      string          `json:"token"`
	SessionPtr json.RawMessage `json:"sessionPtr"`
}

// Inject JSON-LD @context while reusing irmago's request struct
type disclosureJSON struct {
	Context string `json:"@context"`
	*irma.DisclosureRequest
}

func startSession(baseURL string, token string, body any) (*SessionPackage, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/session", bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if token != "" {
		// Raw token (no "Bearer ") as used by Yivi staging server
		req.Header.Set("Authorization", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("start session: %d: %s", resp.StatusCode, string(rb))
	}
	var sp SessionPackage
	if err := json.Unmarshal(rb, &sp); err != nil {
		return nil, err
	}
	return &sp, nil
}

func extractSessionIDFromPtr(sessionPtr json.RawMessage) (string, error) {
	var ptrData struct {
		U string `json:"u"`
	}

	if err := json.Unmarshal(sessionPtr, &ptrData); err != nil {
		return "", fmt.Errorf("failed to unmarshal session pointer: %w", err)
	}

	// Extract session ID from URL like "https://my.app/irma/session/PR3M8VcShuVHKlSaAnAb"
	// Split by "/" and get the last part
	parts := strings.Split(ptrData.U, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid session URL format")
	}

	sessionID := parts[len(parts)-1]
	if sessionID == "" {
		return "", fmt.Errorf("empty session ID")
	}

	return sessionID, nil
}
