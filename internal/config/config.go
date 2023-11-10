package config

import (
	logger "log"
	"strconv"
	"time"

	"github.com/bamdadam/Minefield-Scavenger/internal/db/psql"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Psql psql.Config
}

func New() *Config {
	k := koanf.New(".")
	// envProvider := env.Provider("", ".", func(s string) string {
	// 	return s
	// })
	// if err := k.Load(envProvider, nil); err != nil {
	// 	logger.Fatal("error while loading enviroment variables " + err.Error())
	// }
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		logger.Fatal("error while loading enviroment variables: " + err.Error())
	}
	cfg := new(Config)
	cfg.Psql.ConnectionString = k.String("DB_CONN_STRING")
	d, err := strconv.ParseInt(k.String("DB_TIME_OUT"), 10, 64)
	if err != nil {
		logger.Fatal("error while parsing time out to int: " + err.Error())
	}
	cfg.Psql.Timeout = time.Duration(d) * time.Second
	return cfg
}
