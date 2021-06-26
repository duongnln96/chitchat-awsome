package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chitchat-awsome/config"
	"github.com/chitchat-awsome/internal/data"
	"github.com/chitchat-awsome/internal/handler"
	"github.com/chitchat-awsome/pkg/psqlconnector"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var globalContext context.Context

var rootCmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.GetConfig()
		log.Infof("Run application with config %+v", appConfig)

		// Start DB Connection
		psql := psqlconnector.NewPsqlClient(
			psqlconnector.PsqlDeps{
				Log:    log,
				Ctx:    globalContext,
				Config: appConfig.Psql,
			},
		)
		psql.Start()

		// Data Handler
		datahanler := data.NewDataHanler(
			data.DataHandlerDeps{
				Log: log,
				Ctx: globalContext,
				Db:  psql,
			},
		)

		// Server Handler Init
		server := handler.NewHandler(
			handler.ServerDeps{
				Log:         log,
				Ctx:         globalContext,
				Config:      appConfig.Server,
				DataHandler: datahanler,
			},
		)
		server.Start()
	},
}

func init() {
	prepareLogger()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-osSignals
		fmt.Println("Got os signal ", sig)
		cancel()
		os.Exit(0)
	}()

	globalContext = ctx
}

func prepareLogger() {
	logger, _ := zap.NewDevelopment()
	log = logger.Sugar()
	log.Info("Log is prepared in development mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
