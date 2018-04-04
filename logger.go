package web

import (
	log "github.com/VectorsOrigin/logger"
)

var (
	logger = log.NewLogger("")
)

func init() {
}

func Logger() *log.TLogger {
	return logger
}
