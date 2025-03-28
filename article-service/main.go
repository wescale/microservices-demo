package main

import (
	"article-service/handler"
	"article-service/repository"
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	log.Infoln("-= Article Service =-")
	loadConfig()
	initDatabase()
	if os.Getenv("OTLP_ENDPOINT") != "" {
		initTracer()
	}
	loadApiServer()
}

// loadConfig define the default values and loads the user configuration from config.yaml
func loadConfig() {
	viper.SetDefault("listen", ":8080")
	viper.SetDefault("mongodbUri", "mongodb://localhost:27017/alpha-articles")
	err := viper.BindEnv("mongodbUri", "MONGODB_URI")
	if err != nil {
		log.Warnln(err)
	}
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("/etc/article-service/")
	if err := viper.ReadInConfig(); err != nil {
		log.Warnln(err)
	}
}

// initDatabase initialize the database connection
func initDatabase() {
	mongodbUri := viper.GetString("mongodbUri")
	if err := repository.Initialize(mongodbUri); err != nil {
		log.Errorf("Failed to connect to %s", mongodbUri)
		log.Panicln(err)
	}
	log.Infof("Connected to %s", mongodbUri)
}

func initTracer() func(context.Context) error {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpointURL(os.Getenv("OTLP_ENDPOINT")),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", "article-service"),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithBatcher(exporter),
			trace.WithResource(resources),
		),
	)

	return exporter.Shutdown
}

// loadApiServer initialize the API server with a cors middleware and define routes to be served.
// This function is blocking: it will wait until the server returns an error
func loadApiServer() {
	Router := gin.New()
	Router.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true //return origin == "http://localhost:3001"
		},
		MaxAge: 12 * time.Hour,
	}))

	otelginOption := otelgin.WithPropagators(propagation.TraceContext{})

	Router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/", "/healthz"),
		gin.Recovery(),
		otelgin.Middleware("article-service", otelginOption),
	)

	Router.GET("/", handler.HealthZ)
	Router.GET("/healthz", handler.HealthZ)
	Router.GET("/article/", handler.GetArticle)
	Router.POST("/article/", handler.AddArticle)
	Router.DELETE("/article/:articleId/", handler.DeleteArticle)

	listenAddress := viper.GetString("listen")
	err := Router.Run(listenAddress)
	log.Panicln(err)
}
