package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Infinite-Sum-Games/Am.EventKit/bootstrap"
)

var (
	fs   = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	mode = fs.String("mode", "api", "run mode: `api` or `mcp`")
)

func main() {
	err := fs.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	app := bootstrap.NewApp()

	switch *mode {
	case "api":
		err = app.RunGinServer()
		if err != nil {
			fmt.Printf("error running gin server: %v", err)
		}
		app.Close()
	case "mcp":
		err = app.RunMCPServer(ctx)
		if err != nil {
			fmt.Printf("error running MCP server: %v", err)
			os.Exit(1)
		}
		app.Close()
	default:
		fmt.Printf("unknown mode: %s", *mode)
		flag.Usage()
		os.Exit(1)
	}
}
