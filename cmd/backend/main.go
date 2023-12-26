package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Command line arguments
	port := flag.Int("port", 3000, "Specify the port")
	log := flag.String("Log destination", "server.log", "Specify the log file destination")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port) // use address 0.0.0.0:port
	f, _ := os.Create(*log)
	gin.DefaultWriter = io.MultiWriter(f) // write logs to file f

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")}
}
