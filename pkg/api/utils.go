package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ExtractBody(r *http.Request, to interface{}) error {
	// bad request body.
	body, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		return readErr
	}

	// Extract the request.
	unmarshErr := json.Unmarshal(body, to)
	if unmarshErr != nil {
		return unmarshErr
	}
	return nil
}

type errorObject struct {
	error string `json:"error"`
}

func Error(w http.ResponseWriter, error string, code int) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(&errorObject{
		error: error,
	})
}
