package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/elastic/go-elasticsearch/v8/esutil"
	_ "github.com/joho/godotenv" // 可選：用於載入 .env 檔案
)

const (
	esAddress  = "http://192.168.0.236:9200"
	esUser     = "elastic"
	esPassword = "123456"
)

var es *elasticsearch.Client

func main() {
	// 可選：從 .env 檔案載入配置
	// godotenv.Load()

	cfg := elasticsearch.Config{
		Addresses: []string{esAddress},
		Username:  esUser,
		Password:  esPassword,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var err error
	es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	logger := log.New(os.Stdout, "[MyApp] ", log.LstdFlags)

	// 模擬產生一些日誌訊息
	logger.Println("Application started")
	logInfo("User logged in", "user_id", "123")
	logWarning("Low disk space", "partition", "/")
	logError("Database connection error", "error", "timeout")
	logger.Println("Application finished")
}

func logInfo(message string, args ...string) {
	sendLog("INFO", message, args...)
}

func logWarning(message string, args ...string) {
	sendLog("WARNING", message, args...)
}

func logError(message string, args ...string) {
	sendLog("ERROR", message, args...)
}

func sendLog(level, message string, args ...string) {
	if es == nil {
		log.Println("Elasticsearch client not initialized.")
		return
	}

	logData := map[string]interface{}{
		"@timestamp": time.Now().Format(time.RFC3339),
		"level":      level,
		"message":    message,
	}

	// 添加額外的上下文信息 (key-value pairs)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			logData[args[i]] = args[i+1]
		}
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		return
	}

	indexName := "app-" + time.Now().Format("2006.01.02") // 指定 index

	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    bytes.NewReader(jsonData),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Printf("Error sending log to Elasticsearch: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Printf("Error parsing Elasticsearch response: %s", err)
		} else {
			log.Printf("Elasticsearch error: [%d] %s", res.StatusCode, e["error"].(map[string]interface{})["type"])
		}
	}
}
