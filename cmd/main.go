package main

import (
	"custom-database/cmd/console_mode"
	"custom-database/cmd/http_mode"
	"custom-database/config"
	"custom-database/internal/executor"
	"custom-database/internal/http/handlers"
	"custom-database/internal/lexer"
	"custom-database/internal/storage"
	"flag"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	mode := flag.String("mode", "console", "Режим работы: console или http")
	port := flag.String("port", cfg.Port, "Порт для HTTP сервера")
	flag.Parse()

	storage := storage.NewStorage(cfg)
	executor := executor.NewExecutor(storage)
	lexer := lexer.NewLexer(executor)

	handlers := handlers.NewHttpHandlers(lexer)

	switch *mode {
	case "console":
		console_mode.RunConsoleMode(lexer)
	case "http":
		http_mode.RunHttpServer(handlers, *port)
	default:
		log.Fatal("Неизвестный режим работы. Используйте 'console' или 'http'")
	}
}
