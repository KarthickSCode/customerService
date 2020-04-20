package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer,"Welcome to the ERPLY customer API")
	})

	neg := negroni.Classic()
	neg.UseHandler(mux)

	neg.Run(":8080")
}
