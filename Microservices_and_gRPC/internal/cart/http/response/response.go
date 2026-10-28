package response

import (
	"encoding/json"
	"net/http"
)

func Problem(w http.ResponseWriter, code int, desc error) {
	var problem struct {
		Code int    `json:"code"`
		Desc string `json:"description"`
	}
	problem.Code = code
	problem.Desc = desc.Error()
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(problem)
}

func Json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
