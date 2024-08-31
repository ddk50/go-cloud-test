package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	// HTTPS サーバーの設定
	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		Handler:  http.HandlerFunc(handler),
		ErrorLog: log.New(log.Writer(), "HTTPS Server: ", log.LstdFlags|log.Lshortfile),
	}

	fmt.Println("Starting server on https://localhost:443")

	// サーバー開始とエラーログ出力
	if err := server.ListenAndServeTLS(
		"/app/ssl/server.crt",
		"/app/ssl/server.key"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
