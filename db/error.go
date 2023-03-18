package db

import "errors"

var (
	ErrNoActiveTransaction = errors.New("no transaction is in progress")
)
