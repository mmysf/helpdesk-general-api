package main

import (
	http_agent "app/app/delivery/http/agent"
	http_member "app/app/delivery/http/member"
	"app/app/delivery/http/middleware"
	http_superadmin "app/app/delivery/http/superadmin"
	http_webhook "app/app/delivery/http/webhook"
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	s3Repo "app/app/repository/s3"
	xenditrepo "app/app/repository/xendit"
	usecase_agent "app/app/usecase/agent"
	usecase_webhook "app/app/usecase/webhook"
	"app/helpers"
	"context"

	usecase_member "app/app/usecase/member"
	usecase_superadmin "app/app/usecase/superadmin"

	cronjob "app/cron"

	"github.com/robfig/cron/v3"

	// "context"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	yureka_mongodb "github.com/Yureka-Teknologi-Cipta/yureka/services/mongodb"
	yureka_redis "github.com/Yureka-Teknologi-Cipta/yureka/services/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	timeoutStr := os.Getenv("TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "5"
	}
	timeout, _ := strconv.Atoi(timeoutStr)
	timeoutContext := time.Duration(timeout) * time.Second

	// logger
	writers := make([]io.Writer, 0)
	if logSTDOUT, _ := strconv.ParseBool(os.Getenv("LOG_TO_STDOUT")); logSTDOUT {
		writers = append(writers, os.Stdout)
	}

	if logFILE, _ := strconv.ParseBool(os.Getenv("LOG_TO_FILE")); logFILE {
		logMaxSize, _ := strconv.Atoi(os.Getenv("LOG_MAX_SIZE"))
		if logMaxSize == 0 {
			logMaxSize = 50 //default 50 megabytes
		}

		logFilename := os.Getenv("LOG_FILENAME")
		if logFilename == "" {
			logFilename = "server.log"
		}

		lg := &lumberjack.Logger{
			Filename:   logFilename,
			MaxSize:    logMaxSize,
			MaxBackups: 1,
			LocalTime:  true,
		}

		writers = append(writers, lg)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(io.MultiWriter(writers...))

	// set gin writer to logrus
	gin.DefaultWriter = logrus.StandardLogger().Writer()

	// init mongo database
	mongo := yureka_mongodb.Connect(timeoutContext, os.Getenv("MONGO_URL"), "")

	// init redis database
	var redisClient *redis.Client
	if useRedis, err := strconv.ParseBool(os.Getenv("USE_REDIS")); err == nil && useRedis {
		redisClient = yureka_redis.Connect(timeoutContext, os.Getenv("REDIS_URL"))
	}

	// init s3
	s3Repo := s3Repo.NewS3Repo()

	// init repo
	mongorepo := mongorepo.NewMongodbRepo(mongo)

	// xendit repo
	xenditRepo := xenditrepo.NewXenditRepo()

	// redis repo
	redisrepo := redisrepo.NewRedisRepo(redisClient)

	runType := os.Getenv("APP_RUNTYPE")
	if !helpers.InArrayString(runType, []string{"both", "cron", "api"}) {
		runType = "both"
	}

	if runType == "both" || runType == "cron" {
		// init cron
		c := cron.New(
			cron.WithChain(
				cron.SkipIfStillRunning(
					cron.VerbosePrintfLogger(logrus.StandardLogger()),
				),
				cron.Recover(
					cron.VerbosePrintfLogger(logrus.StandardLogger()),
				),
			),
		)

		// init cronjob
		cj := cronjob.NewCronjob(cronjob.RepoInjection{
			MongoDBRepo: mongorepo,
			Redis:       redisrepo,
			Ctx:         context.TODO(),
			Cron:        c,
		})
		cj.Run(runType == "both")
	}

	if runType == "both" || runType == "api" {
		// init usecase customer
		ucMember := usecase_member.NewAppUsecase(usecase_member.RepoInjection{
			MongoDBRepo: mongorepo,
			Redis:       redisrepo,
			S3Repo:      s3Repo,
			XenditRepo:  xenditRepo,
		}, timeoutContext)

		// init usecase agent
		ucAgent := usecase_agent.NewAppAgentUsecase(usecase_agent.RepoInjection{
			MongoDBRepo: mongorepo,
			Redis:       redisrepo,
			S3Repo:      s3Repo,
		}, timeoutContext)

		// init usecase superadmin
		ucSuperadmin := usecase_superadmin.NewAppSuperadminUsecase(usecase_superadmin.RepoInjection{
			MongoDBRepo: mongorepo,
			Redis:       redisrepo,
			S3Repo:      s3Repo,
		}, timeoutContext)

		// init usecase webhook
		ucWebhook := usecase_webhook.NewAppWebhookUsecase(usecase_webhook.RepoInjection{
			MongoDBRepo: mongorepo,
			Redis:       redisrepo,
		}, timeoutContext)

		// init middleware
		mdl := middleware.NewMiddleware(redisClient, mongorepo)

		// gin mode realease when go env is production
		if os.Getenv("GO_ENV") == "production" || os.Getenv("GO_ENV") == "prod" {
			gin.SetMode(gin.ReleaseMode)
		}

		// init gin
		ginEngine := gin.New()

		// add exception handler
		ginEngine.Use(mdl.Recovery())

		// add logger
		ginEngine.Use(mdl.Logger(io.MultiWriter(writers...)))

		// cors
		ginEngine.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"Content-Length"},
			MaxAge:           12 * time.Hour,
		}))

		// default route
		ginEngine.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, map[string]any{
				"message": "It works",
			})
		})

		// init route
		http_member.NewCustomerHandler(ginEngine.Group("/customer"), mdl, ucMember)
		http_agent.NewAgentRouteHandler(ginEngine.Group("/agent"), mdl, ucAgent)
		http_superadmin.NewSuperadminRouteHandler(ginEngine.Group("/superadmin"), mdl, ucSuperadmin)
		http_webhook.NewWebhookRouteHandler(ginEngine.Group(""), mdl, ucWebhook)

		port := os.Getenv("PORT")

		logrus.Infof("Service running on port %s", port)
		ginEngine.Run(":" + port)
	}
}
