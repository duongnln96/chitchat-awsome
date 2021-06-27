package data

import (
	"context"

	"github.com/chitchat-awsome/config"
	"github.com/chitchat-awsome/pkg/psqlconnector"
	"go.uber.org/zap"
)

var psql psqlconnector.PsqlClientI

type DataHandlerDeps struct {
	Log    *zap.SugaredLogger
	Ctx    context.Context
	Config *config.AppConfig
}

type DataHandlerI interface {
	CreateUser(*User) (User, error)
	GetUserByEmail(string) (User, error)
	GetUserByUUID(string) (User, error)
	GetUserBySession(*Session) (User, error)

	CreateSession(*User) (Session, error)
	GetSessionByUUID(string) (Session, error)
	GetSessionByUser(*User) (Session, error)
	DeleteSessionByUUID(*Session) error

	CreateThread(string, *User) (Thread, error)
	GetAllThreads() ([]Thread, error)
	GetThreadByUUID(string) (Thread, error)
	DeleteThread(Thread) error

	CreatePost(*User, *Thread, string) (Post, error)
	GetPostsIntoThread(Thread) ([]Post, error)
}

type dataHandler struct {
	log *zap.SugaredLogger
	ctx context.Context
}

func NewDataHanler(deps DataHandlerDeps) DataHandlerI {
	// Start DB Connection
	psql = psqlconnector.NewPsqlClient(
		psqlconnector.PsqlDeps{
			Log:    deps.Log,
			Ctx:    deps.Ctx,
			Config: deps.Config.Psql,
		},
	)
	psql.Start()

	return &dataHandler{
		log: deps.Log,
		ctx: deps.Ctx,
	}
}
