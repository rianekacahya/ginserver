package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rianekacahya/logger"
	"go.uber.org/zap"
)

type writter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w writter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request []byte
		dump := &writter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		if c.Request.Body != nil {
			request, _ = ioutil.ReadAll(c.Request.Body)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(request))
		c.Writer = dump
		c.Next()

		response := dump.body.String()

		if c.Writer.Status() >= http.StatusOK && c.Writer.Status() < http.StatusMultipleChoices {
			logger.Info(
				"http",
				zap.Int("status", c.Writer.Status()),
				zap.String("time", time.Now().Format(time.RFC1123Z)),
				zap.String("hostname", c.Request.Host),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.Any("req", json.RawMessage(request)),
				zap.Any("res", json.RawMessage(response)),
			)
		}
	}
}