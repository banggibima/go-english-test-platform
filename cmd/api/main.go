package main

import (
	"log"

	_ "github.com/banggibima/go-english-test-platform/docs"

	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/banggibima/go-english-test-platform/config"
	attemptanswers "github.com/banggibima/go-english-test-platform/internal/attempt_answers"
	"github.com/banggibima/go-english-test-platform/internal/attempts"
	"github.com/banggibima/go-english-test-platform/internal/auth"
	"github.com/banggibima/go-english-test-platform/internal/dashboard"
	"github.com/banggibima/go-english-test-platform/internal/files"
	"github.com/banggibima/go-english-test-platform/internal/questions"
	"github.com/banggibima/go-english-test-platform/internal/results"
	"github.com/banggibima/go-english-test-platform/internal/roles"
	"github.com/banggibima/go-english-test-platform/internal/sections"
	"github.com/banggibima/go-english-test-platform/internal/tests"
	"github.com/banggibima/go-english-test-platform/internal/users"
	"github.com/banggibima/go-english-test-platform/pkg/cache"
	"github.com/banggibima/go-english-test-platform/pkg/database"
	"github.com/banggibima/go-english-test-platform/pkg/logger"
	"github.com/banggibima/go-english-test-platform/pkg/metrics"
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/banggibima/go-english-test-platform/pkg/queue"
	"github.com/banggibima/go-english-test-platform/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Go English Test Platform API
// @version 1.0
// @description Backend API for English proficiency testing platform.
// @host localhost:8081
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logg := logger.New()

	metrics.Register()

	pg, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	rds, err := cache.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer rds.Close()

	rabbit, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close()

	if err := rabbit.DeclareQueue("score-attempt"); err != nil {
		log.Fatal(err)
	}

	minioClient, err := storage.NewMinIO(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_ = minioClient

	router := gin.New()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	router.Use(
		gin.Recovery(),
		middleware.CORS(),
		middleware.Logger(logg),
		middleware.Metrics(),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	authRepository := auth.NewRepository(pg)
	authService := auth.NewService(authRepository, cfg)
	authHandler := auth.NewHandler(authService)
	userRepository := users.NewRepository(pg)
	userService := users.NewService(userRepository)
	userHandler := users.NewHandler(userService)
	roleRepository := roles.NewRepository(pg)
	roleService := roles.NewService(roleRepository)
	roleHandler := roles.NewHandler(roleService)
	testRepository := tests.NewRepository(pg)
	testService := tests.NewService(testRepository, rds)
	testHandler := tests.NewHandler(testService)
	sectionRepository := sections.NewRepository(pg)
	sectionService := sections.NewService(sectionRepository)
	sectionHandler := sections.NewHandler(sectionService)
	questionRepository := questions.NewRepository(pg)
	questionService := questions.NewService(questionRepository)
	questionHandler := questions.NewHandler(questionService)
	attemptRepository := attempts.NewRepository(pg)
	attemptService := attempts.NewService(attemptRepository, rabbit)
	attemptHandler := attempts.NewHandler(attemptService)
	attemptAnswerRepository := attemptanswers.NewRepository(pg)
	attemptAnswerService := attemptanswers.NewService(attemptAnswerRepository)
	attemptAnswerHandler := attemptanswers.NewHandler(attemptAnswerService)
	resultRepository := results.NewRepository(pg)
	resultService := results.NewService(resultRepository)
	resultHandler := results.NewHandler(resultService)
	fileRepository := files.NewRepository(pg)
	fileService := files.NewService(fileRepository, minioClient)
	fileHandler := files.NewHandler(fileService)
	dashboardRepository := dashboard.NewRepository(pg)
	dashboardService := dashboard.NewService(dashboardRepository, rds)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	api := router.Group("/api")
	auth.RegisterRoutes(api, authHandler)
	users.RegisterRoutes(api, userHandler, cfg.JWTSecret)
	roles.RegisterRoutes(api, roleHandler, cfg.JWTSecret)
	tests.RegisterRoutes(api, testHandler, cfg.JWTSecret)
	sections.RegisterRoutes(api, sectionHandler, cfg.JWTSecret)
	questions.RegisterRoutes(api, questionHandler, cfg.JWTSecret)
	attempts.RegisterRoutes(api, attemptHandler, cfg.JWTSecret)
	attemptanswers.RegisterRoutes(api, attemptAnswerHandler, cfg.JWTSecret)
	results.RegisterRoutes(api, resultHandler, cfg.JWTSecret)
	files.RegisterRoutes(api, fileHandler, cfg.JWTSecret)
	dashboard.RegisterRoutes(api, dashboardHandler, cfg.JWTSecret)

	logg.Info("postgres connected")
	logg.Info("redis connected")
	logg.Info("rabbitmq connected")
	logg.Info("minio connected")
	logg.Info("server started", "port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
