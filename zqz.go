package ircpush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
)

func init() {
	cinotify.Register("zqz", zqzHandler{})
}

// Notification is the notification transmitted from a zqzNotify request.
type ZQZNotification struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Author string `json:"author"`
}

// https://zqz.ca/f/0kla3x : phishpic_31231.png
func (n ZQZNotification) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("https://zqz.ca/f/")
	buffer.WriteString(n.ID)
	buffer.WriteString(" - ")
	buffer.WriteString(n.Name)

	if len(n.Author) > 0 {
		buffer.WriteString(" by ")
		buffer.WriteString(n.Author)
	}

	return buffer.String()
}

// zqzHandler implements cinotify.Handler.
type zqzHandler struct {
}

// zqzHandler handles any requests from zqz.
func (_ zqzHandler) Handle(r *http.Request) fmt.Stringer {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var n ZQZNotification
	err := decoder.Decode(&n)
	if err != nil {
		cinotify.DoLog("cinotify/zqz: Failed to decode json payload: ", err)
		return nil
	}

	return n
}

// Route creates a route that only a zqznotify client should hit.
func (_ zqzHandler) Route(r *mux.Route) {
	r.Path("/").Methods("POST").Headers(
		"Content-Type", "application/json",
		"User-Agent", "zqznotify",
	)
}
