package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-template/core/config"
)

func main() {
	var (
		version = flag.Bool("version", false, "Show version information")
		help    = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *version {
		fmt.Printf("Go Template CLI v%s\n", cfg.Version)
		return
	}

	fmt.Printf("Go Template CLI v%s\n", cfg.Version)
	fmt.Println("This is a template CLI application.")
	fmt.Println("Add your CLI commands and functionality here.")
}

func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}
