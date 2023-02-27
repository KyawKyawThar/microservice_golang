package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	log.Println("Send Mail function call")

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		log.Println(err)
		_ = app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.sendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		_ = app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Message: "Send to" + requestPayload.To,

		Error: false,
	}

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		return
	}
}
