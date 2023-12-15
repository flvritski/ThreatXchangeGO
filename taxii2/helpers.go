package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/taxii+json;version=2.1")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) MarshalIndent(data interface{}) interface{} {
	out, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return nil
	}
	return out
}
