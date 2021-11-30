package main

import (
	"log"

	"github.com/windyzoe/study-house/info"
	"github.com/windyzoe/study-house/service"
)

func main() {
	// cron.Start()
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
	service.Start()
	defer service.DB.Close()
	info.Start()
	// rest.Start()
}
