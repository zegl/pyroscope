package log

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func Logger(ctx context.Context) logrus.FieldLogger {
	logger := logrus.StandardLogger()
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		return logger.WithFields(logrus.Fields{
			"trace_id": spanCtx.TraceID().String(),
			"span_id":  spanCtx.SpanID().String(),
		})
	}
	return logger
}

type AppNameFieldDecorator struct {
	AppName   string
	Formatter logrus.Formatter
}

func (d AppNameFieldDecorator) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["app_name"] = d.AppName
	return d.Formatter.Format(entry)
}
