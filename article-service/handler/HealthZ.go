package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func HealthZ(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)))
	defer span.End()

	c.Status(http.StatusOK)
}
