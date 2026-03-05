package main

import (
	"flag"
	"fmt"
	"os"

	apiclient "github.com/justEstif/todo-open/internal/client/api"
)

func main() {
	baseURL := flag.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	flag.Parse()

	client := apiclient.New(*baseURL)
	if err := client.Health(); err != nil {
		fmt.Fprintf(os.Stderr, "server health check failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("todo.open server is healthy")
}
