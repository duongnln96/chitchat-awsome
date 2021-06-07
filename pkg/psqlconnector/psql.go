package psqlconnector

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type RowsHandlerFunc func(*sql.Rows) bool
type RowHandlerFunc func(*sql.Row) error

type PsqlDeps struct {
	Log    *zap.SugaredLogger
	Ctx    context.Context
	Config PsqlConfigurations
}

type PsqlClientI interface {
	Start()
	Stop()
	DBExec(string, ...interface{}) error
	DBQueryRow(RowHandlerFunc, string, ...interface{}) error
	DBQueryRows(RowsHandlerFunc, string, ...interface{}) error
}

type psqlClient struct {
	log                 *zap.SugaredLogger
	ctx                 context.Context
	host                string
	port                int
	username            string
	password            string
	dbname              string
	querytimeout        time.Duration
	healthcheckInterval time.Duration
	isOK                bool
	client              *sql.DB
	stopChannel         chan bool
}

func NewPsqlClient(deps PsqlDeps) PsqlClientI {
	dbInfo := deps.Config.GetConfig()
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		deps.Log.Fatalf("cannot connect to psql: %s error: %s", dbInfo, err)
	}

	timeoutCtx, cancel := context.WithTimeout(deps.Ctx, deps.Config.QueryTimeout)
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		deps.Log.Fatalf("cannot ping to psql: %s error: %s", dbInfo, err)
	}

	return &psqlClient{
		log:                 deps.Log,
		ctx:                 deps.Ctx,
		host:                deps.Config.Host,
		port:                deps.Config.Port,
		username:            deps.Config.Username,
		password:            deps.Config.Password,
		dbname:              deps.Config.DBname,
		querytimeout:        deps.Config.QueryTimeout,
		healthcheckInterval: deps.Config.HealthcheckInterval,
		isOK:                true,
		client:              db,
		stopChannel:         make(chan bool),
	}
}

func (pc *psqlClient) startHealthCheck() {
	go func() {
		ticker := time.NewTicker(5)
		for {
			select {
			case <-pc.stopChannel:
				return
			case <-ticker.C:
				timeoutCtx, cancel := context.WithTimeout(pc.ctx, pc.querytimeout)
				pc.isOK = pc.client.PingContext(timeoutCtx) == nil
				cancel()
			}

		}
	}()
}

func (pc *psqlClient) connIsOK() bool {
	return pc.isOK
}

func (pc *psqlClient) Start() {
	pc.startHealthCheck()
}

func (pc *psqlClient) Stop() {
	pc.stopChannel <- true
}

func (pc *psqlClient) DBExec(cmd string, args ...interface{}) error {
	if !pc.connIsOK() {
		return fmt.Errorf("cannot connect to database")
	}

	timeoutCtx, cancel := context.WithTimeout(pc.ctx, pc.querytimeout)
	defer cancel()
	_, err := pc.client.ExecContext(timeoutCtx, cmd, args...)

	return err
}

func (pc *psqlClient) DBQueryRow(handlerFunc RowHandlerFunc, cmd string, args ...interface{}) error {
	if !pc.connIsOK() {
		return fmt.Errorf("cannot connect to database")
	}

	timeoutCtx, cancel := context.WithTimeout(pc.ctx, pc.querytimeout)
	defer cancel()

	return handlerFunc(pc.client.QueryRowContext(timeoutCtx, cmd, args...))
}

func (pc *psqlClient) DBQueryRows(handlerFunc RowsHandlerFunc, cmd string, args ...interface{}) error {
	if !pc.connIsOK() {
		return fmt.Errorf("cannot connect to database")
	}

	timeoutCtx, cancel := context.WithTimeout(pc.ctx, pc.querytimeout)
	defer cancel()

	rows, err := pc.client.QueryContext(timeoutCtx, cmd, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if !handlerFunc(rows) {
			return fmt.Errorf("Error while processing rows")
		}
	}

	return nil
}
