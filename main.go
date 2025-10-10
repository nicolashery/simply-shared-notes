package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/nicolashery/simply-shared-notes/app"
)

//go:embed sql/pragmas.sql
var pragmasSQL string

//go:embed all:dist
var distFS embed.FS

//go:embed locales/*.toml
var localesFS embed.FS

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, distFS, pragmasSQL, localesFS); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
