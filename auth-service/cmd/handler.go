package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)



// Auth is a simple test handler for the broker
func (app *Config) Auth(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Message: "Hits the Auth Service",
		Error:   false,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}


func (app *Config) Authenticated(w http.ResponseWriter, r *http.Request) {

	log.Println("Authenticated func call")

	var resultPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &resultPayload)

	if err != nil {
		err = app.errorJSON(w, err, http.StatusBadRequest)

		return
	}

	//validate the use against DB
	user, err := app.Model.User.GetByEmail(resultPayload.Email)

	if user != nil {
		valid, err := user.PasswordMatches(resultPayload.Password)

		if err != nil || !valid {
			log.Println("code runnning.......")
			err = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
			return
		}
		log.Printf("logged in user %s, user data %v", user.Email, user)

		//log authentication
		err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))

		if err != nil {
			err = app.errorJSON(w, err, http.StatusBadRequest)
		}
		payload := jsonResponse{
			Message: fmt.Sprintf("Logged in user %s", user.Email),
			Data:    user,
			Error:   false,
		}

		_ = app.writeJSON(w, http.StatusAccepted, payload)
	} else {
		log.Println("user is nil")
	}

}

func (app *Config) logRequest(name, data string) error {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := fmt.Sprintf("http://%s/log", "logger-service")

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("logRequest %v\n", err)
		return err
	}

	client := &http.Client{}

	_, err = client.Do(request)

	if err != nil {

		log.Printf("logRequest client %v\n", err)
	}
	return nil
}
