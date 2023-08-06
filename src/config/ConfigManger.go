package config

import "fmt"

func InitConfigManger(path string) {

	if !LoadServerConfig(path) {
		panic(fmt.Sprintf("LoadServerConfig(path) error"))
	}

	if !LoadActivityNpcConfig(path) {
		panic(fmt.Sprintf("LoadActivityNpcConfig(path) error"))
	}

	if !LoadActivityInfoConfig(path) {
		panic(fmt.Sprintf("LoadActivityInfoConfig(path) error"))
	}

	if !LoadActivitypassawardConfig(path) {
		panic(fmt.Sprintf("LoadActivitypassawardConfig(path) error"))
	}

	if !LoadActivityNpcGroupConfig(path) {
		panic(fmt.Sprintf("LoadActivityNpcGroupConfig(path) error"))
	}

}
