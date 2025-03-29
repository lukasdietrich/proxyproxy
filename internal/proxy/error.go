package proxy

import (
	"fmt"
	"net/http"
)

var (
	_ error = Error{}
)

type Error struct {
	Status int
	Text   string
}

func (e Error) Error() string {
	return fmt.Sprintf("http error status=%d %s", e.Status, e.Text)
}

func (e Error) write(w http.ResponseWriter) {
	http.Error(w, e.Text, e.Status)
}
