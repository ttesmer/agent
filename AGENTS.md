# Agent Development Guide

This is an AI agent written in Go that interfaces with OpenRouter (LLM API). It provides a REPL interface for chatting with the agent, which can execute bash commands as tools.

## Project Structure

```
.
├── main.go          # Single-file application (Agent, Client, Tool handling)
├── go.mod           # Go module (if present)
└── AGENTS.md        # This file
```

## Quick Start

### Prerequisites
- Go 1.21+ installed
- `OPENROUTER_API_KEY` environment variable set

### Running the Agent
```bash
go run main.go
```

### Environment Variables
- `OPENROUTER_API_KEY` - Required. Your OpenRouter API key
- `MODEL` - Optional. Defaults to `moonshotai/kimi-k2.5`
- `DEBUG=1` - Optional. Prints request/response JSON for debugging

## Architecture

The codebase is organized in `main.go`:

| Component | Description |
|-----------|-------------|
| `Agent` | Main loop, conversation management, tool orchestration |
| `Client` | HTTP client for OpenRouter API |
| `Message` / `Tool` / `ToolCall` | OpenAI-compatible message types |
| `executeTool()` | Tool implementations (run_command, etc.) |
| `handleToolCall()` | User approval flow for tool execution |

### Adding a New Tool

1. **Add tool definition** in `runInference()`:
   ```go
   {
       Type: "function",
       Function: FunctionDef{
           Name: "my_tool",
           Description: "What it does",
           Parameters: map[string]any{...},
       },
   }
   ```

2. **Add handler** in `executeTool()`:
   ```go
   case "my_tool":
       // Unmarshal params, execute logic, return result
   ```

3. **Test it**: Run the agent and ask it to use your tool.

## Development Workflow

### Branching
- Create feature branches: `git checkout -b feature/my-feature`
- Make focused commits with clear messages
- Merge to `main` when done, then delete feature branch

### Testing Changes
1. Run `go run main.go`
2. Test tool calls by asking the agent to execute commands
3. Set `DEBUG=1` to see full API request/response

### Code Style
- Keep the single-file structure for now (it's simple enough)
- Print user-facing output with color codes (see existing `\u001b[38;5;XXXm` patterns)
- Always return meaningful error messages (include output, not just exit codes)

## Common Tasks

### Check why a tool failed
```bash
DEBUG=1 go run main.go
# Then ask the agent to run the failing command
```

### Reset conversation state
Just quit (Ctrl+C) and restart - conversation is in-memory only.

### Modify tool behavior
Edit `executeTool()` function. The `run_command` tool uses `cmd.CombinedOutput()` so both stdout and stderr are captured.

## File Locations

The current working directory for bash commands is wherever you ran `go run main.go` from.

## Gotchas

- Tool calls require user approval (press Enter to approve, `n` to refuse)
- The agent streams input via stdin - one line = one user message
- JSON unmarshaling in `executeTool` uses struct tags - match the LLM's output keys exactly
