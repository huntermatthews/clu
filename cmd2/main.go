package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/huntermatthews/clu/cmd2/global"
)

var version = "unset"

func main() {
	// Create a custom handler with the Debug level
	handler := &CustomHandler{level: slog.LevelDebug}

	// Replace the default logger
	slog.SetDefault(slog.New(handler))

	slog.Debug("debug logging is enabled")
	slog.Debug("version: " + version)

	// ----------

	appInfo := global.GetAppInfo(version)
	buildInfo := global.GetBuildInfo()

	if buildInfo == nil {
		fmt.Println("Build info not available")
		return
	}

	appInfo.Print(os.Stdout)
	buildInfo.Print(os.Stdout)
}
