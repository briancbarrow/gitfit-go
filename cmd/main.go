package main

import (
	"flag"

	"github.com/briancbarrow/gitfit-go/internal/server"
)

func main() {
	isProd := flag.Bool("isProd", false, "Determines if the server is in production mode")

	server := server.NewServer(*isProd)
	err := server.ListenAndServe()
	if err != nil {
		panic("cannot start server")
	}
}
