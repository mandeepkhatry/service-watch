package main

import (
	"log"
	"service-watch/watch"
)

var configPath = "config/watch.json"

func main() {

	serviceWatch, err := watch.NewServiceWatcher(configPath)

	if err != nil {
		log.Fatalln(err)
	}

	err = serviceWatch.Watch()

	if err != nil {
		log.Fatalln(err)
	}

}
