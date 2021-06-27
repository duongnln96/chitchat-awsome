package handler

import (
	"fmt"
	"net/http"

	"github.com/chitchat-awsome/internal/utils"
)

// GET /threads/new
// Show the new thread form page
func (s *server) NewThread(w http.ResponseWriter, r *http.Request) {
	_, err := s.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		utils.GenerateHTML(w, nil, "layout", "private.navbar", "new.thread")
	}
}

// POST /thread/create
// Create Thread
func (s *server) CreateThread(w http.ResponseWriter, r *http.Request) {
	ss, err := s.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err := r.ParseForm()
		if err != nil {
			s.log.Errorf("Cannot parse form create thread: %+v", err)
		}
		user, err := s.datahandler.GetUserBySession(&ss)
		if err != nil {
			s.log.Errorf("Cannot get user by session into create thread: %+v", err)
		}
		topic := r.PostFormValue("topic")
		if _, err := s.datahandler.CreateThread(topic, &user); err != nil {
			s.log.Errorf("Cannot create thread: %+v", err)
		}

		http.Redirect(w, r, "/", 302)
	}
}

// GET /thread/read
// Show the details of the thread, including the posts and the form to write a post
func (s *server) ReadThread(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	uuid := vals.Get("id")
	thread, err := s.datahandler.GetThreadByUUID(uuid)
	if err != nil {
		s.errorMessage(w, r, "Cannot read thread")
	} else {
		_, err := s.GetSession(w, r)
		if err != nil {
			utils.GenerateHTML(w, &thread, "layout", "public.navbar", "public.thread")
		} else {
			utils.GenerateHTML(w, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}

// POST /thread/post
// Create the post
func (s *server) PostThread(w http.ResponseWriter, r *http.Request) {
	ss, err := s.GetSession(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			s.log.Errorf("Cannot parse form %+v", err)
		}
		user, err := s.datahandler.GetUserBySession(&ss)
		if err != nil {
			s.log.Errorf("Cannot get user from session %+v", err)
		}
		body := r.PostFormValue("body")
		uuid := r.PostFormValue("uuid")
		thread, err := s.datahandler.GetThreadByUUID(uuid)
		if err != nil {
			s.errorMessage(w, r, "Cannot read thread")
		}
		if _, err := s.datahandler.CreatePost(&user, &thread, body); err != nil {
			s.log.Errorf("Cannot create post %+v", err)
		}
		url := fmt.Sprint("/thread/read?id=", uuid)
		http.Redirect(w, r, url, 302)
	}
}
