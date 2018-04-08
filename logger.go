package web

import (
	log "github.com/VectorsOrigin/logger"
)

/*
	logger 负责Web框架的日志打印
	@不提供给其他程序使用
*/
var (
	logger = log.NewLogger("")
)

func init() {
}

func Logger() *log.TLogger {
	return logger
}
