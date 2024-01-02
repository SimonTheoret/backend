package back

import (
	"flag"
	"fmt"
	"os"
)

// This fonctions sets the log file and the port with the help of the CLI. It
// has default port 3000 and default logFile server.log
func SetUp() (logFile *os.File, address string) {
	port := flag.Int("port", 3000, "Specify the port")
	log := flag.String("Log destination", "server.log", "Specify the log file destination")
	flag.Parse()

	address = fmt.Sprintf("0.0.0.0:%d", *port) // use address 0.0.0.0:port
	f, err := os.Create(*log)
	if err != nil {
		fmt.Println(err) // print err if logging to file is impossible
		os.Exit(1)       //Stops program
	}
	return f, address
}
