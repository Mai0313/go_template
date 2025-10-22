package main

import (
	"flag"
	"fmt"

	"go-template/core/version"
)

func main() {
	var showVersion = flag.Bool("version", false, "Show version information")

	// Parse flags before use
	flag.Parse()

	// If --version flag is set, print version and exit
	if *showVersion {
		versionInfo := version.Get()
		fmt.Printf("Version: %s\n", versionInfo.Version)
		fmt.Printf("Build Time: %s\n", versionInfo.BuildTime)
		fmt.Printf("Git Commit: %s\n", versionInfo.GitCommit)
		fmt.Printf("Go Version: %s\n", versionInfo.GoVersion)
		return
	}
}
