package handler

import (
	"context"
	"time"

	"github.com/chitchat-awsome/config"
	"github.com/chitchat-awsome/pkg/psqlconnector"
	"go.uber.org/zap"
)

type HandlerI interface{}

type HandlerServer struct {
	Log    *zap.SugaredLogger
	Ctx    context.Context
	Config config.AppConfig
	Psql   psqlconnector.PsqlClientI
}

type handlerServer struct {
	log          *zap.SugaredLogger
	ctx          context.Context
	address      string
	readtimeout  time.Duration
	writetimeout time.Duration
	static       string
	psql         psqlconnector.PsqlClientI
}

func NewHandler(deps HandlerServer) HandlerI {
	return &handlerServer{
		log:          deps.Log,
		ctx:          deps.Ctx,
		address:      deps.Config.Server.Address,
		readtimeout:  deps.Config.Server.ReadTimeout,
		writetimeout: deps.Config.Server.WriteTimeout,
		static:       deps.Config.Server.Static,
		psql:         deps.Psql,
	}
}

func (handler *handlerServer) StartServer() error {
	return nil
}

func (handler *handlerServer) Routing() error {
	return nil
}
