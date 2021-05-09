package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

const (
	middlewareLevel = "Middleware"
	usecaseLevel    = "Usecase"
	deliveryLevel   = "Delivery"
	repositoryLevel = "Repository"
	responseLevel   = "Response"
	startLevel      = "Start"
	utilsLevel      = "Utils"

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
	level    string
	funcName string
}

func (entry *EntryLog) AddFuncName(name string) *EntryLog {
	entry.funcName = name
	return entry
}

func (entry *EntryLog) createMetaInfo(ctx context.Context) []interface{} {
	requestId := ctx.Value("request_id")
	if requestId == nil {
		requestId = defaultRequestId
	}

	metaInfo := make([]interface{}, 0, 2)
	metaInfo = append(metaInfo, "[id: ", requestId, "] ", entry.level)

	if entry.funcName != "" {
		metaInfo = append(metaInfo, " [", entry.funcName, "]")
	}

	return metaInfo
}

func (entry *EntryLog) Debug(ctx context.Context, fields Fields) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields(fields)).
		Debug(metaInfo...)
}

func (entry *EntryLog) Info(ctx context.Context, fields Fields) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields(fields)).
		Info(metaInfo...)
}

func (entry *EntryLog) Error(ctx context.Context, err error) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Warn(metaInfo...)
}

func (entry *EntryLog) InlineInfo(ctx context.Context, data ...interface{}) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields{
		"info": data,
	}).Info(metaInfo...)
}

func (entry *EntryLog) InlineDebug(ctx context.Context, data ...interface{}) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields{
		"data": data,
	}).Debug(metaInfo...)
}

func (entry *EntryLog) Fatal(ctx context.Context, err error) {
	metaInfo := entry.createMetaInfo(ctx)

	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Error(metaInfo...)
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

func Start() *EntryLog {
	return &EntryLog{
		level: startLevel,
	}
}
