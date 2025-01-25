package config

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tutup-lapak/db"

	"github.com/sirupsen/logrus"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func NewDatabase(log *logrus.Logger) *db.Postgres {
	dbUrl := os.Getenv("DATABASE_URL")
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	db, err := db.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatal("unable connect to db", err.Error())
		os.Exit(1)
	}

	return db
}
