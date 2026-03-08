package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/justEstif/todo-open/internal/tui"
)

func main() {
	addr := flag.String("server", "http://localhost:8080", "todo.open server address")
	flag.Parse()

	if err := tui.Run(*addr); err != nil {
		fmt.Fprintln(os.Stderr, "tui error:", err)
		os.Exit(1)
	}
}
