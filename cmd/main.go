package main

import (
	"fmt"
	"log"
	"service-watch/watch"
)

func main() {
	fmt.Println("Running main")

	watcher := watch.ServiceWatcher{}
	watcher.Init("static/gateway-config/test.toml")
	err := watcher.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(watcher.AppConfig.Api)
}
