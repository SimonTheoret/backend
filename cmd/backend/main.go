package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	p := flag.Int("port", 3000, "Specify the port")
	flag.Parse()
	port := fmt.Sprintf(":%d", *p)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(port))
}
