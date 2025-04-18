package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://192.168.0.236:9200"},
		Username:  "elastic",
		Password:  "123456",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 定義要索引的文檔
	logData := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     "INFO",
		"message":   "This is a test log message from Golang",
		"service":   "golang-app",
	}

	// 將文檔編碼為 JSON
	jsonData, err := json.Marshal(logData)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %s", err)
	}

	// 定義索引名稱
	indexName := "golang-logs-" + time.Now().Format("2006.01.02")

	// 創建 IndexRequest
	req := esapi.IndexRequest{
		Index:   indexName,
		Body:    bytes.NewReader(jsonData),
		Refresh: "true", // 使索引操作立即可見
	}

	// 執行請求
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response from Elasticsearch: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			log.Fatalf("Elasticsearch error: [%d] %s", res.StatusCode, e["error"].(map[string]interface{})["type"])
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	fmt.Printf("Indexed document ID: %v\n", r["_id"])
	fmt.Println("Log message sent to Elasticsearch!")
}
