package ircpush

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
)

func init() {
	cinotify.Register("cave", caveHandler{})
}

// CaveNotification is the notification transmitted from a cavenotify request.
type CaveNotification struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Author string `json:"uploader"`
	Source string `json:"source"`
}

// phishcave.com/u/123123 : [mobile] phishpic_31231.png by Chetic
func (n CaveNotification) String() string {
	return fmt.Sprintf(
		"http://phishcave.com/u/%v : [%v] %v by %v",
		n.Id,
		n.Source,
		n.Name,
		n.Author,
	)
}

// caveHandler implements cinotify.Handler
type caveHandler struct {
}

// caveHandler handles any requests from phishcave.
func (_ caveHandler) Handle(r *http.Request) fmt.Stringer {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var n CaveNotification
	err := decoder.Decode(&n)
	if err != nil {
		cinotify.DoLog("cinotify/cave: Failed to decode json payload: ", err)
		return nil
	}

	return n
}

// Route creates a route that only a cavenotify client should hit.
func (_ caveHandler) Route(r *mux.Route) {
	r.Path("/").Methods("POST").Headers(
		"Content-Type", "application/json",
		"User-Agent", "cavenotify",
	)
}
