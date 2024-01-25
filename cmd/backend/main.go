package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"

	back "github.com/SimonTheoret/backend/pkg"
)

func main() {
	// flags, logs, and gin
	logFile, address := back.SetUp()
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout) // write logs to file f
	rf := back.DefaultFormatter()
	r := gin.Default()

	// Model(s) and modelMapper
<<<<<<< HEAD
	model := back.NewHttpModel("TestModel", 0, "127.0.0.1")
	mapper := back.SetUpModels([]back.modeler{model}, rf) // Build the models and start them
=======
    h := back.HttpModel{Dest : "127.0.0.1" }
	model := back.NewBase("TestModel", &h)
	mapper := back.SetUpModels([]back.Modeler{&model}, rf) // Build the models and start them
>>>>>>> 4cf7017 (added ability to add model at runtime)

	// Router
	r.POST("/post", back.HandlerModelPost(mapper))
	r.GET("/get", back.HandlerModelGet(mapper))

    err := r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")}
    fmt.Println(err)
}
