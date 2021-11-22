package cron

import (
	"fmt"

	"github.com/robfig/cron"
)

func Start() {
	c := cron.New()
	c.AddFunc("*/3 * * * * *", func() {
		fmt.Println("every 3 seconds executing")
	})
	go c.Start()
}
