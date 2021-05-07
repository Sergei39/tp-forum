package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	middlewareLevel = "Middleware Level"
	usecaseLevel    = "Usecase Level"
	deliveryLevel   = "Delivery Level"
	repositoryLevel = "Repository Level"
	responseLevel   = "Response Level"
	utilsLevel      = "Utils Level"

	defaultRequestId = "000"
)

type Fields map[string]interface{}

func InitLogger() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		PadLevelText:  true,
	})
}

type EntryLog struct {
	level string
}

func (entry *EntryLog) Debug(ctx context.Context, fields Fields) {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}
	logrus.WithFields(logrus.Fields(fields)).
		Debug("[id: ", requestId, "] ", entry.level)
}

func (entry *EntryLog) Info(ctx context.Context, fields Fields) {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}
	logrus.WithFields(logrus.Fields(fields)).
		Info("[id: ", requestId, "] ", entry.level)
}

func (entry *EntryLog) Error(ctx context.Context, err error) {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}
	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Warn("[id: ", requestId, "] ", entry.level)
}

func (entry *EntryLog) InlineInfo(ctx context.Context, data ...interface{}) {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}
	logrus.WithFields(logrus.Fields{
		"info": data,
	}).Info("[id: ", requestId, "] ", entry.level)
}

func (entry *EntryLog) InlineDebug(ctx context.Context, data ...interface{}) {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}
	logrus.WithFields(logrus.Fields{
		"data": data,
	}).Debug("[id: ", requestId, "] ", entry.level)
}

func Middleware() *EntryLog {
	return &EntryLog{
		level: middlewareLevel,
	}
}

func Delivery() *EntryLog {
	return &EntryLog{
		level: deliveryLevel,
	}
}

func Usecase() *EntryLog {
	return &EntryLog{
		level: usecaseLevel,
	}
}

func Repo() *EntryLog {
	return &EntryLog{
		level: repositoryLevel,
	}
}

func Response() *EntryLog {
	return &EntryLog{
		level: responseLevel,
	}
}
