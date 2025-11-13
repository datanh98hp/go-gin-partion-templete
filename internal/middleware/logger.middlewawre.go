package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"
	"user-management-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CustomeResponse struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomeResponse) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)

}
func LoggerMiddleware(httpLogger *zerolog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		//get total time request
		start := time.Now()
		contentType := c.GetHeader("Content-Type")
		requestBody := make(map[string]any) // define request body map
		var formFiles []map[string]any
		//check content type
		if strings.HasPrefix(contentType, "multipart/form-data") {
			// check form data size max = 32MB
			if err := c.Request.ParseMultipartForm(32 << 20); err == nil && c.Request.MultipartForm != nil {
				// for values
				for k, v := range c.Request.MultipartForm.Value {
					if len(v) == 1 {
						requestBody[k] = v[0]
					} else {
						requestBody[k] = v
					}

				}
				//  for files
				for field, files := range c.Request.MultipartForm.File {
					for _, file := range files {
						formFiles = append(formFiles, map[string]any{
							"field":        field,
							"file":         file,
							"size":         formatFileSize(file.Size),
							"content_type": file.Header.Get("Content-Type"),
						})
					}
				}
				// write log
				if len(formFiles) > 0 {
					requestBody["files"] = formFiles
				}
			}
			log.Println("multipart/form-data")
		} else {
			// Content-Type: application/x-www-form-urlencoded
			//Content-Type: application/json
			/// chuyen toan bo du lieu tu body
			bodyByte, err := io.ReadAll(c.Request.Body) // read body request => bodyByte
			if err != nil {
				httpLogger.Error().Err(err).Msg("read body request error")
			}
			/// nap lai du lieu tra ve body
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
			//log.Printf("%s", bodyByte)

			if strings.HasPrefix(contentType, "application/json") {
				//Content-Type: application/json
				_ = json.Unmarshal(bodyByte, &requestBody) // convert json to map
			} else {
				// Content-Type: application/x-www-form-urlencoded
				value, _ := url.ParseQuery(string(bodyByte))

				for k, v := range value {
					// fix show value
					if len(v) == 1 {
						requestBody[k] = v[0]
					} else {
						requestBody[k] = v
					}

				}
				//log.Printf("=========%s", requestBody)
			}
		}
		/// Response Writer
		customeWriter := &CustomeResponse{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = customeWriter // ghi vao customeWriter response
		c.Next()                 // hanle request

		statusCode := c.Writer.Status()

		duration := time.Since(start)

		// response
		/// multipart/form-data
		responseContentType := c.Writer.Header().Get("Content-Type")

		responseBodyRaw := customeWriter.body.String()

		var responseBodyParsed interface{}
		if strings.HasPrefix(responseContentType, "image/") {
			responseBodyParsed = "[BINARY DATA]"
		} else if strings.HasPrefix(responseContentType, "application/json") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "{") ||
			strings.HasPrefix(strings.TrimSpace(responseBodyRaw), "[") {
			if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParsed); err != nil {
				responseBodyParsed = responseBodyRaw
			}
		} else {
			responseBodyParsed = responseBodyRaw
		}
		//log.Printf("responseBodyRaw: %s", responseBodyRaw)
		logEvent := httpLogger.Info()
		if statusCode >= 500 {
			logEvent = httpLogger.Error()
		} else if statusCode >= 400 {
			logEvent = httpLogger.Warn()
		}
		logEvent.
			Str("trace_id", logger.GetTraceId(c.Request.Context())).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("query", c.Request.URL.RawQuery).
			Str("ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Str("referer", c.Request.Referer()).
			Str("protocol", c.Request.Proto).
			Str("host", c.Request.Host).
			Str("remote_address", c.Request.RemoteAddr).
			Str("request_uri", c.Request.RequestURI).
			Int64("content_length", c.Request.ContentLength).
			Int("status_code", statusCode).
			Str("time", time.Now().Format("2006-01-02T15:04:05+07:00")).
			Interface("header", c.Request.Header).
			Interface("request_body", requestBody).
			Interface("response_body", responseBodyParsed).
			Int64("duration", duration.Microseconds()).
			Msg("HTTP request logs------")
	}
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<10))
	default:
		return fmt.Sprintf("%d B", (size))
	}
}
