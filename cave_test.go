package ircpush

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	. "testing"

	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
)

var testNotification = CaveNotification{
	Name:   "name",
	Id:     "id",
	Author: "author",
	Source: "source",
}

func TestString(t *T) {
	expect := "http://phishcave.com/u/id : [source] name by author"

	if got := testNotification.String(); got != expect {
		t.Error("Expected: %s, got: %s", expect, got)
	}
}

func TestHandle(t *T) {
	var err error
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	if err = encoder.Encode(testNotification); err != nil {
		t.Error("Failed to jsonify payload: ", err)
	}

	var req *http.Request
	if req, err = http.NewRequest("POST", "/", buf); err != nil {
		t.Error("Error creating mock request: ", err)
	}

	c := caveHandler{}
	note := c.Handle(req)
	if note == nil {
		t.Error("Expected to get a notification, got nil.")
	}

	if notification, ok := note.(CaveNotification); !ok {
		t.Error("Expected to get a Notification type.")
	} else if notification != testNotification {
		t.Error("Expected an unaltered payload.")
	}
}

func TestHandleFail(t *T) {
	var err error
	buf := bytes.NewBufferString("{!$@($*&@&$)(*$)*&@$)")
	logger := &bytes.Buffer{}

	cinotify.Logger = log.New(logger, "", log.LstdFlags)

	var req *http.Request
	if req, err = http.NewRequest("POST", "/", buf); err != nil {
		t.Error("Error creating mock request: ", err)
	}

	if 0 != logger.Len() {
		t.Error("How could something be logged at this point?")
	}

	d := caveHandler{}
	note := d.Handle(req)
	if note != nil {
		t.Error("Expected an error to occur.")
	}

	if 0 == logger.Len() {
		t.Error("Expected something to be written to the log.")
	}
}

func TestRoute(t *T) {
	var err error

	d := caveHandler{}
	router := mux.NewRouter()
	r := router.NewRoute()

	d.Route(r)
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	resp := httptest.NewRecorder()
	var req *http.Request
	if req, err = http.NewRequest("POST", "/", nil); err != nil {
		t.Error("Error creating mock request: ", err)
	}
	req.Header.Add("User-Agent", "cavenotify")
	req.Header.Add("Content-Type", "application/json")

	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Error("Route did not match request.")
	}
}
