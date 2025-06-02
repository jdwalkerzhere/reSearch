package agent

import (
	"fmt"
	"io"
	"os"
)

// ConsoleClient implements a console-based output handler for the agent
type ConsoleClient struct {
	writer io.Writer
}

// NewConsoleClient creates a new console client with the specified writer
func NewConsoleClient(writer io.Writer) *ConsoleClient {
	if writer == nil {
		writer = os.Stdout
	}
	return &ConsoleClient{writer: writer}
}

// Run starts processing output events from the agent
func (c *ConsoleClient) Run(agent *Agent) {
	for event := range agent.OutputChannel() {
		switch event.Type {
		case EventAssistantPrefix:
			fmt.Fprint(c.writer, "\u001b[93m"+event.Content+"\u001b[0m")
		case EventUserPrefix:
			fmt.Fprint(c.writer, "\u001b[94m"+event.Content+"\u001b[0m")
		case EventContent:
			fmt.Fprint(c.writer, event.Content)
		case EventNewline:
			fmt.Fprintln(c.writer)
		}
	}
}
