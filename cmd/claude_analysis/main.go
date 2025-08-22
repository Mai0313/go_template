package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"claude_analysis/core/config"
	"claude_analysis/core/telemetry"
)

// readStdinAndSave reads JSON data from stdin, sends it to API and returns response
func readStdinAndSave(baseURL string) map[string]interface{} {
	// Load configuration
	cfg := config.Default()

	// Override API endpoint if baseURL is provided
	if baseURL != "" {
		cfg.API.Endpoint = baseURL
	}

	// Create telemetry client
	client := telemetry.New(cfg)

	// 讀取 stdin
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Printf("[ERROR] Failed to read from stdin: %v", err)
		return map[string]interface{}{"status": "error", "message": "failed to read stdin"}
	}

	var aggregated []telemetry.ApiConversationStats
	if strings.EqualFold(cfg.Mode, "POST_TOOL") {
		// 支援直接吃任意一行 JSONL（或整段文字含多行），逐行解析、彙整
		reader := bufio.NewScanner(strings.NewReader(string(stdinData)))
		raw := make([]map[string]interface{}, 0)
		for reader.Scan() {
			line := strings.TrimSpace(reader.Text())
			if line == "" {
				continue
			}
			var obj map[string]interface{}
			if err := json.Unmarshal([]byte(line), &obj); err == nil {
				raw = append(raw, obj)
			} else {
				log.Printf("[WARN] JSON parsing failed, skipping line: %v", err)
			}
		}
		if len(raw) == 0 {
			log.Printf("[ERROR] No valid JSON lines found in POST_TOOL mode")
			return map[string]interface{}{"status": "error", "message": "no valid JSON lines found"}
		}
		aggregated = telemetry.AggregateConversationStats(raw)
	} else { // STOP (default)
		path, err := telemetry.ExtractTranscriptPath(string(stdinData))
		if err != nil {
			log.Printf("[ERROR] Failed to extract transcript path: %v", err)
			return map[string]interface{}{"status": "error", "message": "failed to extract transcript path"}
		}
		log.Printf("[INFO] Extracted transcript path: %s", path)
		data, err := telemetry.ReadJSONL(path)
		if err != nil {
			log.Printf("[ERROR] Failed to read JSONL file: %v", err)
			return map[string]interface{}{"status": "error", "message": "failed to read JSONL file"}
		}
		aggregated = telemetry.AggregateConversationStats(data)
	}

	// 透過解析器聚合統計，包裝成單一物件 {user, records, ...}
	payload := map[string]interface{}{
		"user":            cfg.UserName,
		"records":         aggregated,
		"extensionName":   cfg.ExtensionName,
		"machineId":       cfg.MachineID,
		"insightsVersion": cfg.InsightsVersion,
	}

	// 送出
	response, err := client.Submit(payload)
	if err != nil {
		log.Printf("[ERROR] API call failed (endpoint: %s): %v", cfg.API.Endpoint, err)
		return map[string]interface{}{"status": "error", "message": "API call failed", "endpoint": cfg.API.Endpoint}
	}

	log.Printf("[INFO] Successfully sent telemetry data to %s", cfg.API.Endpoint)
	return response
}

func main() {
	// 配置 logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[Claude-Analysis] ")

	// 获取环境变量中的 base URL，如果没有则使用默认值
	defaultBaseURL := "https://gaia.mediatek.inc/o11y/upload_locs"
	if envURL := os.Getenv("O11Y_BASE_URL"); envURL != "" {
		defaultBaseURL = envURL
		log.Printf("[INFO] Read API endpoint from environment variable O11Y_BASE_URL: %s", envURL)
	}

	// Parse command line flags (命令行参数优先级最高)
	var o11yBaseURL = flag.String("o11y_base_url", defaultBaseURL, "Base URL for o11y API endpoint")
	flag.Parse()

	// 确定最终使用的 URL
	finalURL := *o11yBaseURL
	if finalURL != defaultBaseURL && os.Getenv("O11Y_BASE_URL") != "" {
		log.Printf("[INFO] Command line argument --o11y_base_url overrides environment variable, using: %s", finalURL)
	}

	log.Printf("[INFO] claude_analysis starting...")
	inputData := readStdinAndSave(finalURL)

	log.Printf("[INFO] readStdinAndSave completed, preparing output...")
	if len(inputData) > 0 {
		jsonOutput, err := json.MarshalIndent(inputData, "", "  ")
		if err != nil {
			log.Printf("[ERROR] JSON marshaling failed: %v", err)
			// 即使序列化失敗，也輸出簡單的錯誤訊息而不是中斷程序
			fmt.Println(`{"status": "error", "message": "JSON marshaling failed"}`)
		} else {
			fmt.Println(string(jsonOutput))
		}
	} else {
		log.Printf("[WARN] No data to output")
		fmt.Println(`{"status": "no_data", "message": "no data to output"}`)
	}
	log.Printf("[INFO] claude_analysis completed")
}
