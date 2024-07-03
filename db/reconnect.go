package db

import (
	"context"
	"errors"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"io"
)

var ReconnectCount = 1

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type ExecQuerierReconnect struct {
	P Querier
}

func GetExecQuerierReconnect(e Querier) Querier {
	return &ExecQuerierReconnect{
		P: e,
	}
}

func (e *ExecQuerierReconnect) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	return e.P.Exec(ctx, sql, arguments...)
}

func (e *ExecQuerierReconnect) Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	for i := 0; i <= ReconnectCount; i++ {
		rows, err = e.P.Query(ctx, sql, args...)
		if err != nil {
			log.Errorf("TestLog:sql Query error :%v", err.Error())
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			log.Errorf("sql Query EOF, reconnect...")
			continue
		}
		return
	}
	return
}

func (e *ExecQuerierReconnect) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return &queryRow{
		ctx:  ctx,
		sql:  sql,
		args: args,
		p:    e.P,
	}
}

type queryRow struct {
	ctx  context.Context
	sql  string
	args []interface{}
	p    Querier
}

func (e *queryRow) Scan(dest ...interface{}) (err error) {
	for i := 0; i <= ReconnectCount; i++ {
		err = e.p.QueryRow(e.ctx, e.sql, e.args...).Scan(dest...)
		if err != nil {
			log.Infof("TestLog:sql QueryRow error :%v", err.Error())
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			log.Errorf("sql QueryRow EOF, and reconnect...")
			continue
		}
		return
	}
	return
}
