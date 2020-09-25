package main

import (
	"fmt"
	"service-watch/watch"
)

var configPath = "config/watch.json"

func main() {
	fmt.Println("Testing")

	serviceWatch, err := watch.NewServiceWatcher(configPath)

	if err != nil {
		panic(err)
	}

	err = serviceWatch.Watch()

	if err != nil {
		panic(err)
	}

}
