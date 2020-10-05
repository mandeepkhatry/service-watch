package main

import (
	"flag"
	"log"
	"service-watch/watch"
)

func main() {

	configPath := flag.String("config", "", "config path")

	flag.Parse()

	serviceWatch, err := watch.NewServiceWatcher(string(*configPath))

	if err != nil {
		log.Fatalln(err)
	}

	err = serviceWatch.Watch()

	if err != nil {
		log.Fatalln(err)
	}

}
