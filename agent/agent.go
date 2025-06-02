package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// OutputEvent represents different types of output from the agent
type OutputEvent struct {
	// Type indicates the kind of output event
	Type string // "assistant_prefix", "user_prefix", "content", "newline", etc.
	// Content contains the event data
	Content string
	// Raw contains the raw event for specialized handlers
	Raw anthropic.MessageStreamEventUnion
}

// Constants for event types
const (
	EventAssistantPrefix = "assistant_prefix"
	EventUserPrefix      = "user_prefix"
	EventContent         = "content"
	EventNewline         = "newline"
	EventRaw             = "raw"
)

// ToolDefinition defines a tool that can be used by the AI
type ToolDefinition struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error)
}

// GenerateSchema creates a JSON schema for the given type
func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T

	schema := reflector.Reflect(v)

	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// ReadFileInput defines the input for the read_file tool
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory"`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

// ReadFileDefinition defines the read_file tool
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFile,
}

// ListFilesInput defines the input for the list_files tool
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

// ListFilesDefinition defines the list_files tool
var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: ListFilesInputSchema,
	Function:    ListFiles,
}

// EditFileInput defines the input for the edit_file tool
type EditFileInput struct {
	Path   string `json:"path" jsonschema_definition:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

var EditFileInputSchema = GenerateSchema[EditFileInput]()

// EditFileDefinition defines the edit_file tool
var EditFileDefinition = ToolDefinition{
	Name: "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_stir' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
	InputSchema: EditFileInputSchema,
	Function:    EditFile,
}

func NewAgent(client anthropic.Client, getUserMessage func() (string, bool), systemPrompt string) *Agent {
	// Create a buffered channel for output events
	outputChan := make(chan OutputEvent, 100)

	// Define available tools
	tools := []ToolDefinition{ReadFileDefinition, ListFilesDefinition, EditFileDefinition}

	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		SystemPrompt:   systemPrompt,
		outputChan:     outputChan,
		done:           make(chan struct{}),
		tools:          tools,
	}
}

type Agent struct {
	client         anthropic.Client
	getUserMessage func() (string, bool)
	SystemPrompt   string `json:"system,omitzero"`
	outputChan     chan OutputEvent
	done           chan struct{}
	tools          []ToolDefinition
}

// OutputChannel returns the channel that emits output events
func (a *Agent) OutputChannel() <-chan OutputEvent {
	return a.outputChan
}

// Close shuts down the agent and closes the output channel
func (a *Agent) Close() {
	close(a.done)
	close(a.outputChan)
}

// executeTool executes the specified tool with the given input
func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef ToolDefinition
	var found bool
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}

	a.outputChan <- OutputEvent{Type: EventContent, Content: fmt.Sprintf("Using tool: %s\n", name)}
	response, err := toolDef.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

func (a *Agent) Run(ctx context.Context) error {
	var messages []anthropic.MessageParam
	startConversation := "Start the conversation" // Create user message
	userContent := anthropic.ContentBlockParamUnion{
		OfText: &anthropic.TextBlockParam{Text: startConversation},
	}
	userMessage := anthropic.MessageParam{
		Role:    "user",
		Content: []anthropic.ContentBlockParamUnion{userContent},
	}
	messages = append(messages, userMessage)

	for {
		// We're letting Claude speak first, with streaming for responsiveness
		a.outputChan <- OutputEvent{Type: EventAssistantPrefix, Content: "Claude: "}
		message, err := a.runInference(ctx, messages)
		if err != nil {
			return err
		}
		a.outputChan <- OutputEvent{Type: EventNewline}

		// Convert the message to MessageParam for appending to the messages slice
		var messageContent []anthropic.ContentBlockParamUnion
		var toolResults []anthropic.ContentBlockParamUnion

		for _, block := range message.Content {
			if block.Type == "text" {
				textBlock := anthropic.TextBlockParam{Text: block.Text}
				messageContent = append(messageContent, anthropic.ContentBlockParamUnion{OfText: &textBlock})
			} else if block.Type == "tool_use" {
				// Process tool calls
				result := a.executeTool(block.ID, block.Name, block.Input)
				toolResults = append(toolResults, result)
			}
		}

		assistantMessage := anthropic.MessageParam{
			Role:    "assistant",
			Content: messageContent,
		}

		messages = append(messages, assistantMessage)

		// If we have tool results, add them to the conversation
		if len(toolResults) > 0 {
			toolResultMessage := anthropic.MessageParam{
				Role:    "user",
				Content: toolResults,
			}
			messages = append(messages, toolResultMessage)
			// Skip getting user input and go straight to next Claude response
			continue
		}

		// Now it's the user's turn
		a.outputChan <- OutputEvent{Type: EventUserPrefix, Content: "You: "}
		userInput, ok := a.getUserMessage()
		if !ok {
			break
		}

		userContent := anthropic.ContentBlockParamUnion{
			OfText: &anthropic.TextBlockParam{Text: userInput},
		}
		userMessage := anthropic.MessageParam{
			Role:    "user",
			Content: []anthropic.ContentBlockParamUnion{userContent},
		}
		messages = append(messages, userMessage)
	}

	return nil
}
func (a *Agent) runInference(ctx context.Context, messages []anthropic.MessageParam) (*anthropic.Message, error) {
	// Create the message creation request
	sysPrompt := []anthropic.TextBlockParam{
		{Text: a.SystemPrompt},
	}

	// Convert tool definitions to Anthropic tool parameters
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	req := anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 1024,
		System:    sysPrompt,
		Messages:  messages,
		Tools:     anthropicTools,
	}

	// Use streaming API for better responsiveness
	stream := a.client.Messages.NewStreaming(ctx, req)
	if stream.Err() != nil {
		return nil, stream.Err()
	}
	// Don't defer stream.Close() here - we'll handle it explicitly in each case

	// Initialize an empty message to accumulate content
	finalMessage := &anthropic.Message{}

	// Create channels for communication with the streaming goroutine
	done := make(chan struct{})
	errCh := make(chan error, 1)

	// Process streaming in a separate goroutine to allow for clean cancellation
	go func() {
		// Ensure the stream is closed when the goroutine exits
		defer func() {
			stream.Close()
			close(done)
		}()

		for stream.Next() {
			// Get the current event
			event := stream.Current()

			// Send the raw event to the output channel
			a.outputChan <- OutputEvent{Type: EventRaw, Raw: event}

			// Extract and send text content for convenience
			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					a.outputChan <- OutputEvent{Type: EventContent, Content: deltaVariant.Text}
				}
			}

			// Accumulate the event into the message
			err := finalMessage.Accumulate(event)
			if err != nil {
				errCh <- fmt.Errorf("error accumulating event: %w", err)
				return
			}

			// Check if context was cancelled between events
			if ctx.Err() != nil {
				errCh <- ctx.Err()
				return
			}
		}

		// Check for any errors in the stream
		if stream.Err() != nil {
			errCh <- fmt.Errorf("stream error: %w", stream.Err())
			return
		}
	}()

	// Wait for either completion, error, or cancellation
	select {
	case <-done:
		// Stream completed successfully
		return finalMessage, nil
	case err := <-errCh:
		// Error occurred during streaming
		return nil, err
	case <-ctx.Done():
		// Context was cancelled - explicitly close the stream to force any blocking operations to stop
		stream.Close()
		return nil, fmt.Errorf("operation cancelled: %w", ctx.Err())
	case <-a.done:
		// Agent is being shut down
		stream.Close()
		return nil, fmt.Errorf("agent shutting down")
	}
}

// Tool implementations

func ReadFile(input json.RawMessage) (string, error) {
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	err := json.Unmarshal(input, &listFilesInput)
	if err != nil {
		return "", err
	}

	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func EditFile(input json.RawMessage) (string, error) {
	editFileInput := EditFileInput{}
	err := json.Unmarshal(input, &editFileInput)
	if err != nil {
		return "", err
	}

	if editFileInput.Path == "" || editFileInput.OldStr == editFileInput.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(editFileInput.Path)
	if err != nil {
		if os.IsNotExist(err) && editFileInput.OldStr == "" {
			return createNewFile(editFileInput.Path, editFileInput.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, editFileInput.OldStr, editFileInput.NewStr, -1)

	if oldContent == newContent && editFileInput.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	err = os.WriteFile(editFileInput.Path, []byte(newContent), 0644)
	if err != nil {
		return "", err
	}

	return "OK", nil
}

func createNewFile(filePath, content string) (string, error) {
	dir := filepath.Dir(filePath)
	if dir != "." {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
