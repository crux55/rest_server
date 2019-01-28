package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var URL_BASE = "http://localhost:8000"
var URL_LOGIN = URL_BASE + "/login"
var URL_REGISTER = URL_BASE + "/register"
var URL_ACCOUNT = URL_BASE + "/account?token="

var CONTENT_TYPE = "application/json"

var DEFAULT_EMAIL_ADDRESS = "newuser@somedomain.com"
var DEFAULT_PASSWORD = "defaultpassword"

//Test a user can be created a logged in
func TestLogin(t *testing.T) {
	message := map[string]interface{}{
		"emailaddress": DEFAULT_EMAIL_ADDRESS,
		"password":  DEFAULT_PASSWORD,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(URL_LOGIN, CONTENT_TYPE, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["token"])
	assert.Equal(t, 200, resp.StatusCode, "Expected: Status OK")
}

// Test a user can be created
func TestCreateUsers(t *testing.T) {
	message := map[string]interface{}{
		"emailaddress": DEFAULT_EMAIL_ADDRESS,
		"password":  DEFAULT_PASSWORD,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(URL_REGISTER, CONTENT_TYPE, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["message"])
	assert.Equal(t, 200, resp.StatusCode, "Expected: Status OK")
}

// Test a logged in user can access the protected endpoint
func TestProtectedEndpoint(t *testing.T) {

	message := map[string]interface{}{
		"emailaddress": DEFAULT_EMAIL_ADDRESS,
		"password":  DEFAULT_PASSWORD,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(URL_LOGIN, CONTENT_TYPE, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["token"])

	token := result["token"]

	resp, err = http.Get(URL_ACCOUNT + token.(string))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s %d", "Status code: ", resp.StatusCode))
	log.Println(body)
	assert.Equal(t, 200, resp.StatusCode, "Expected: Method not allowed")
}

// Check the health check works
func TestHealthCheckHandler(t *testing.T) {
	resp, err := http.Get("http://localhost:8000/healthcheck")
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s %d", "Status code: ", resp.StatusCode))
	assert.Equal(t, 200, resp.StatusCode, "Expected: Status OK")
}
