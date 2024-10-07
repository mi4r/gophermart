package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	gzipHeader = "gzip"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// For the versatility of the content-type
func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	if grw.Header().Get(echo.HeaderContentType) == "" {
		grw.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}
	return grw.Writer.Write(b)
}

// func (s *Server) AuthWithCookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {

// 		return nil
// 	}
// }

// func (s *Server) SetContentTypeJSONMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		return nil
// 	}
// }

func GzipMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check headers
		acceptEncoding := c.Request().Header.Get(echo.HeaderAcceptEncoding)
		if !strings.Contains(acceptEncoding, gzipHeader) {
			return next(c)
		}

		contentEncoding := c.Request().Header.Get(echo.HeaderContentEncoding)
		if strings.Contains(contentEncoding, gzipHeader) {
			// Create a gzip reader for the request body
			gzipReader, err := gzip.NewReader(c.Request().Body)
			if err != nil {
				return err
			}
			defer gzipReader.Close()

			// Replace the request body with the gzip reader
			c.Request().Body = io.NopCloser(gzipReader)
		}

		c.Response().Header().Set(echo.HeaderContentEncoding, gzipHeader)
		// Create gzip writer
		gzipWriter := gzip.NewWriter(c.Response().Writer)
		defer gzipWriter.Close()

		// Replace the response writer with a gzip response writer
		grw := &gzipResponseWriter{
			Writer:         gzipWriter,
			ResponseWriter: c.Response().Writer,
		}
		c.Response().Writer = grw

		// Call the next handler
		if err := next(c); err != nil {
			c.Error(err)
		}

		// Flush the gzip writer to ensure all data is written
		if err := gzipWriter.Flush(); err != nil {
			return err
		}

		return nil
	}
}
