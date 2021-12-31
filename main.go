package main

import (
	"io"
	"log"
	"os"

	"github.com/windyzoe/study-house/db"
	"github.com/windyzoe/study-house/rest"
	"github.com/windyzoe/study-house/spider"
)

func init() {
	file := "./" + "message" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
}

func main() {
	// cron.Start()
	db.Start()
	defer db.DB.Close()
	spider.Start()
	rest.Start()

}
