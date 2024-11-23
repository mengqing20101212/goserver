package config

import (
	"fmt"
	"logger"
)

var log *logger.Logger

func InitConfigManger(logger *logger.Logger, path string) {
	if log == nil {
		log = logger
	}

	if !LoadActivitypassawardConfig(path) {
		panic(fmt.Sprintf("LoadActivitypassawardConfig(path) error"))
	}

	if !LoadActivityNpcGroupConfig(path) {
		panic(fmt.Sprintf("LoadActivityNpcGroupConfig(path) error"))
	}

	if !LoadServerConfig(path) {
		panic(fmt.Sprintf("LoadServerConfig(path) error"))
	}

	if !LoadActivityNpcConfig(path) {
		panic(fmt.Sprintf("LoadActivityNpcConfig(path) error"))
	}

	if !LoadActivityInfoConfig(path) {
		panic(fmt.Sprintf("LoadActivityInfoConfig(path) error"))
	}

}
