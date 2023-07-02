package main

import "goserver/common/logger"

func main() {
	log := logger.Logger{}
	dir := "../logs"
	filename := "test.log"
	log.Init(&dir, &filename)
}
