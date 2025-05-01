package main

import (
	"custom-database/cmd/console_mode"
	"custom-database/cmd/http_mode"
	"custom-database/config"
	"custom-database/internal/backend"
	"custom-database/internal/http/handlers"
	"custom-database/internal/parser"
	"flag"
	"log"
)

// @title Custom Database API
// @version 1.0
// @description API для работы с кастомной базой данных
// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	mode := flag.String("mode", "console", "Режим работы: console или http")
	port := flag.String("port", cfg.Port, "Порт для HTTP сервера")
	flag.Parse()

	parser := parser.NewParser()
	mb, err := backend.NewMemoryBackend(cfg)
	if err != nil {
		log.Fatal("Error creating memory backend:", err)
	}

	handlers := handlers.NewHttpHandlers(parser, mb)

	switch *mode {
	case "console":
		console_mode.RunConsoleMode(parser, mb)
	case "http":
		http_mode.RunHttpServer(handlers, *port)
	default:
		log.Fatal("Неизвестный режим работы. Используйте 'console' или 'http'")
	}
}
