package main

import (
	"custom-database/cmd/mode"
	"custom-database/config"
	"custom-database/internal/executor"
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

	modeArg := flag.String("mode", "console", "Режим работы: console или http")
	portArg := flag.String("port", cfg.Port, "Порт для HTTP сервера")
	flag.Parse()

	storage := storage.NewStorage(cfg)
	executor := executor.NewExecutor(storage)
	lexer := lexer.NewLexer(executor)

	switch *modeArg {
	case "console":
		mode.RunConsoleMode(lexer)
	case "http":
		mode.RunHttpServer(lexer, *portArg)
	default:
		log.Fatal("Неизвестный режим работы. Используйте 'console' или 'http'")
	}
}
