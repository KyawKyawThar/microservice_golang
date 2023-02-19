package main

import (
	"logger-serviice/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog is a simple test handler for the broker
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	//read json into var
	var requestPayload jsonPayload
	_ = app.readJSON(w, r, &requestPayload)

	//insert data

	event := data.LogEntry{

		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Model.LogEntry.Insert(event)

	if err != nil {
		err = app.errorJSON(w, err)
		return
	}
	resp := jsonResponse{

		Data:  "Logged",
		Error: false,
	}

	_ = app.writeJSON(w, http.StatusAccepted, resp)
}
