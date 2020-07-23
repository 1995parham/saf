package service

import (
	"time"

	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New()

	e.GET("/", Job)
	e.Logger.Fatal(e.Start(":1373"))
}

func Job(c echo.Context) error {
	time.Sleep(5 * time.Second)

	return nil
}