package helper

import (
	"assessment/apperror"
	"assessment/model"
	"encoding/json"
	"log"
	"net/http"
)

func ReturnFailure(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var err2 error
	switch err {
	case apperror.BadRequest:
		w.WriteHeader(apperror.BadRequest.Status)
		err2 = json.NewEncoder(w).Encode(map[string]string{"message": "invalid request received"})
	case apperror.NotFound:
		w.WriteHeader(apperror.NotFound.Status)
		err2 = json.NewEncoder(w).Encode(map[string]string{"message": "the requested resource was not found on this server"})
	case apperror.ServerError:
		w.WriteHeader(apperror.ServerError.Status)
		err2 = json.NewEncoder(w).Encode(map[string]string{"message": "an error occurred while processing that request"})

	}

	if err2 != nil {
		log.Printf("Error encoding JSON: %v", err2)
	}
}

func ReturnSuccess(w http.ResponseWriter, data model.Result) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := struct {
		Message string       `json:"message"`
		Data    model.Result `json:"result"`
	}{"success", data}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}
