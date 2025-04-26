package main

import (
	"custom-database/internal/lexer"
	"fmt"
	"log"
	"net/http"
)

func runHttpServer(lexer lexer.Lexer, port string) {
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.FormValue("query")
		if query == "" {
			http.Error(w, "Query is required", http.StatusBadRequest)
			return
		}

		err := lexer.ParseQuery(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Query executed successfully")
	})

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("HTTP сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Error starting HTTP server:", err)
	}
}
