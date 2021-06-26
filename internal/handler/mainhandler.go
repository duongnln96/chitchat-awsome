package handler

import (
	"net/http"
	"strings"

	"github.com/chitchat-awsome/internal/utils"
)

// GET /err?msg=
// shows the error message page
func (s *server) Error(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	_, err := s.GetSession(w, r)
	if err != nil {
		utils.GenerateHTML(w, vals.Get("msg"), "layout", "public.navbar", "error")
	} else {
		utils.GenerateHTML(w, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}

func (s *server) errorMessage(w http.ResponseWriter, r *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(w, r, strings.Join(url, ""), 302)
}

func (s *server) Index(w http.ResponseWriter, r *http.Request) {
	threads, err := s.datahandler.GetAllThreads()
	if err != nil {
		s.errorMessage(w, r, "Cannot get all threads")
	} else {
		_, err := s.GetSession(w, r)
		if err != nil {
			utils.GenerateHTML(w, threads, "layout", "public.navbar", "index")
		} else {
			utils.GenerateHTML(w, threads, "layout", "private.navbar", "index")
		}
	}
}
