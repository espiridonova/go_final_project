package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540

func Run() error {
	port := defaultPort
	portParam := os.Getenv("TODO_PORT")
	if portParam != "" {
		var err error
		port, err = strconv.Atoi(portParam)
		if err != nil {
			return err
		}
	}

	fmt.Printf("server start; port:%d\n", port)

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
