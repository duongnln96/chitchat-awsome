package handler

import (
	"net/http"

	"github.com/chitchat-awsome/internal/data"
	"github.com/chitchat-awsome/internal/utils"
)

// GET /login
// Show the login page
func (s *server) Login(w http.ResponseWriter, r *http.Request) {
	t := utils.ParseTemplateFiles("login.layout", "public.navbar", "login")
	t.Execute(w, nil)
}

// GET /signup
// Show the signup page
func (s *server) Signup(w http.ResponseWriter, r *http.Request) {
	utils.GenerateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /authenticate
// Authenticate the user given the email and password
func (s *server) Authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.log.Errorf("Cannot parse form %+v", err)
	}

	user, err := s.datahandler.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		s.log.Debugf("User not found %+v", err)
	} else {
		s.log.Debugf("User found: %+v", user)
	}

	if user.Password == utils.Encrypt(r.PostFormValue("password")) {
		cookie := http.Cookie{
			Name:     "_cookie",
			HttpOnly: true,
		}
		userSs, err := s.datahandler.GetSessionByUser(&user)
		if err != nil {
			s.log.Debugf("Session for this user not valid, creating new: %+v", err)
			ss, err := s.datahandler.CreateSession(&user)
			if err != nil {
				s.log.Errorf("Cannot create session for this user %+v", ss)
			} else {
				cookie.Value = ss.Uuid
			}
		} else {
			cookie.Value = userSs.Uuid
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		s.log.Debug("Password is not correct")
		http.Redirect(w, r, "/login", 302)
	}
}

// POST /signup
// Create the user account
func (s *server) SignupUserAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.log.Errorf("Cannot parse form %+v", err)
	}

	user := data.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}

	usrRet, err := s.datahandler.CreateUser(&user)
	if err != nil {
		s.log.Errorf("Cannot create user: %+v, err: %+v", user, err)
	} else {
		s.log.Debugf("Created user: %+v", usrRet)
	}
	http.Redirect(w, r, "/", 302)
}

// GET /logout
// Logs the user out
func (s *server) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		ss := data.Session{
			Uuid: cookie.Value,
		}
		err := s.datahandler.DeleteSessionByUUID(&ss)
		if err != nil {
			s.log.Errorf("Cannot logout, cannot delete session %+v", err)
		}
	}
	http.Redirect(w, r, "/", 302)
}

// Get session
func (s *server) GetSession(w http.ResponseWriter, r *http.Request) (data.Session, error) {
	ss := data.Session{}
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		ss, err = s.datahandler.GetSessionByUUID(cookie.Value)
		if err != nil {
			s.log.Debugf("Session not Found %+v", err)
		}
	}

	return ss, err
}
