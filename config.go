package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DSN  string
	Port string
}

var ErrDbDsnNotSet = errors.New("could not find DB in env vars")

func configFromEnv() (cfg config, err error) {
	err = godotenv.Load()
	if err != nil {
		fmt.Errorf("Loading .env file: %w", err)
	}
	cfg.DSN = os.Getenv("DSN")
	cfg.Port = os.Getenv("PORT")
	if cfg.DSN == "" {
		err = ErrDbDsnNotSet
		return
	}
	return
}
