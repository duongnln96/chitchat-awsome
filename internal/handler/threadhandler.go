package handler

import (
	"net/http"
)

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

// 
