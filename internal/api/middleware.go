package api

import (
    "bytes"
    "fmt"
    "io"
    "time"

    "github.com/gin-gonic/gin"
)

type ResponseBodyWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (r ResponseBodyWriter) Write(b []byte) (int, error) {
    r.body.Write(b)
    return r.ResponseWriter.Write(b)
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()

        var requestBody []byte
        if c.Request.Body != nil {
            requestBody, _ = io.ReadAll(c.Request.Body)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
        }

        w := &ResponseBodyWriter{
            ResponseWriter: c.Writer,
            body:           &bytes.Buffer{},
        }
        c.Writer = w

        c.Next()

        endTime := time.Now()
        latency := endTime.Sub(startTime)
        statusCode := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path

        if path == "/api/find-recipes" {
            fmt.Printf("[DEBUG] Request: %s\n", string(requestBody))
            responseStr := w.body.String()
            if len(responseStr) > 500 {
                responseStr = responseStr[:500] + "... (dipotong)"
            }
            fmt.Printf("[DEBUG] Response: %s\n", responseStr)
        }

        if statusCode >= 400 {
            c.Error(fmt.Errorf("[ERROR] %s | %d | %s | %s %s | %s", endTime.Format("2006-01-02 15:04:05"), statusCode, latency, method, path, clientIP))
        } else {
            fmt.Printf("[INFO] %s | %d | %s | %s %s | %s\n", endTime.Format("2006-01-02 15:04:05"), statusCode, latency, method, path, clientIP)
        }
    }
}