# Claude Analysis

A telemetry tool that collects and analyzes your Claude Code development activity, providing insights into your coding patterns and productivity.

## What does this tool do?

Claude Analysis automatically:
1. **Tracks your coding activity** - Monitors file reads, writes, and code changes
2. **Analyzes development patterns** - Counts lines written, characters processed, and tool usage
3. **Sends analytics data** - Uploads aggregated statistics to telemetry servers for insights
4. **Generates usage reports** - Returns structured data about your development session

## How it works

The tool operates in two modes:

### STOP Mode (Default)
- Reads a Python dictionary from stdin containing a `transcript_path`
- Loads and parses the JSONL conversation transcript file
- Aggregates all development activities from the session
- Sends analytics to the telemetry server

### POST_TOOL Mode
- Reads JSON lines directly from stdin
- Processes tool usage events in real-time
- Aggregates statistics and sends to server immediately

## Usage

### Basic Usage
```bash
# STOP mode (default) - reads transcript path from stdin
echo "{'transcript_path': '/path/to/conversation.jsonl'}" | ./claude_analysis

# POST_TOOL mode - reads JSON lines directly
MODE=POST_TOOL ./claude_analysis < tool_events.jsonl

# Custom API endpoint
./claude_analysis --o11y_base_url https://custom-server.com/api/upload < input.json
```

### Command Line Options
- `--o11y_base_url`: Override the default API endpoint URL (default: `https://gaia.mediatek.inc/o11y/upload_locs`)

### Environment Variables
- `MODE`: Set to `POST_TOOL` for direct JSON processing, or leave unset for STOP mode
- Can also create a `.env` file in the working directory with `MODE=POST_TOOL`

### Input Format

**STOP Mode Input:**
```
{'transcript_path': '/absolute/path/to/conversation.jsonl'}
```

**POST_TOOL Mode Input (JSONL):**
```json
{"type":"assistant","uuid":"msg1","cwd":"/workspace","sessionId":"session1","timestamp":"2025-01-01T00:00:00Z","message":{"content":[{"type":"tool_use","name":"Read"}]}}
{"parentUuid":"msg1","timestamp":"2025-01-01T00:00:01Z","toolUseResult":{"filePath":"file.txt","content":"Hello World"}}
```

## What gets tracked?

The tool analyzes and reports:

### File Operations
- **Read operations**: Files opened and content read
- **Write operations**: Files created or modified
- **Diff operations**: Code patches and changes applied

### Statistics Generated
- Total unique files accessed
- Total lines written
- Total characters read/written/modified
- Tool usage counts (Read, Write, ApplyDiff, etc.)
- Session metadata (workspace path, git repository, timestamps)

### Output Format

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

## Configuration

The tool uses these default settings:
- **API Endpoint**: `https://gaia.mediatek.inc/o11y/upload_locs` (can be overridden with `--o11y_base_url`)
- **Timeout**: 10 seconds
- **Extension Name**: "Claude-Code"
- **Insights Version**: "0.0.1"

Most configuration is automatically loaded from your system (username, machine ID). The API endpoint can be customized using the `--o11y_base_url` command line option.

## Integration

This tool is typically used as a hook in Claude Code:
1. Claude Code generates conversation transcripts
2. The transcript path is passed to claude_analysis
3. Analytics are processed and sent to the telemetry server
4. Results are returned for further processing

## Troubleshooting

**Problem**: Tool fails to read transcript file
**Solution**: Ensure the transcript path in your input is absolute and the file exists

**Problem**: Network timeout errors
**Solution**: Check your internet connection and firewall settings for the telemetry endpoint

**Problem**: JSON parsing errors
**Solution**: Verify your input format matches the expected structure for your chosen mode

**Problem**: Empty output
**Solution**: Check that your transcript file contains valid conversation data with tool usage events
