package api

import (
	"encoding/json"
	"net/http"
)

// Error a struct to return on error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Message a struct to return on error
type Message struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, rsp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"code": 500, "message": "Could not write response"}`))
		return
	}
	w.WriteHeader(200)
	w.Write(response)
	return
}

func WriteJSONWithStatus(w http.ResponseWriter, rsp interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(rsp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"code": 500, "message": "Could not write response"}`))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
	return
}

// WriteOK writes the given interface as JSON to the given writer
func WriteOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"ok":true}`))
	return
}

// WriteJSONError writes the given error as JSON to the given writer
func WriteJSONError(w http.ResponseWriter, message string, internalCode string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(Error{
		Code:    internalCode,
		Message: message,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"code": 500, "message": "Could not write response"}`))
		return
	}
	w.WriteHeader(statusCode)
	w.Write(response)
	return
}

// WriteJSONError writes the given error as JSON to the given writer
func HandleRateLimiting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(Error{
		Code:    RATE_LIMIT_INTERNAL_CODE,
		Message: RATE_LIMIT_MESSAGE,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(`{"code": 500, "message": "Could not write response"}`))
		return
	}
	w.WriteHeader(RATE_LIMIT_ERROR_CODE)
	w.Write(response)
	return
}
