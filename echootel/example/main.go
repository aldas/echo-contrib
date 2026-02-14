package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo-contrib/echootel"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type user struct {
	ID   string
	Name string
}

// Try with `curl -v http://localhost:8080/users/123`
//
// * Host localhost:8080 was resolved.
// * IPv6: ::1
// * IPv4: 127.0.0.1
// *   Trying [::1]:8080...
// * Established connection to localhost (::1 port 8080) from ::1 port 36360
// * using HTTP/1.x
// > GET /users/123 HTTP/1.1
// > Host: localhost:8080
// > User-Agent: curl/8.18.0
// > Accept: */*
// >
// * Request completely sent off
// < HTTP/1.1 200 OK
// < Content-Type: application/json
// < Date: Sun, 08 Feb 2026 15:09:21 GMT
// < Content-Length: 38
// <
// {"ID":"123","Name":"otelecho tester"}
func main() {
	tp, err := initTracer()
	if err != nil {
		slog.Error("Failed to initialize otel tracer", "error", err)
		return
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("Failed to shutdown tracer provider", "error", err)
		}
	}()

	e := echo.New()
	e.Use(echootel.NewMiddlewareWithConfig(echootel.Config{
		ServerName:     "my-server",
		TracerProvider: tp,
		OnError: func(c *echo.Context, err error) {
			e.Logger.Error("otel middleware", "error", err)
		},
	}))

	e.GET("/users/:id", func(c *echo.Context) error {
		u := user{
			ID:   c.Param("id"),
			Name: "",
		}
		u.Name, _ = traceGetUser(c, u.ID)
		return c.JSON(http.StatusOK, u)
	})
	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start echo server", "error", err)
	}
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func traceGetUser(c *echo.Context, id string) (string, error) {
	tp, err := echo.ContextGet[trace.Tracer](c, echootel.TracerKey)
	if err != nil {
		return "", err
	}

	_, span := tp.Start(c.Request().Context(), "getUser", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	if id == "123" {
		return "otelecho tester", nil
	}
	return "unknown", nil
}
