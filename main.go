package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	projectID  = "pfj-test-434203"
	bucketName = "pfj-test-434203-company1"
)

type response struct {
	OK bool `json:"ok"`
}

func bqtest() {
	ctx := context.Background()

	start := time.Now()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// TODO エラー処理
	defer client.Close()

	// データセットとテーブルのIDを設定
	// これは事前に用意しとかないといけないの？
	datasetID := "pfjtest"
	tableID := "ioexecute"

	client.Jobs(ctx)

	// GCSから読み込むCSVファイルのURIを設定
	gcsURI := "gs://pfj-test-434203-company1/ioexecute.csv"

	// GCSのCSVファイルからBigQueryテーブルにデータをロードするジョブを設定
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.SkipLeadingRows = 1 // CSVのヘッダ行をスキップ

	gcsRef.Schema = bigquery.Schema{
		{Name: "execute_id", Type: bigquery.IntegerFieldType},
		{Name: "executed_at", Type: bigquery.TimestampFieldType},
		{Name: "type", Type: bigquery.StringFieldType},
		{Name: "misc", Type: bigquery.StringFieldType},
	}

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteEmpty

	// ジョブを実行してデータをロード
	job, err := loader.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to start job: %v", err)
	}

	// ジョブの完了を待機
	status, err := job.Wait(ctx)
	if err != nil {
		log.Fatalf("Job failed: %v", err)
	}

	// ジョブのステータスを確認
	if status.Err() != nil {
		log.Fatalf("Job completed with error: %v", status.Err())
	}

	log.Printf("Data loaded successfully in %v\n", time.Since(start))
}

func bqTest2() {
	ctx := context.Background()

	start := time.Now()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// TODO エラー処理
	defer client.Close()

	// データセットとテーブルのIDを設定
	// これは事前に用意しとかないといけないの？
	datasetID := "sample"
	tableID := "sales"

	client.Jobs(ctx)

	// gs://pfj-test-434203-company1/dummy_data_test_name.csv
	// gs://pfj-test-434203-company1/dummy_data_timestamp_milliseconds.csv
	// GCSから読み込むCSVファイルのURIを設定
	gcsURI := "gs://pfj-test-434203-company1/dummy_data_test_name.csv"

	// GCSのCSVファイルからBigQueryテーブルにデータをロードするジョブを設定
	gcsRef := bigquery.NewGCSReference(gcsURI)
	gcsRef.SkipLeadingRows = 1 // CSVのヘッダ行をスキップ

	gcsRef.Schema = bigquery.Schema{
		{Name: "timestamp", Type: bigquery.TimestampFieldType},
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "product", Type: bigquery.StringFieldType},
		{Name: "unitprice", Type: bigquery.NumericFieldType},
		{Name: "quantity", Type: bigquery.NumericFieldType},
	}

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(gcsRef)
	loader.WriteDisposition = bigquery.WriteTruncate

	// ジョブを実行してデータをロード
	job, err := loader.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to start job: %v", err)
	}

	// ジョブの完了を待機
	status, err := job.Wait(ctx)
	if err != nil {
		log.Fatalf("Job failed: %v", err)
	}

	// ジョブのステータスを確認
	if status.Err() != nil {
		log.Fatalf("Job completed with error: %v", status.Err())
	}

	log.Printf("Data loaded successfully in %v\n", time.Since(start))
}

func handlerUploadFileToGCS(w http.ResponseWriter, r *http.Request) {
	// ファイルのアップロード処理
	ctx := context.Background()

	start := time.Now()

	// Google Cloud Storage クライアントの作成
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("keys/gs-accessor.json"))
	if err != nil {
		http.Error(w, "Failed to create GCS client", http.StatusInternalServerError)
		log.Printf("Failed to create GCS client: %v", err)
		return
	}
	defer client.Close()

	// POSTされたファイルを取得
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from request", http.StatusBadRequest)
		log.Printf("Failed to retrieve file: %v", err)
		return
	}
	defer file.Close()

	// GCS バケット内のオブジェクト名を設定（ここではアップロード時のファイル名を使用）
	objectName := header.Filename

	// GCS バケットにファイルをアップロード
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		http.Error(w, "Failed to upload file to GCS", http.StatusInternalServerError)
		log.Printf("Failed to upload file to GCS: %v", err)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Failed to finalize file upload", http.StatusInternalServerError)
		log.Printf("Failed to finalize file upload: %v", err)
		return
	}

	log.Printf("File uploaded to GCS: %s in %v", objectName, time.Since(start))
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
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

func handlerBigQuery(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// JSON レスポンスを生成
	w.Header().Set("Content-Type", "application/json")
	resp := response{OK: true}

	bqtest()

	// エンコード中にエラーが発生した場合はログに出力
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
		return
	}

	log.Printf("Handled request: %s %s from %s in %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
}

func handlerBigQuery2(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// JSON レスポンスを生成
	w.Header().Set("Content-Type", "application/json")
	resp := response{OK: true}

	bqTest2()

	// エンコード中にエラーが発生した場合はログに出力
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v", err)
		return
	}

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

	http.HandleFunc("/", handlerHello)
	http.HandleFunc("/upload", handlerUploadFileToGCS)
	http.HandleFunc("/bigquery", handlerBigQuery)
	http.HandleFunc("/bigquery2", handlerBigQuery2)

	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
