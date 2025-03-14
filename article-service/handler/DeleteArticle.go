package handler

import (
	"article-service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func DeleteArticle(c *gin.Context) {
	ctx := c.Request.Context()
	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)))
	defer span.End()

	articleId := c.Param("articleId")

	if err := repository.DeleteArticle(ctx, articleId); err != nil {
		log.Warnf("DeleteArticle Error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     err.Error(),
			"articleId": articleId,
		})
		return
	}

	c.Status(http.StatusAccepted)
}
