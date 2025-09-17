package main

import (
	log "boarding-pass/logging"
	"fmt"
	"net/http"
)

const ErrorInternal = "error:internal"

func respondWithErr(w http.ResponseWriter, code int, responseBody string, logMsg string, err error) {
	message := fmt.Sprintf("%v: %v", logMsg, err)
	log.Error.Printf("%s\n -> returning statuscode %d with message %v", message, code, responseBody)
	w.WriteHeader(code)
	if _, writeErr := w.Write([]byte(responseBody)); writeErr != nil {
		log.Error.Printf("failed to write body to http response: %v", writeErr)
	}
}
