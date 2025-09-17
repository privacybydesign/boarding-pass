package main

type sessionResultPayload struct {
	Status      string `json:"status"`
	ProofStatus string `json:"proofStatus"`
	Err         *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Disclosed [][]*disclosedAttribute `json:"disclosed"`
}

type disclosedAttribute struct {
	RawValue   *string `json:"rawvalue"`
	Identifier string  `json:"id"`
}

func extractDocumentNumber(result *sessionResultPayload, expectedAttr string) (string, bool) {
	for _, conj := range result.Disclosed {
		for _, attr := range conj {
			if attr == nil || attr.RawValue == nil {
				continue
			}
			if attr.Identifier == expectedAttr {
				return *attr.RawValue, true
			}
		}
	}
	return "", false
}
