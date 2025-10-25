package pgx

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"user-management-api/pkg/logger"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type PgxZerologTracer struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}
type QueryInfo struct {
	QueryName     string
	OperationType string
	CleanSQL      string
	OriginalSQL   string
}

var (
	sqlcNameRegex = regexp.MustCompile(`-- name:\s*(\w+)\s*:(\w+)`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	commentRegex  = regexp.MustCompile(`-- [^\r\n]*`)
)

func formatArgs(arg any) string {

	val := reflect.ValueOf(arg)                                 // Get the value of the argument passed
	if arg == nil || val.Kind() == reflect.Ptr && val.IsNil() { // Check if the argument is nil
		return "NULL"
	}
	if val.Kind() == reflect.Ptr { // Check if the argument is a pointer
		arg = val.Elem().Interface()
	}
	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return fmt.Sprintf("%s", v.Format("2006-01-02T15:04:05Z07:00"))
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''"))
	}
}
func replacePlaceHolder(sql string, args []any) string {
	for i, arg := range args {
		placehoder := fmt.Sprintf("$%d", i+1)
		sql = strings.ReplaceAll(sql, placehoder, formatArgs(arg))
	}
	return sql
}

func parseSQL(sql string) QueryInfo {
	info := QueryInfo{
		OriginalSQL: sql,
	}

	if matches := sqlcNameRegex.FindStringSubmatch(sql); len(matches) == 3 {
		info.QueryName = matches[1]
		info.OperationType = strings.ToUpper(matches[2])
	}

	// replace comment by " "
	cleanSQL := commentRegex.ReplaceAllString(sql, " ")
	cleanSQL = strings.TrimSpace(cleanSQL)
	cleanSQL = spaceRegex.ReplaceAllString(cleanSQL, " ")

	info.CleanSQL = cleanSQL

	return info
}

func (t *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {

	sql, _ := data["sql"].(string)
	args, _ := data["args"].([]any)
	duration, _ := data["time"].(time.Duration)

	queryInfo := parseSQL(sql)
	var finalSql string
	if len(args) > 0 {
		finalSql = replacePlaceHolder(queryInfo.CleanSQL, args)
	} else {
		finalSql = queryInfo.CleanSQL
	}
	baseLogger := t.Logger.With().
		Str("trace_id", logger.GetTraceId(ctx)).
		Dur("duration", duration).
		Str("sql_original", queryInfo.OriginalSQL).
		Str("sql", finalSql).
		Str("query_name", queryInfo.QueryName).
		Str("operation", queryInfo.OperationType).
		Interface("args", args)

	logger := baseLogger.Logger()

	if msg == "Query" && duration > t.SlowQueryLimit {
		logger.Warn().Str("event", "Slow Query").Msg("Slow SQL Query")
		return
	}

	if msg == "Query" {
		logger.Info().Str("event", "Query").Msg("Executed SQL")
		return
	}

}
