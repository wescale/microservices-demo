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

func GetArticle(c *gin.Context) {
	log.Println(c.Request.Header)

	ctx := c.Request.Context()

	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)))
	defer span.End()

	articles, err := repository.GetArticles(ctx, nil)
	if err != nil {
		log.Warnf("GetArticle Error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
	})
}
