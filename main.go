package main

import (
	"fmt"
	"log"

	"github.com/windyzoe/study-house/cron"
	"github.com/windyzoe/study-house/info"
	"github.com/windyzoe/study-house/rest"
	"github.com/windyzoe/study-house/service"
)

func main() {
	cron.Start()
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	fmt.Println("123123")
	service.Start()
	info.Start()
	rest.Start()

}
