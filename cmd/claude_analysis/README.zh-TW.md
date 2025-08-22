# Claude Analysis（Claude 分析工具）

一個遙測工具，用於收集和分析您的 Claude Code 開發活動，提供程式撰寫模式和生產力的深入分析。

## 這個工具能做什麼？

Claude Analysis 自動：
1. **追蹤您的程式撰寫活動** - 監控檔案讀取、寫入和程式碼變更
2. **分析開發模式** - 統計撰寫的行數、處理的字元數和工具使用情況
3. **傳送分析資料** - 將彙整統計資料上傳到遙測伺服器以獲得洞察
4. **產生使用報告** - 回傳關於您開發會話的結構化資料

## 運作原理

該工具有兩種操作模式：

### STOP 模式（預設）
- 從標準輸入讀取包含 `transcript_path` 的 Python 字典
- 載入並解析 JSONL 對話記錄檔案
- 彙整會話中的所有開發活動
- 將分析資料傳送到遙測伺服器

### POST_TOOL 模式
- 直接從標準輸入讀取 JSON 行
- 即時處理工具使用事件
- 彙整統計資料並立即傳送到伺服器

## 使用方法

### 基本用法
```bash
# STOP 模式（預設）- 從標準輸入讀取記錄路徑
echo "{'transcript_path': '/path/to/conversation.jsonl'}" | ./claude_analysis

# POST_TOOL 模式 - 直接讀取 JSON 行
MODE=POST_TOOL ./claude_analysis < tool_events.jsonl

# 自訂 API 端點
./claude_analysis --o11y_base_url https://custom-server.com/api/upload < input.json
```

### 命令列選項
- `--o11y_base_url`: 覆蓋預設的 API 端點 URL（預設值：`https://gaia.mediatek.inc/o11y/upload_locs`）

### 環境變數
- `MODE`: 設定為 `POST_TOOL` 進行直接 JSON 處理，或保持未設定使用 STOP 模式
- 也可以在工作目錄中建立包含 `MODE=POST_TOOL` 的 `.env` 檔案

### 輸入格式

**STOP 模式輸入：**
```
{'transcript_path': '/absolute/path/to/conversation.jsonl'}
```

**POST_TOOL 模式輸入（JSONL）：**
```json
{"type":"assistant","uuid":"msg1","cwd":"/workspace","sessionId":"session1","timestamp":"2025-01-01T00:00:00Z","message":{"content":[{"type":"tool_use","name":"Read"}]}}
{"parentUuid":"msg1","timestamp":"2025-01-01T00:00:01Z","toolUseResult":{"filePath":"file.txt","content":"Hello World"}}
```

## 追蹤什麼內容？

該工具分析和報告：

### 檔案操作
- **讀取操作**：開啟和讀取的檔案內容
- **寫入操作**：建立或修改的檔案
- **差異操作**：套用的程式碼修補和變更

### 產生的統計資料
- 存取的唯一檔案總數
- 撰寫的總行數
- 讀取/寫入/修改的總字元數
- 工具使用次數（Read、Write、ApplyDiff 等）
- 會話中繼資料（工作區路徑、git 儲存庫、時間戳記）

### 輸出格式

- [Example Output](./examples/claude_code_log.json)

```json
{
  "user": "your-username",
  "records": [{
    "totalUniqueFiles": 5,
    "totalWriteLines": 120,
    "totalReadCharacters": 2500,
    "totalWriteCharacters": 1800,
    "totalDiffCharacters": 350,
    "toolCallCounts": {"Read": 8, "Write": 3, "ApplyDiff": 1},
    "taskId": "session-id",
    "timestamp": 1704067200000,
    "folderPath": "/path/to/workspace",
    "gitRemoteUrl": "https://github.com/user/repo.git"
  }],
  "extensionName": "Claude-Code",
  "machineId": "unique-machine-id",
  "insightsVersion": "0.0.1"
}
```

## 設定

工具使用這些預設設定：
- **API 端點**：`https://gaia.mediatek.inc/o11y/upload_locs`（可透過 `--o11y_base_url` 覆蓋）
- **逾時時間**：10 秒
- **擴充套件名稱**："Claude-Code"
- **洞察版本**："0.0.1"

大部分設定會自動從您的系統載入（使用者名稱、機器 ID）。API 端點可以透過 `--o11y_base_url` 命令列選項進行自訂。

## 整合

此工具通常用作 Claude Code 中的掛鉤：
1. Claude Code 產生對話記錄
2. 記錄路徑傳遞給 claude_analysis
3. 分析資料被處理並傳送到遙測伺服器
4. 回傳結果供進一步處理

## 疑難排解

**問題**：工具無法讀取記錄檔案
**解決方案**：確保輸入中的記錄路徑是絕對路徑且檔案存在

**問題**：網路逾時錯誤
**解決方案**：檢查您的網路連線和遙測端點的防火牆設定

**問題**：JSON 解析錯誤
**解決方案**：驗證您的輸入格式是否符合所選模式的預期結構

**問題**：空輸出
**解決方案**：檢查您的記錄檔案是否包含帶有工具使用事件的有效對話資料
