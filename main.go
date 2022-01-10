package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/windyzoe/study-house/db"
	"github.com/windyzoe/study-house/rest"
	"github.com/windyzoe/study-house/spider"
	"github.com/windyzoe/study-house/util"
)

// func init() {
// 	file := "./" + "message" + ".log"
// 	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
// 	if err != nil {
// 		panic(err)
// 	}
// 	mw := io.MultiWriter(os.Stdout, logFile)
// 	log.SetOutput(mw)
// 	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
// }

func initLog() {
	//global
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.TimeFieldFormat = time.RFC3339
	//console
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
	consoleWriter.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	consoleWriter.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	// file
	file := "./" + "message" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	//multi
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
	log.Info().Msg("Hello World")
}

func main() {
	initLog()
	util.InitConfig()
	// cron.Start()
	db.Start()
	defer db.DB.Close()
	spider.Start()
	rest.Start()
}
