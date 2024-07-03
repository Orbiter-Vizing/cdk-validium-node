package state

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgconn"
)

type execQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

var ReconnectCount = 1

type ExecQuerierReconnect struct {
	p execQuerier
}

func (e ExecQuerierReconnect) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	return e.p.Exec(ctx, sql, arguments)
}

func (e ExecQuerierReconnect) Query(ctx context.Context, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	for i := 0; i <= ReconnectCount; i++ {
		if rows, err = e.p.Query(ctx, sql, args); errors.Is(err, pgx.ErrNoRows) {
			return
		} else if err != nil {
			continue
		}
		return
	}
	return
}

func (e ExecQuerierReconnect) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return queryRow{
		ctx:  ctx,
		sql:  sql,
		args: args,
		p:    e.p,
	}
}

type queryRow struct {
	ctx  context.Context
	sql  string
	args []interface{}
	p    execQuerier
}

func (e queryRow) Scan(dest ...interface{}) (err error) {
	for i := 0; i <= ReconnectCount; i++ {
		if err = e.p.QueryRow(e.ctx, e.sql, e.args).Scan(dest); errors.Is(err, pgx.ErrNoRows) {
			return
		} else if err != nil {
			continue
		}
		return
	}
	return
}
