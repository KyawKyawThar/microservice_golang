package main

import (
	"broker-service/event"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from,omitempty"`
	To      string `json:"to,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
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

	read, _ := io.ReadAll(r.Body)

	//{
	//	"action": "auth",
	//	"email": "admin@example.com",
	//	"password": "verysecret"
	//}

	//{
	//	"action": "logs",
	//	"name": "event",
	//	"data": "Some kind of data"
	//
	//}
	//jsonData := fmt.Sprintf("%s", read)

	var requestedPayload RequestPayload
	var auth AuthPayload
	var logs LogPayload
	var mails MailPayload

	_ = json.Unmarshal(read, &auth)
	_ = json.Unmarshal(read, &requestedPayload)
	_ = json.Unmarshal(read, &logs)
	_ = json.Unmarshal(read, &mails)

	//requestedPayload.Log = LogPayload{
	//	Name: "event",
	//	Data: "Some kind of data",
	//}

	switch requestedPayload.Action {
	case "auth":
		app.authenticate(w, auth)
	case "logs":
		//app.logItem(w, logs)
		app.logEventViaRabbit(w, logs)
	case "mails":
		app.SendMail(w, mails)
	default:
		_ = app.errorJSON(w, errors.New("unknown action"))
	}

	//var requestedPayload RequestPayload
	////
	//err := app.readJSON(w, r, &requestedPayload)
	////
	//if err != nil {
	//	err = app.errorJSON(w, err)
	//	return
	//}
	//
	//requestedPayload.Auth = AuthPayload{
	//	Email:    "admin@example.com",
	//	Password: "verysecret",
	//}
	//
	//requestedPayload.Log = LogPayload{
	//	Name: "event",
	//	Data: "Some kind of data",
	//}

	//{
	//	"action": "mails",
	//	"from": "me@example.com",
	//	"to": "kyawkyaw.thar84@gmail.com",
	//	"subject":"Test Mail",
	//	"message": "Hello Blockchain"
	//}
	//
	//switch requestedPayload.Action {
	//case "auth":
	//	app.authenticate(w, requestedPayload.Auth)
	//case "log":
	//	app.logItem(w, requestedPayload.Log)
	//default:
	//	_ = app.errorJSON(w, errors.New("unknown action"))
	//}
}

// logItem logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {

	log.Println("kkt", l)
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

	log.Println("authenticate", a)
	//{admin@example.com verysecret}

	// create json we'll send to the authentication-service
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	//call the service
	//authServiceURL := fmt.Sprintf("http://%s/authenticated", "127.0.0.1:81")
	authServiceURL := fmt.Sprintf("http://%s/authenticated", "auth-service")

	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
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

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {

	jsonData, err := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := fmt.Sprintf("http://%s/send", "mail-service")

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("err func call", err)
		err = app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-type", "application/json")
	client := http.Client{}

	response, err := client.Do(request)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	fmt.Printf("Response of status code %d\n", response.StatusCode)
	if response.StatusCode != http.StatusAccepted {
		err = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "Message sent to: " + msg.To

	log.Println("value of mail-service payload is %v\n", payload)
	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {

	err := app.pushToQueue(l.Name, l.Data)

	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)

	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	err = emitter.Push(string(j), "log.INFO")

	if err != nil {
		return err
	}
	return nil
}
