package data

import (
	"context"

	"github.com/chitchat-awsome/pkg/psqlconnector"
	"go.uber.org/zap"
)

type DataHandlerDeps struct {
	Log *zap.SugaredLogger
	Ctx context.Context
	Db  psqlconnector.PsqlClientI
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
	DeleteThread(Thread) error
}

type dataHandler struct {
	log *zap.SugaredLogger
	ctx context.Context
	db  psqlconnector.PsqlClientI
}

func NewDataHanler(deps DataHandlerDeps) DataHandlerI {
	return &dataHandler{
		log: deps.Log,
		ctx: deps.Ctx,
		db:  deps.Db,
	}
}
