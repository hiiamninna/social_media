package main

import (
	"fmt"
	"os"
	"os/signal"
	"social_media/config"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:01:02 "

func main() {
	app := config.AppRoute()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println(time.Now().Format(TIME_FORMAT), "Gracefully shutting down...")
		_ = app.Shutdown()
	}()
	if err := app.Listen(":8080"); err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), "error on http listens : "+err.Error())
	}
}
