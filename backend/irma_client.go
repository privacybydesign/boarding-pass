package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	irma "github.com/privacybydesign/irmago"
)

type SessionPackage struct {
	Token      string          `json:"token"`
	SessionPtr json.RawMessage `json:"sessionPtr"`
}

func makeDisclosureRequest(credentialConfig *CredentialConfig) *irma.DisclosureRequest {
	irma.NewRequestorIdentifier("boarding-pass")

	disclosureRequest := irma.NewDisclosureRequest()
	disclosureRequest.Disclose = irma.AttributeConDisCon{
		irma.AttributeDisCon{
			irma.AttributeCon{irma.NewAttributeRequest("pbdf-staging." + credentialConfig.IssuerId + "." + credentialConfig.Credential + "." + credentialConfig.Attribute)},
		},
	}
	return disclosureRequest
}

func startChainedSession(baseURL string, req irma.ServiceProviderRequest, requestorID string, priv *rsa.PrivateKey) (*SessionPackage, error) {
	jwtStr, err := irma.SignRequestorRequest(&req, jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()), priv, requestorID)
	if err != nil {
		return nil, fmt.Errorf("sign requestor request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, baseURL+"/session", bytes.NewBufferString(jwtStr))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "text/plain")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("start chained session: %d: %s", resp.StatusCode, string(rb))
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
