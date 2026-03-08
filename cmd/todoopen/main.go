package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/justEstif/todo-open/internal/app"
	apiclient "github.com/justEstif/todo-open/internal/client/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/info"
	"github.com/justEstif/todo-open/internal/tui"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		return runHelp(stdout)
	}

	switch args[0] {
	case "-h", "--help", "help":
		return runHelp(stdout)
	case "-v", "--version", "version":
		return runVersion(stdout)
	case "--agent-info", "-A":
		return runAgentInfo(args[1:], stdout, stderr)
	case "validate":
		return runValidate(args[1:], stdout, stderr)
	case "task":
		return runTask(args[1:], stdout, stderr)
	case "adapters":
		return runAdapters(args[1:], stdout, stderr)
	case "web", "gui":
		return runWeb(args[1:], stdout, stderr)
	case "tui":
		return runTui(args[1:], stdout, stderr)
	default:
		return runHealth(args, stdout, stderr)
	}
}

func runHelp(stdout io.Writer) int {
	fmt.Fprintln(stdout, "todoopen - server-first local task client")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "Usage:")
	fmt.Fprintln(stdout, "  todoopen --help")
	fmt.Fprintln(stdout, "  todoopen --version")
	fmt.Fprintln(stdout, "  todoopen --agent-info [--server URL]   # print agent-info JSON and exit")
	fmt.Fprintln(stdout, "  todoopen [--server URL]                # health check")
	fmt.Fprintln(stdout, "  todoopen web [--addr ADDR] [--no-open] # launch web app")
	fmt.Fprintln(stdout, "  todoopen tui [--addr ADDR] [--server URL] # launch terminal UI")
	fmt.Fprintln(stdout, "  todoopen validate [flags]")
	fmt.Fprintln(stdout, "  todoopen task <create|list|get|update|delete> [flags]")
	fmt.Fprintln(stdout, "  todoopen adapters [--workspace PATH] [--json]")
	return 0
}

func runVersion(stdout io.Writer) int {
	fmt.Fprintf(stdout, "todoopen %s\n", resolvedVersion())
	return 0
}

func runAgentInfo(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("agent-info", flag.ContinueOnError)
	fs.SetOutput(stderr)
	baseURL := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	printJSON(stdout, info.Build(resolvedVersion(), *baseURL))
	return 0
}

func resolvedVersion() string {
	if version != "" && version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
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

func runWeb(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("web", flag.ContinueOnError)
	fs.SetOutput(stderr)
	addr := fs.String("addr", "127.0.0.1:8080", "address to bind local server")
	baseURL := fs.String("server", "", "use an existing server URL instead of starting one")
	noOpen := fs.Bool("no-open", false, "do not open browser automatically")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	url := *baseURL
	if url == "" {
		url = "http://" + *addr
		srv, err := app.NewServer(*addr)
		if err != nil {
			fmt.Fprintf(stderr, "server setup failed: %v\n", err)
			return 1
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(stderr, "server failed: %v\n", err)
			}
		}()
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
		}()
	}

	if err := waitForHealthy(url, 5*time.Second); err != nil {
		fmt.Fprintf(stderr, "web launch failed: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "todo.open web is available at %s/\n", url)
	if !*noOpen {
		if err := openBrowser(url + "/"); err != nil {
			fmt.Fprintf(stderr, "failed to open browser automatically: %v\n", err)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	return 0
}

func runTui(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("tui", flag.ContinueOnError)
	fs.SetOutput(stderr)
	addr := fs.String("addr", "127.0.0.1:8080", "address to bind local server")
	baseURL := fs.String("server", "", "use an existing server URL instead of starting one")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	url := *baseURL
	if url == "" {
		url = "http://" + *addr
		srv, err := app.NewServer(*addr)
		if err != nil {
			fmt.Fprintf(stderr, "server setup failed: %v\n", err)
			return 1
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Fprintf(stderr, "server failed: %v\n", err)
			}
		}()
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = srv.Shutdown(ctx)
		}()
	}

	if err := waitForHealthy(url, 5*time.Second); err != nil {
		fmt.Fprintf(stderr, "tui launch failed: %v\n", err)
		return 1
	}

	if err := tui.Run(url); err != nil {
		fmt.Fprintf(stderr, "tui error: %v\n", err)
		return 1
	}
	return 0
}

func waitForHealthy(baseURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := apiclient.New(baseURL)
	for {
		if err := client.Health(); err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for healthy server at %s", baseURL)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func runAdapters(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("adapters", flag.ContinueOnError)
	fs.SetOutput(stderr)
	workspace := fs.String("workspace", "", "workspace root path (defaults to TODOOPEN_WORKSPACE_ROOT or cwd)")
	asJSON := fs.Bool("json", false, "print adapter status as JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	workspaceRoot := *workspace
	if workspaceRoot == "" {
		workspaceRoot = os.Getenv("TODOOPEN_WORKSPACE_ROOT")
	}
	if workspaceRoot == "" {
		var err error
		workspaceRoot, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(stderr, "failed to resolve workspace root: %v\n", err)
			return 1
		}
	}

	if _, err := app.LoadWorkspaceMeta(workspaceRoot); err != nil {
		fmt.Fprintf(stderr, "failed to load workspace metadata: %v\n", err)
		return 1
	}
	adapterCfg, err2 := app.LoadAdapterFileConfig(workspaceRoot)
	if err2 != nil {
		fmt.Fprintf(stderr, "failed to load adapter config: %v\n", err2)
		return 1
	}
	viewRegistry, err := app.NewViewRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "failed to load view adapters: %v\n", err)
		return 1
	}
	syncRegistry, err := app.NewSyncRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "failed to load sync adapters: %v\n", err)
		return 1
	}

	runtime := app.BuildAdapterRuntimeFromConfig(context.Background(), adapterCfg, viewRegistry, syncRegistry) //nolint:contextcheck
	if *asJSON {
		printJSON(stdout, runtime)
	} else {
		for _, s := range runtime.Status {
			health := "healthy"
			if !s.Healthy {
				health = "unhealthy"
			}
			fmt.Fprintf(stdout, "%s\t%s\tsource=%s\tenabled=%t\t%s\n", s.Kind, s.Name, s.Source, s.Enabled, health)
			if s.Message != "" {
				fmt.Fprintf(stdout, "  %s\n", s.Message)
			}
		}
		if runtime.Ready {
			fmt.Fprintln(stdout, "ready=true")
		} else {
			fmt.Fprintln(stdout, "ready=false")
			for _, err := range runtime.Errors {
				fmt.Fprintf(stdout, "error: %s\n", err)
			}
		}
	}

	if !runtime.Ready {
		return 1
	}
	return 0
}

type taskSubcommand func(args []string, stdout io.Writer, stderr io.Writer) int

func runTask(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: todoopen task <create|list|get|update|delete> [flags]")
		return 2
	}

	commands := map[string]taskSubcommand{
		"create": runTaskCreate,
		"list":   runTaskList,
		"get":    runTaskGet,
		"update": runTaskUpdate,
		"delete": runTaskDelete,
	}

	handler, ok := commands[args[0]]
	if !ok {
		fmt.Fprintf(stderr, "unknown task subcommand %q\n", args[0])
		return 2
	}
	return handler(args[1:], stdout, stderr)
}

func runTaskCreate(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("task create", flag.ContinueOnError)
	fs.SetOutput(stderr)
	server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	title := fs.String("title", "", "task title")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	task, err := apiclient.New(*server).CreateTask(*title)
	if err != nil {
		fmt.Fprintf(stderr, "create failed: %v\n", err)
		return 1
	}
	printJSON(stdout, task)
	return 0
}

func runTaskList(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("task list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	if err := fs.Parse(args); err != nil {
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
}

func runTaskGet(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("task get", flag.ContinueOnError)
	fs.SetOutput(stderr)
	server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	id := fs.String("id", "", "task id")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	task, err := apiclient.New(*server).GetTask(*id)
	if err != nil {
		fmt.Fprintf(stderr, "get failed: %v\n", err)
		return 1
	}
	printJSON(stdout, task)
	return 0
}

func runTaskUpdate(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("task update", flag.ContinueOnError)
	fs.SetOutput(stderr)
	server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	id := fs.String("id", "", "task id")
	title := fs.String("title", "", "task title")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	task, err := apiclient.New(*server).UpdateTask(*id, *title)
	if err != nil {
		fmt.Fprintf(stderr, "update failed: %v\n", err)
		return 1
	}
	printJSON(stdout, task)
	return 0
}

func runTaskDelete(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("task delete", flag.ContinueOnError)
	fs.SetOutput(stderr)
	server := fs.String("server", "http://127.0.0.1:8080", "todo.open server base URL")
	id := fs.String("id", "", "task id")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if err := apiclient.New(*server).DeleteTask(*id); err != nil {
		fmt.Fprintf(stderr, "delete failed: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, "deleted")
	return 0
}

func printJSON(w io.Writer, v any) {
	_ = json.NewEncoder(w).Encode(v)
}
