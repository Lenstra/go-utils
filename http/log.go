package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type LogOptions struct {
	SkipBody bool
}

func (o *LogOptions) optsOrDefault() *LogOptions {
	if o != nil {
		return o
	}
	return &LogOptions{}
}

// LogRequest logs a http.Request and reset the body properly if needed.
func LogRequest(event *zerolog.Event, req *http.Request, opts *LogOptions) error {
	opts = opts.optsOrDefault()
	if event.Enabled() {
		event = event.Str("method", req.Method).Str("path", req.URL.Path)

		if !opts.SkipBody && req.Body != nil {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return err
			}
			event = event.Str("body", string(body))
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		event.Msg("request")
	}

	return nil
}

// LogResponse logs a http.Response and reset the body properly if needed.
func LogResponse(event *zerolog.Event, resp *http.Response, opts *LogOptions) error {
	opts = opts.optsOrDefault()
	if event.Enabled() {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		event.Str("method", resp.Request.Method).Str("path", resp.Request.URL.Path).Int("status_code", resp.StatusCode).Str("body", string(body)).Msg("response")
	}
	return nil
}
