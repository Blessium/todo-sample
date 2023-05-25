package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New() 
    RegisterRoutes(e)
    e.Logger.Fatal(e.Start(":1234"))
}
