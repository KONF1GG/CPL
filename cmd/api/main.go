package main

import (
	"contralPlane/internal/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("load config %w", err)
	}

	fmt.Println(cfg)
}
