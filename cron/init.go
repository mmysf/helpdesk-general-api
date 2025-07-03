package cronjob

// import (
// 	"context"
// 	"time"

// 	"github.com/robfig/cron/v3"
// 	"github.com/sirupsen/logrus"

// 	mongorepo "app/app/repository/mongo"
// 	redisrepo "app/app/repository/redis"
// 	s3repo "app/app/repository/s3"
// 	xenditrepo "app/app/repository/xendit"
// )

// type CronjobHandler interface {
// 	Run(runInBackground bool)
// }

// type cronjob struct {
// 	ctx            context.Context
// 	cron           *cron.Cron
// 	mongodbRepo    mongorepo.MongoDBRepo
// 	contextTimeout time.Duration
// 	redisRepo      redisrepo.RedisRepo
// 	s3Repo         s3repo.S3Repo
// 	xenditRepo     xenditrepo.XenditRepo
// }

// func NewCronjob(cj *cronjob, logger *logrus.Logger) CronjobHandler {
// 	// init cron
// 	c := cron.New(
// 		cron.WithChain(
// 			cron.SkipIfStillRunning(
// 				cron.VerbosePrintfLogger(logger),
// 			),
// 			cron.Recover(
// 				cron.VerbosePrintfLogger(logger),
// 			),
// 		),
// 	)

// 	cj := cronjob{
// 		ctx:  context.TODO(),
// 		cron: c,
// 	}

// 	return &cj
// }

// func (cj *cronjob) Run(runInBackground bool) {
// 	// starting cron
// 	logrus.Info("Cronjob started")
// 	if runInBackground {
// 		cj.cron.Start() // run with goroutine
// 	} else {
// 		cj.cron.Run()
// 	}
// }

import (
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	"context"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type cronjob struct {
	ctx         context.Context
	cron        *cron.Cron
	mongodbRepo mongorepo.MongoDBRepo
	redisRepo   redisrepo.RedisRepo
}

type RepoInjection struct {
	Ctx         context.Context
	Cron        *cron.Cron
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
}

func NewCronjob(r RepoInjection) CronjobHandler {
	return &cronjob{
		ctx:         r.Ctx,
		cron:        r.Cron,
		mongodbRepo: r.MongoDBRepo,
		redisRepo:   r.Redis,
	}
}

type CronjobHandler interface {
	Run(runInBackground bool)
}

func (cj *cronjob) Run(runInBackground bool) {
	cj.SyncExpiredSubscription()
	cj.AutoCloseResolvedTickets()


	// starting cron
	logrus.Info("Cronjob started")
	if runInBackground {
		cj.cron.Start() // run with goroutine
	} else {
		cj.cron.Run()
	}
}
