package handler

import (
	"context"
	"time"

	"github.com/chitchat-awsome/config"
	"github.com/chitchat-awsome/pkg/psqlconnector"
	"go.uber.org/zap"
)

type HandlerI interface{}

type Handler struct {
	Log    *zap.SugaredLogger
	Ctx    context.Context
	Config config.AppConfig
}

type handerServer struct {
	log          *zap.SugaredLogger
	ctx          context.Context
	address      string
	readtimeout  time.Duration
	writetimeout time.Duration
	static       string
	psql         psqlconnector.PsqlClientI
}

func NewHandler(deps Handler) HandlerI {
	psql := psqlconnector.NewPsqlClient(
		psqlconnector.PsqlDeps{
			Log:    deps.Log,
			Ctx:    deps.Ctx,
			Config: deps.Config.Psql,
		},
	)
	psql.Start()

	return &handerServer{
		log:          deps.Log,
		ctx:          deps.Ctx,
		address:      deps.Config.Server.Address,
		readtimeout:  deps.Config.Server.ReadTimeout,
		writetimeout: deps.Config.Server.WriteTimeout,
		static:       deps.Config.Server.Static,
		psql:         psql,
	}
}
