# Claude Analysis（Claude 分析工具）

一个遥测工具，用于收集和分析您的 Claude Code 开发活动，提供编程模式和生产力的深入分析。

## 这个工具能做什么？

Claude Analysis 自动：
1. **跟踪您的编程活动** - 监控文件读取、写入和代码更改
2. **分析开发模式** - 统计编写的行数、处理的字符数和工具使用情况
3. **发送分析数据** - 将聚合统计数据上传到遥测服务器以获得洞察
4. **生成使用报告** - 返回关于您开发会话的结构化数据

## 工作原理

该工具有两种操作模式：

### STOP 模式（默认）
- 从标准输入读取包含 `transcript_path` 的 Python 字典
- 加载并解析 JSONL 对话记录文件
- 聚合会话中的所有开发活动
- 将分析数据发送到遥测服务器

### POST_TOOL 模式
- 直接从标准输入读取 JSON 行
- 实时处理工具使用事件
- 聚合统计数据并立即发送到服务器

## 使用方法

### 基本用法
```bash
# STOP 模式（默认）- 从标准输入读取记录路径
echo "{'transcript_path': '/path/to/conversation.jsonl'}" | ./claude_analysis

# POST_TOOL 模式 - 直接读取 JSON 行
MODE=POST_TOOL ./claude_analysis < tool_events.jsonl

# 自定义 API 端点
./claude_analysis --o11y_base_url https://custom-server.com/api/upload < input.json
```

### 命令行选项
- `--o11y_base_url`: 覆盖默认的 API 端点 URL（默认值：`https://gaia.mediatek.inc/o11y/upload_locs`）

### 环境变量
- `MODE`: 设置为 `POST_TOOL` 进行直接 JSON 处理，或保持未设置使用 STOP 模式
- 也可以在工作目录中创建包含 `MODE=POST_TOOL` 的 `.env` 文件

### 输入格式

**STOP 模式输入：**
```
{'transcript_path': '/absolute/path/to/conversation.jsonl'}
```

**POST_TOOL 模式输入（JSONL）：**
```json
{"type":"assistant","uuid":"msg1","cwd":"/workspace","sessionId":"session1","timestamp":"2025-01-01T00:00:00Z","message":{"content":[{"type":"tool_use","name":"Read"}]}}
{"parentUuid":"msg1","timestamp":"2025-01-01T00:00:01Z","toolUseResult":{"filePath":"file.txt","content":"Hello World"}}
```

## 跟踪什么内容？

该工具分析和报告：

### 文件操作
- **读取操作**：打开和读取的文件内容
- **写入操作**：创建或修改的文件
- **差异操作**：应用的代码补丁和更改

### 生成的统计数据
- 访问的唯一文件总数
- 编写的总行数
- 读取/写入/修改的总字符数
- 工具使用次数（Read、Write、ApplyDiff 等）
- 会话元数据（工作区路径、git 仓库、时间戳）

### 输出格式

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

## 配置

工具使用这些默认设置：
- **API 端点**：`https://gaia.mediatek.inc/o11y/upload_locs`（可通过 `--o11y_base_url` 覆盖）
- **超时时间**：10 秒
- **扩展名称**："Claude-Code"
- **洞察版本**："0.0.1"

大部分配置会自动从您的系统加载（用户名、机器 ID）。API 端点可以通过 `--o11y_base_url` 命令行选项进行自定义。

## 集成

此工具通常用作 Claude Code 中的钩子：
1. Claude Code 生成对话记录
2. 记录路径传递给 claude_analysis
3. 分析数据被处理并发送到遥测服务器
4. 返回结果供进一步处理

## 故障排除

**问题**：工具无法读取记录文件
**解决方案**：确保输入中的记录路径是绝对路径且文件存在

**问题**：网络超时错误
**解决方案**：检查您的网络连接和遥测端点的防火墙设置

**问题**：JSON 解析错误
**解决方案**：验证您的输入格式是否符合所选模式的预期结构

**问题**：空输出
**解决方案**：检查您的记录文件是否包含带有工具使用事件的有效对话数据
