package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var a App

var createdUser User
var createdSnippet Snippet

const extensionQueryUUID = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

const tableCreationQueryUsers = `
CREATE TABLE IF NOT EXISTS users (
	user_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	username VARCHAR(25) NOT NULL,
    mail VARCHAR(50) NOT NULL,
	password VARCHAR(60) NOT NULL,
    UNIQUE(mail)
)`

const tableCreationQuerySnippets = `
CREATE TABLE IF NOT EXISTS snippets (
	snippet_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	owner UUID REFERENCES users(user_id) NOT NULL,
	language VARCHAR(20) NULL,
	title VARCHAR(30) NULL,
	category VARCHAR(50) NULL,
	code VARCHAR NULL
)`

func TestMain(m *testing.M) {
	a.Initialize("admin", "123", "SnippetsTest", false)
	ensureExtensionExists()
	ensureTablesExist()
	clearTable()

	code := m.Run()

	clearTable()
	os.Exit(code)
}

func ensureExtensionExists() {
	if _, err := a.DB.Exec(extensionQueryUUID); err != nil {
		log.Fatal(err)
	}
}

func ensureTablesExist() {
	if _, err := a.DB.Exec(tableCreationQueryUsers); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableCreationQuerySnippets); err != nil {
		log.Fatal(err)
	}
}

// Order is important, as tables with fks need to be deleted first
func clearTable() {
	a.DB.Exec("DELETE FROM snippets")
	a.DB.Exec("DELETE FROM users")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

/* ******************************
 * Testing functions
 *******************************/

func TestWelcomeMessage(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body == "" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"mail":"test@mail", "username": "User1", "password": "123"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["mail"] != "test@mail" {
		t.Errorf("Expected snippet name to be 'test@mail'. Got '%v'", m["mail"])
	}

	if m["username"] != "User1" {
		t.Errorf("Expected snippet name to be 'User1'. Got '%v'", m["username"])
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&createdUser); err != nil {
		t.Errorf("Invalid request payload")
		return
	}
}

func TestCreateSnippet(t *testing.T) {
	bodyString := fmt.Sprintf(`{"owner":"%s", "title": "snippet1", "language": "python", "category": "Hello World", "code": "print(\"Hello World from Go Rest API\")"}`, createdUser.ID)
	jsonStr := []byte(bodyString)
	req, _ := http.NewRequest("POST", fmt.Sprintf(`/user/%s/snippets`, createdUser.ID), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["owner"] != createdUser.ID {
		t.Errorf("Expected owner to be 'createdUserID'. Got '%v'", m["owner"])
	}

	if m["title"] != "snippet1" {
		t.Errorf("Expected snippet title to be 'snippet1'. Got '%v'", m["title"])
	}

	if m["language"] != "python" {
		t.Errorf("Expected snippet language to be 'python'. Got '%v'", m["language"])
	}

	if m["category"] != "Hello World" {
		t.Errorf("Expected snippet about to be 'Hello World'. Got '%v'", m["category"])
	}

	if m["code"] != `print("Hello World from Go Rest API")` {
		t.Errorf(`Expected snippet about to be 'print("Hello World from Go Rest API")'. Got '%v'`, m["code"])
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&createdSnippet); err != nil {
		t.Errorf("Invalid request payload")
		return
	}
}

func TestCreateSnippet2(t *testing.T) {
	bodyString := fmt.Sprintf(`{"owner":"%s", "title": "snippet2", "language": "python", "category": "Hello World", "code": "print(\"Hello World from Go Rest API\")"}`, createdUser.ID)
	jsonStr := []byte(bodyString)
	req, _ := http.NewRequest("POST", fmt.Sprintf(`/user/%s/snippets`, createdUser.ID), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["owner"] != createdUser.ID {
		t.Errorf("Expected owner to be 'createdUserID'. Got '%v'", m["owner"])
	}

	if m["title"] != "snippet2" {
		t.Errorf("Expected snippet title to be 'snippet2'. Got '%v'", m["title"])
	}

	if m["language"] != "python" {
		t.Errorf("Expected snippet language to be 'python'. Got '%v'", m["language"])
	}

	if m["category"] != "Hello World" {
		t.Errorf("Expected snippet about to be 'Hello World'. Got '%v'", m["category"])
	}

	if m["code"] != `print("Hello World from Go Rest API")` {
		t.Errorf(`Expected snippet about to be 'print("Hello World from Go Rest API")'. Got '%v'`, m["code"])
	}
}

func TestGetSnippet(t *testing.T) {
	connectionString := fmt.Sprintf(`/user/%s/snippets/%s`, createdUser.ID, createdSnippet.ID)
	req, _ := http.NewRequest("GET", connectionString, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	snippet := Snippet{}
	err := json.Unmarshal(response.Body.Bytes(), &snippet)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	if !reflect.DeepEqual(snippet, createdSnippet) {
		t.Errorf("Error while deepequal, object values are not the same")
	}
}

func TestGetSnippets(t *testing.T) {
	connectionString := fmt.Sprintf(`/user/%s/snippets/`, createdUser.ID)
	req, _ := http.NewRequest("GET", connectionString, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	snippets := make([]Snippet, 0)
	err := json.Unmarshal(response.Body.Bytes(), &snippets)

	if err != nil {
		t.Errorf("Error while unmarshaling: " + err.Error())
	}

	createdUser.Snippets = []Snippet{}
}
