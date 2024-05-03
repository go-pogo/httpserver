// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelaccesslog

import (
	"context"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

// Middleware wraps a [http.Handler] so it's request uri is added to the
// [trace.Span] derived from the [http.Request]'s context.
// This is a workaround for https://github.com/open-telemetry/opentelemetry-go/commit/7b749591320bfcdef2061f4d4f5aa533ab76b47f
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		trace.SpanFromContext(req.Context()).SetAttributes(
			semconv.URLPath(req.URL.Path),
			semconv.URLQuery(req.URL.RawQuery),
			// keep to stay backwards compatible with otelhttp
			semconv.HTTPTargetKey.String(req.RequestURI),
		)
		next.ServeHTTP(wri, req)
	})
}

// WithHandlerName adds name as value to the [http.Request]'s context. It
// should be used on a per route/handler basis.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		SetHandlerName(req.Context(), name)
		next.ServeHTTP(wri, req)
	})
}

func SetHandlerName(ctx context.Context, name string) {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.SetAttributes(semconv.CodeFunctionKey.String(name))
	}
}
