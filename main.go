package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rezaAmiri123/scalable-backend/internal/cache"
	"github.com/rezaAmiri123/scalable-backend/internal/controller"
	"github.com/rezaAmiri123/scalable-backend/internal/database"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
)

func main() {
	// loading .env file for the local development
	_ = godotenv.Load()
	// setup mysql connections
	masterDSN := os.Getenv("MASTER_DSN")
	replicaDSN := strings.Split(os.Getenv("REPLICA_DSN"), ",")
	replicaDSN = lo.Map(replicaDSN, func(item string, _ int) string {
		return strings.TrimSpace(item)
	})
	db, err := database.NewGorm(strings.TrimSpace(masterDSN), replicaDSN...)
	if err != nil {
		logrus.WithError(err).Panicln("the database connection setup failed")
	}
	// set up the database
	gdb := database.NewGormDatabase(db)
	if err := gdb.Migrate(); err != nil {
		logrus.WithError(err).Panicln("error while migrating the database")
	}

	// set up the cache
	rdb := cache.NewRedis(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"))
	defer rdb.Close()
	synchronize := cache.NewSync(rdb, gdb)
	_ = synchronize
	rc := cache.NewRedisCache(rdb)
	inmem := cache.NewInMemoryCache(rc)

	// set up Prometheus exposer
	http.Handle("/metric", promhttp.Handler())
	logrus.Info("starting the metric server on port 8081")
	go func() {
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			logrus.WithError(err).Error("the metric server has stopped")
		}
	}()
	// set up the http apis
	e := echo.New()
	controller.NewEchoController(e, gdb, rc, inmem)
	logrus.Info("starting the api server on port 8080")
	logrus.WithError(err).Error(e.Start(":8080"))
}
