package main

import (
	"flag"
	"fmt"
	"os"

	apiclient "github.com/justEstif/todo-open/internal/client/api"
	"github.com/justEstif/todo-open/internal/core"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "validate" {
		os.Exit(runValidate(os.Args[2:]))
	}
	os.Exit(runHealth(os.Args[1:]))
}

func runHealth(args []string) int {
	fs := flag.NewFlagSet("todoopen", flag.ContinueOnError)
	baseURL := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	client := apiclient.New(*baseURL)
	if err := client.Health(); err != nil {
		fmt.Fprintf(os.Stderr, "server health check failed: %v\n", err)
		return 1
	}

	fmt.Println("todo.open server is healthy")
	return 0
}

func runValidate(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	path := fs.String("file", "tasks.jsonl", "path to JSONL task file")
	mode := fs.String("mode", string(core.ValidationModeStrict), "validation mode: strict|compat")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	validationMode := core.ValidationMode(*mode)
	if validationMode != core.ValidationModeStrict && validationMode != core.ValidationModeCompat {
		fmt.Fprintf(os.Stderr, "invalid mode %q (expected strict|compat)\n", *mode)
		return 2
	}

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", *path, err)
		return 1
	}
	defer f.Close()

	issues, err := core.ValidateTaskJSONL(f, validationMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "validation failed: %v\n", err)
		return 1
	}
	if len(issues) == 0 {
		fmt.Printf("OK: %s is valid (%s mode)\n", *path, validationMode)
		return 0
	}

	fmt.Fprintf(os.Stderr, "Found %d validation issue(s) in %s:\n", len(issues), *path)
	for _, issue := range issues {
		fmt.Fprintf(os.Stderr, "- line %d field %s: %s\n", issue.Line, issue.Field, issue.Message)
		fmt.Fprintf(os.Stderr, "  context: %s\n", issue.Context)
	}
	return 1
}
