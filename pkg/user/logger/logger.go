package logger

import "context"

type Logger interface {
	LogNewUser(context.Context) error
}

type Log struct {
	Logger
}
