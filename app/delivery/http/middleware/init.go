package middleware

import (
	mongorepo "app/app/repository/mongo"
	"app/helpers"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type appMiddleware struct {
	secretKeyCustomer   string
	secretKeyAgent      string
	secretKeySuperuser  string
	secretKeySuperadmin string
	xenditWebhookToken  string
	cache               CacheConfig
	mongo               mongorepo.MongoDBRepo
}

type CacheConfig struct {
	enabled     bool
	headerKeys  []string
	store       *redis.Client
	storeTTL    time.Duration
	cachePrefix string
}

func NewMiddleware(redis *redis.Client, mongo mongorepo.MongoDBRepo) Middleware {
	ttl, _ := time.ParseDuration(os.Getenv("REDIS_TTL"))
	// default ttl redis
	if ttl == 0 {
		ttl = 1 * time.Minute
	}

	useRedis, _ := strconv.ParseBool(os.Getenv("USE_REDIS"))
	redisKeyPrefix := os.Getenv("REDIS_KEY_PREFIX")

	return &appMiddleware{
		mongo:               mongo,
		secretKeyCustomer:   helpers.GetJWTSecretKeyCustomer(),
		secretKeyAgent:      helpers.GetJWTSecretKeyAgent(),
		secretKeySuperuser:  helpers.GetJWTSecretKeySuperuser(),
		secretKeySuperadmin: helpers.GetJWTSecretKeySuperadmin(),
		xenditWebhookToken:  os.Getenv("XENDIT_WEBHOOK_VERIFICATION_TOKEN"),
		cache: CacheConfig{
			enabled:     useRedis,
			store:       redis,
			cachePrefix: redisKeyPrefix + "gin:",
			storeTTL:    ttl,
			headerKeys: []string{
				"User-Agent",
				"Accept",
				"Accept-Encoding",
				"Accept-Language",
				"Cookie",
			},
		},
	}
}

type Middleware interface {
	AuthCustomer() gin.HandlerFunc
	AuthAgent() gin.HandlerFunc
	AuthSuperadmin() gin.HandlerFunc
	Role(allowedRoles ...string) gin.HandlerFunc
	Logger(writer io.Writer) gin.HandlerFunc
	Recovery() gin.HandlerFunc
	Cache(expiry ...time.Duration) gin.HandlerFunc
	VerifyXenditWebhookToken() gin.HandlerFunc
}
