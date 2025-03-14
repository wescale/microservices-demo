package handler

import (
	"article-service/model"
	"article-service/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func AddArticle(c *gin.Context) {
	ctx := c.Request.Context()

	span := trace.SpanFromContext(otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header)))
	defer span.End()

	var addArticleRequest AddArticleRequest
	if err := c.ShouldBindJSON(&addArticleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := addArticleRequest.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := repository.AddArticle(c, addArticleRequest.toArticle())
	if err != nil {
		log.Warnf("AddArticle Error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

type AddArticleRequest struct {
	ArticleName        string `json:"name"`
	ArticleDescription string `json:"description"`
}

func (r AddArticleRequest) validate() error {
	if r.ArticleName == "" {
		return fmt.Errorf("articleName is required")
	}
	if r.ArticleDescription == "" {
		return fmt.Errorf("articleDescription is required")
	}
	return nil
}

func (r AddArticleRequest) toArticle() *model.Article {
	return &model.Article{
		Title:       r.ArticleName,
		Description: r.ArticleDescription,
	}
}
