package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/chitchat-awsome/config"
	"github.com/chitchat-awsome/internal/data"
	"go.uber.org/zap"
)

type ServerI interface {
	Start() error
	Stop()
}

type ServerDeps struct {
	Log         *zap.SugaredLogger
	Ctx         context.Context
	Config      config.ServerConfig
	DataHandler data.DataHandlerI
}

type server struct {
	log          *zap.SugaredLogger
	ctx          context.Context
	address      string
	readtimeout  time.Duration
	writetimeout time.Duration
	static       string
	datahandler  data.DataHandlerI
}

func NewHandler(deps ServerDeps) ServerI {
	return &server{
		log:          deps.Log,
		ctx:          deps.Ctx,
		address:      deps.Config.Address,
		readtimeout:  deps.Config.ReadTimeout,
		writetimeout: deps.Config.WriteTimeout,
		static:       deps.Config.Static,
		datahandler:  deps.DataHandler,
	}
}

func (s *server) Start() error {
	mux := s.routingFunc()
	// starting up the server
	server := &http.Server{
		Addr:           s.address,
		Handler:        mux,
		ReadTimeout:    time.Duration(int64(s.readtimeout) * int64(time.Second)),
		WriteTimeout:   time.Duration(int64(s.writetimeout) * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	err := server.ListenAndServe()
	if err != nil {
		s.log.Panicf("Cannot Start Server %+v\n", err)
	}
	return nil
}

func (s *server) Stop() {

}

func (s *server) routingFunc() *http.ServeMux {
	// handle static assets
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(s.static))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// all route patterns matched here
	// route handler functions defined in other files
	//

	// index
	mux.HandleFunc("/", s.Index)
	// error
	mux.HandleFunc("/err", s.Error)

	// defined in authhandler.go
	mux.HandleFunc("/login", s.Login)
	mux.HandleFunc("/logout", s.Logout)
	mux.HandleFunc("/signup", s.Signup)
	mux.HandleFunc("/signup_account", s.SignupUserAccount)
	mux.HandleFunc("/authenticate", s.Authenticate)

	// defined in threadhandler.go
	mux.HandleFunc("/thread/new", s.NewThread)
	mux.HandleFunc("/thread/create", s.CreateThread)
	mux.HandleFunc("/thread/post", s.PostThread)
	mux.HandleFunc("/thread/read", s.ReadThread)

	return mux
}
