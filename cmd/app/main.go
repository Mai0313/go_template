package main

import (
	"fmt"
	"log"

	"go-template/core/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Go Template App v%s\n", cfg.Version)
	fmt.Println("This is a template application.")
	fmt.Println("Replace this with your actual application logic.")
}
