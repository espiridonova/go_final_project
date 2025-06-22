package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540

func Run() error {
	println("server start")

	port := defaultPort
	portParam := os.Getenv("TODO_PORT")
	if portParam != "" {
		var err error
		port, err = strconv.Atoi(portParam)
		if err != nil {
			return err
		}
	}

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
