package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, statusCode int, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}

type ErrorRespMsg struct {
	Message string `json:"message" example:"error message"`
	Code    string `json:"status,omitempty" example:"IR001"`
} // @name Error

func WriteError(rw http.ResponseWriter, statusCode int, err error) {
	_ = WriteJSON(rw, statusCode, ErrorRespMsg{
		Message: err.Error(),
	})
}

func DecodeJSON(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return json.NewDecoder(req.Body).Decode(obj)
}

func ParseQuery(req *http.Request, obj interface{}) error {
	params := req.URL.Query()

	m := map[string]string{}
	for k, v := range params {
		m[k] = v[0]
	}
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}
