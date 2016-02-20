package api

import (
	"encoding/json"
	"net/http"
)

type errorResp struct {
	Errors []JSONError `json:"errors"`
}

type JSONError struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

func sendError(w http.ResponseWriter, httpCode int, code string, msg string) {
	err := errorResp{
		Errors: []JSONError{
			{
				Code:  code,
				Title: msg,
			},
		},
	}
	w.WriteHeader(httpCode)
	enc := json.NewEncoder(w)
	enc.Encode(err)
}
