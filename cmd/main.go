package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"stpManager/controller"
	"stpManager/repo"

	"time"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("hello from stp-manager")

	if err := initConfig(); err != nil {
		logrus.Fatal("coud not load config file")
	}

	logrus.Debug("creating postgres repo")
	repo, err := repo.NewPostgresRepo(&repo.PostgresConfig{
		Host:     viper.GetString("repo.host"),
		Port:     viper.GetString("repo.port"),
		Username: viper.GetString("repo.username"),
		Password: viper.GetString("repo.password"),
		DBName:   viper.GetString("repo.dbname"),
		SSLMode:  viper.GetString("repo.sslmode"),
		SchemaName: viper.GetString("repo.schema"),
	})

	if err != nil {
		logrus.Fatal("could not connect to postgres: ", err.Error())
	}
	logrus.Debug("created postgres repo")

	c := controller.NewController(repo)

	server := &http.Server{
		Addr:           ":" + viper.GetString("port"),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		Handler:        c.InitRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.Fatal("could not start http server: ", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("/Users/aryabov/projects/stp/manager/config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
