package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	apiclient "github.com/justEstif/todo-open/internal/client/api"
	"github.com/justEstif/todo-open/internal/core"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "validate":
			return runValidate(args[1:], stdout, stderr)
		case "task":
			return runTask(args[1:], stdout, stderr)
		}
	}
	return runHealth(args, stdout, stderr)
}

func runHealth(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("todoopen", flag.ContinueOnError)
	fs.SetOutput(stderr)
	baseURL := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	client := apiclient.New(*baseURL)
	if err := client.Health(); err != nil {
		fmt.Fprintf(stderr, "server health check failed: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, "todo.open server is healthy")
	return 0
}

func runValidate(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	path := fs.String("file", "tasks.jsonl", "path to JSONL task file")
	mode := fs.String("mode", string(core.ValidationModeStrict), "validation mode: strict|compat")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	validationMode := core.ValidationMode(*mode)
	if validationMode != core.ValidationModeStrict && validationMode != core.ValidationModeCompat {
		fmt.Fprintf(stderr, "invalid mode %q (expected strict|compat)\n", *mode)
		return 2
	}

	f, err := os.Open(*path)
	if err != nil {
		fmt.Fprintf(stderr, "failed to open %s: %v\n", *path, err)
		return 1
	}
	defer f.Close()

	issues, err := core.ValidateTaskJSONL(f, validationMode)
	if err != nil {
		fmt.Fprintf(stderr, "validation failed: %v\n", err)
		return 1
	}
	if len(issues) == 0 {
		fmt.Fprintf(stdout, "OK: %s is valid (%s mode)\n", *path, validationMode)
		return 0
	}

	fmt.Fprintf(stderr, "Found %d validation issue(s) in %s:\n", len(issues), *path)
	for _, issue := range issues {
		fmt.Fprintf(stderr, "- line %d field %s: %s\n", issue.Line, issue.Field, issue.Message)
		fmt.Fprintf(stderr, "  context: %s\n", issue.Context)
	}
	return 1
}

func runTask(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: todoopen task <create|list|get|update|delete> [flags]")
		return 2
	}

	subcmd := args[0]
	subArgs := args[1:]
	switch subcmd {
	case "create":
		fs := flag.NewFlagSet("task create", flag.ContinueOnError)
		fs.SetOutput(stderr)
		server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
		title := fs.String("title", "", "task title")
		if err := fs.Parse(subArgs); err != nil {
			return 2
		}
		task, err := apiclient.New(*server).CreateTask(*title)
		if err != nil {
			fmt.Fprintf(stderr, "create failed: %v\n", err)
			return 1
		}
		printJSON(stdout, task)
		return 0

	case "list":
		fs := flag.NewFlagSet("task list", flag.ContinueOnError)
		fs.SetOutput(stderr)
		server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
		if err := fs.Parse(subArgs); err != nil {
			return 2
		}
		tasks, err := apiclient.New(*server).ListTasks()
		if err != nil {
			fmt.Fprintf(stderr, "list failed: %v\n", err)
			return 1
		}
		for _, task := range tasks {
			fmt.Fprintf(stdout, "%s\t%s\t%s\n", task.ID, task.Status, task.Title)
		}
		return 0

	case "get":
		fs := flag.NewFlagSet("task get", flag.ContinueOnError)
		fs.SetOutput(stderr)
		server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
		id := fs.String("id", "", "task id")
		if err := fs.Parse(subArgs); err != nil {
			return 2
		}
		task, err := apiclient.New(*server).GetTask(*id)
		if err != nil {
			fmt.Fprintf(stderr, "get failed: %v\n", err)
			return 1
		}
		printJSON(stdout, task)
		return 0

	case "update":
		fs := flag.NewFlagSet("task update", flag.ContinueOnError)
		fs.SetOutput(stderr)
		server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
		id := fs.String("id", "", "task id")
		title := fs.String("title", "", "task title")
		if err := fs.Parse(subArgs); err != nil {
			return 2
		}
		task, err := apiclient.New(*server).UpdateTask(*id, *title)
		if err != nil {
			fmt.Fprintf(stderr, "update failed: %v\n", err)
			return 1
		}
		printJSON(stdout, task)
		return 0

	case "delete":
		fs := flag.NewFlagSet("task delete", flag.ContinueOnError)
		fs.SetOutput(stderr)
		server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
		id := fs.String("id", "", "task id")
		if err := fs.Parse(subArgs); err != nil {
			return 2
		}
		if err := apiclient.New(*server).DeleteTask(*id); err != nil {
			fmt.Fprintf(stderr, "delete failed: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, "deleted")
		return 0
	default:
		fmt.Fprintf(stderr, "unknown task subcommand %q\n", subcmd)
		return 2
	}
}

func printJSON(w io.Writer, v any) {
	_ = json.NewEncoder(w).Encode(v)
}
