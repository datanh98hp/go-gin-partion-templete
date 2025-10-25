package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Level     string
	FileName  string
	MaxSize   int
	MaxBackUp int
	MaxAge    int
	Compress  bool
	IsDev     string
}
type ContextKey string

const TraceIdKey ContextKey = "trace_id"

func NewLogger(config LoggerConfig) *zerolog.Logger {

	zerolog.TimeFieldFormat = time.RFC3339 //set time format in log
	lvl, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	var writer io.Writer

	// If dev mode, use console writer
	if config.IsDev == "development" {
		writer = PrettyLogJSONWriter{Writer: os.Stdout}
	} else {
		writer = &lumberjack.Logger{
			Filename:   config.FileName,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackUp,
			MaxAge:     config.MaxAge,   //days
			Compress:   config.Compress, // disabled by default
			LocalTime:  true,
		}
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()
	return &logger
}

type PrettyLogJSONWriter struct {
	Writer io.Writer
}

func (w PrettyLogJSONWriter) Write(p []byte) (n int, err error) {
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, p, "", " ")
	if err != nil {
		return w.Writer.Write(p)
	}

	return w.Writer.Write(prettyJSON.Bytes())
}

func GetTraceId(c context.Context) string {
	if traceId, ok := c.Value(TraceIdKey).(string); ok && traceId != "" {
		return traceId
	}
	return ""
}
