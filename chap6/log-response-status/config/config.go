package config

import (
	"io"
	"log"
)

type AppConfig struct {
	Logger *log.Logger
}

func InitConfig(w io.Writer) AppConfig {
	return AppConfig{
		Logger: log.New(w, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

}
