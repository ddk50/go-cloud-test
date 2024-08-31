package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type response struct {
	OK bool `json:"ok"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// JSON レスポンスを生成
	w.Header().Set("Content-Type", "application/json")
	resp := response{OK: true}

	// エンコード中にエラーが発生した場合はログに出力
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
		return
	}

	// リクエストの詳細をログに記録
	log.Printf("Handled request: %s %s from %s in %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
}

func main() {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		if _, err := strconv.Atoi(envPort); err != nil {
			port = envPort
		} else {
			log.Printf("Invalid PORT environment variable value: %s, using default port %s", envPort, port)
		}
	}

	http.HandleFunc("/", handler)

	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
