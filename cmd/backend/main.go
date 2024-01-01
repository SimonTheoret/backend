package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Command line arguments
	port := flag.Int("port", 3000, "Specify the port")
	log := flag.String("Log destination", "server.log", "Specify the log file destination")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port) // use address 0.0.0.0:port
	f, err := os.Create(*log)
	if err != nil {
		fmt.Println(err) // print err if logging to file is impossible
		os.Exit(1)       //Stops program
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // write logs to file f

	r := gin.Default()
	r.POST("predict", ModelPredict)
	r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")}
}
