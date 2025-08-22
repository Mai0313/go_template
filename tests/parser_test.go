package tests

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"claude_analysis/core/config"
	"claude_analysis/core/telemetry"
)

func TestParser_FromTestConversationJSONL_PrintsFullPayload(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to get caller info")
	}
	jsonlPath := filepath.Join(filepath.Dir(thisFile), "test_conversation.jsonl")

	records, err := telemetry.ReadJSONL(jsonlPath)
	if err != nil {
		t.Fatalf("ReadJSONL error: %v", err)
	}

	statsArr := telemetry.AggregateConversationStats(records)
	if len(statsArr) != 1 {
		t.Fatalf("expected 1 stats record, got %d", len(statsArr))
	}

	cfg := config.Default()
	payload := map[string]interface{}{
		"user":            cfg.UserName,
		"records":         statsArr,
		"extensionName":   cfg.ExtensionName,
		"machineId":       cfg.MachineID,
		"insightsVersion": cfg.InsightsVersion,
	}
	pretty, err := json.MarshalIndent(payload, "", "  ")
	if err == nil {
		t.Logf("Full transformed payload:\n%s", string(pretty))
	} else {
		t.Logf("Full transformed payload (marshal error: %v)", err)
	}
}

func TestParser_ComprehensiveSyntheticEvents(t *testing.T) {
	t.Helper()

	tmpDir := t.TempDir()
	// prepare .git/config to validate gitRemoteUrl detection
	gitCfgDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitCfgDir, 0o755); err != nil {
		t.Fatalf("mkdir .git dir: %v", err)
	}
	remoteURL := "git@github.com:org/repo.git"
	cfgContent := "[remote \"origin\"]\n    url = " + remoteURL + "\n"
	if err := os.WriteFile(filepath.Join(gitCfgDir, "config"), []byte(cfgContent), 0o644); err != nil {
		t.Fatalf("write .git/config: %v", err)
	}

	// Build synthetic records covering read/write/apply_diff and various field placements
	recs := []map[string]interface{}{
		{
			"type":      "assistant",
			"uuid":      "a1",
			"cwd":       tmpDir,
			"sessionId": "sess123",
			"timestamp": "2025-01-01T00:00:00Z",
			"message": map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"type": "tool_use", "name": "Read"},
				},
			},
		},
		{
			"parentUuid": "a1",
			"timestamp":  "2025-01-01T00:00:01Z",
			"toolUseResult": map[string]interface{}{
				// top-level fields (no nested file object)
				"filePath": "fileA.txt",
				"content":  "hello世界\n", // 8 runes (5+2+1)
			},
		},
		{
			"type":      "assistant",
			"uuid":      "b1",
			"cwd":       tmpDir,
			"sessionId": "sess123",
			"timestamp": "2025-01-01T00:00:02Z",
			"message": map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"type": "tool_use", "name": "Write"},
				},
			},
		},
		{
			"parentUuid": "b1",
			"timestamp":  "2025-01-01T00:00:03Z",
			"toolUseResult": map[string]interface{}{
				// nested file object
				"file": map[string]interface{}{
					"filePath": "fileB.txt",
					"content":  "line1\nline2\n", // 12 runes, 3 lines
				},
			},
		},
		{
			"type":      "assistant",
			"uuid":      "c1",
			"cwd":       tmpDir,
			"sessionId": "sess123",
			"timestamp": "2025-01-01T00:00:04Z",
			"message": map[string]interface{}{
				"content": []interface{}{
					map[string]interface{}{"type": "tool_use", "name": "ApplyDiff"},
				},
			},
		},
		{
			"parentUuid": "c1",
			"timestamp":  "2025-01-01T00:00:05Z",
			"toolUseResult": map[string]interface{}{
				// structuredPatch fallback & top-level filePath
				"structuredPatch": []string{"@@ -1,1 +1,1 @@"},
				"filePath":        "fileB.txt",
			},
		},
	}

	statsArr := telemetry.AggregateConversationStats(recs)
	if len(statsArr) != 1 {
		t.Fatalf("expected 1 stats record, got %d", len(statsArr))
	}
	stats := statsArr[0]

	// Totals
	if stats.TotalUniqueFiles != 2 {
		t.Errorf("TotalUniqueFiles expected 2, got %d", stats.TotalUniqueFiles)
	}
	if stats.TotalReadCharacters != 8 {
		t.Errorf("TotalReadCharacters expected 8, got %d", stats.TotalReadCharacters)
	}
	if stats.TotalWriteLines != 3 {
		t.Errorf("TotalWriteLines expected 3, got %d", stats.TotalWriteLines)
	}
	if stats.TotalWriteCharacters != 12 {
		t.Errorf("TotalWriteCharacters expected 12, got %d", stats.TotalWriteCharacters)
	}
	if stats.TotalDiffCharacters <= 0 {
		t.Errorf("TotalDiffCharacters expected > 0, got %d", stats.TotalDiffCharacters)
	}

	// Details
	if len(stats.ReadFileDetails) != 1 || stats.ReadFileDetails[0].FilePath != "fileA.txt" {
		t.Errorf("ReadFileDetails mismatch: %+v", stats.ReadFileDetails)
	}
	if len(stats.WriteToFileDetails) != 1 || stats.WriteToFileDetails[0].FilePath != "fileB.txt" {
		t.Errorf("WriteToFileDetails mismatch: %+v", stats.WriteToFileDetails)
	}
	if stats.WriteToFileDetails[0].LineCount != 3 {
		t.Errorf("Write LineCount expected 3, got %d", stats.WriteToFileDetails[0].LineCount)
	}
	if len(stats.ApplyDiffDetails) != 1 || stats.ApplyDiffDetails[0].FilePath != "fileB.txt" {
		t.Errorf("ApplyDiffDetails mismatch: %+v", stats.ApplyDiffDetails)
	}

	// Tool call counts
	if stats.ToolCallCounts["Read"] != 1 || stats.ToolCallCounts["Write"] != 1 || stats.ToolCallCounts["ApplyDiff"] != 1 {
		t.Errorf("ToolCallCounts mismatch: %+v", stats.ToolCallCounts)
	}

	// Context fields
	if stats.TaskID != "sess123" {
		t.Errorf("taskId expected 'sess123', got %s", stats.TaskID)
	}
	if stats.FolderPath != tmpDir {
		t.Errorf("folderPath expected %s, got %s", tmpDir, stats.FolderPath)
	}
	// last timestamp should be from 00:00:05Z
	wantLast, _ := time.Parse(time.RFC3339, "2025-01-01T00:00:05Z")
	if stats.Timestamp != wantLast.UnixNano()/int64(time.Millisecond) {
		t.Errorf("timestamp mismatch, got %d", stats.Timestamp)
	}
	if stats.GitRemoteURL != remoteURL {
		t.Errorf("gitRemoteUrl expected %s, got %s", remoteURL, stats.GitRemoteURL)
	}
}

func TestParser_EmptyRecords_ReturnsEmpty(t *testing.T) {
	statsArr := telemetry.AggregateConversationStats(nil)
	if len(statsArr) != 0 {
		t.Fatalf("expected empty stats for empty input, got %d", len(statsArr))
	}
}

// Integration tests that execute the binary and hit network are purposely omitted
// to keep tests hermetic. End-to-end behavior is covered by unit tests using
// telemetry.AggregateConversationStats and real sample JSONL lines.

func TestParser_POST_TOOL_FromSubsetJSONLines(t *testing.T) {
	// Read a subset of lines (assistant tool_use Read + its toolUseResult) and aggregate
	_, thisFile, _, _ := runtime.Caller(0)
	jsonlPath := filepath.Join(filepath.Dir(thisFile), "test_conversation.jsonl")

	f, err := os.Open(jsonlPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var assistant map[string]interface{}
	var toolres map[string]interface{}
	var assistantUUID string
	for scanner.Scan() {
		line := scanner.Text()
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		if assistant == nil {
			if obj["type"] == "assistant" {
				// check contains tool_use name Read
				if msg, ok := obj["message"].(map[string]interface{}); ok {
					if arr, ok := msg["content"].([]interface{}); ok {
						for _, it := range arr {
							m, _ := it.(map[string]interface{})
							if m != nil && m["type"] == "tool_use" && m["name"] == "Read" {
								assistant = obj
								if u, _ := obj["uuid"].(string); u != "" {
									assistantUUID = u
								}
							}
						}
					}
				}
			}
		} else if toolres == nil {
			if p, _ := obj["parentUuid"].(string); p != "" && p == assistantUUID {
				if _, ok := obj["toolUseResult"].(map[string]interface{}); ok {
					toolres = obj
					break
				}
			}
		}
	}
	if assistant == nil || toolres == nil {
		t.Fatalf("failed to locate assistant/toolUseResult lines in %s", jsonlPath)
	}
	subset := []map[string]interface{}{assistant, toolres}
	statsArr := telemetry.AggregateConversationStats(subset)
	if len(statsArr) != 1 {
		t.Fatalf("expected 1, got %d", len(statsArr))
	}
	stats := statsArr[0]
	if len(stats.ReadFileDetails) < 1 {
		t.Errorf("expected at least 1 read detail, got %d", len(stats.ReadFileDetails))
	}
}
