package main

import (
	"io"
	"os"

	back "github.com/SimonTheoret/backend/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	// flags, logs, and gin
	logFile, address := back.SetUp()
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout) // write logs to file f
	rf := back.DefaultFormatter()
	r := gin.Default()

	// Model(s) and modelMapper
	model := back.NewHttpModel("TestModel", 0, "127.0.0.1")
	mapper := back.SetUpModels([]back.modeler{model}, rf) // Build the models and start them

	// Router
	r.POST("/post", back.HandlerModelPost(mapper))
	r.GET("/get", back.HandlerModelGet(mapper))

	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")}
}
