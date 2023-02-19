package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

//type LogPayload struct {}

// Broker is a simple test handler for the broker
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Message: "Hits the Broker Service",
		Error:   false,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// HandleSubmit handles a JSON payload that describes an action to take,
// processes it, and sends it where it needs to go
func (app *Config) HandleSubmit(w http.ResponseWriter, r *http.Request) {

	//read, _ := io.ReadAll(r.Body)
	//var auth AuthPayload
	//json.Unmarshal(read, &auth)
	//
	//log.Printf("new Data %T", auth.Password)

	var requestedPayload RequestPayload

	//
	err := app.readJSON(w, r, &requestedPayload)

	if err != nil {
		err = app.errorJSON(w, err)
		return
	}
	//
	requestedPayload.Auth = AuthPayload{
		Email:    "admin@example.com",
		Password: "verysecret",
	}

	requestedPayload.Log = LogPayload{
		Name: "event",
		Data: "Some kind of data",
	}

	switch requestedPayload.Action {
	case "auth":
		app.authenticate(w, requestedPayload.Auth)
	case "log":
		app.logItem(w, requestedPayload.Log)
	default:
		_ = app.errorJSON(w, errors.New("unknown action"))
	}
}

// logItem logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {

	jsonData, _ := json.MarshalIndent(l, "", "\t")

	logServiceURL := fmt.Sprintf("http://%s/log", "logger-service")

	// convert byte array to string
	//log.Println("logItem", fmt.Sprintf("%s", jsonData))

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	//
	if err != nil {
		log.Println("err func call", err)
		err = app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-type", "application/json")
	client := http.Client{}

	response, err := client.Do(request)
	defer response.Body.Close()
	log.Println("client", response)
	//after response call looger service route
	if response.StatusCode != http.StatusAccepted {
		err = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

// authenticate tries to log a user in through the authentication-service. It receives a json payload
// of type requestPayload, with AuthPayload embedded.
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {

	//{admin@example.com verysecret}

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	//call the service
	//authServiceURL := fmt.Sprintf("http://%s/authenticated", "127.0.0.1:81")
	testServiceURL := fmt.Sprintf("http://%s/authenticated", "auth-service")

	request, err := http.NewRequest("POST", testServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("err func call")
		err = app.errorJSON(w, err)
		return
	}
	//
	client := &http.Client{}

	response, err := client.Do(request)

	log.Printf("authenticate request: %v", response)

	if err != nil {
		log.Println("err func call--1")
		err = app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()
	//defer func(Body io.ReadCloser) {
	//	err := Body.Close()
	//	if err != nil {
	//		err = app.errorJSON(w, err)
	//	}
	//}(response.Body)

	//make sure we get back the correct status
	if response.StatusCode == http.StatusUnauthorized {
		err = app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		err = app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	//create a variable we'll read response.body into
	var jsonFromService jsonResponse

	//decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		log.Println("err func call----2")
		err = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
