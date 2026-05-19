// Command askuserquestion-demo wires the askuserquestion tool to stdin/stdout
// so you can see the full request/answer loop without a real model. It
// simulates a model call with a hard-coded Params payload and prints the
// canonical tool result that would be sent back to the assistant.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"charm.land/fantasy"

	"github.com/dreamware-nz/askuserquestion-go"
)

func main() {
	tool := askuserquestion.NewTool(askuserquestion.StdinResolver{
		In:  os.Stdin,
		Out: os.Stdout,
	})

	params := askuserquestion.Params{
		Questions: []askuserquestion.Question{
			{
				Question: "Which auth method should we wire up first?",
				Header:   "Auth",
				Options: []askuserquestion.Option{
					{Label: "OAuth (Recommended)", Description: "Browser-based, no secrets in env"},
					{Label: "API key", Description: "Long-lived static token"},
					{Label: "mTLS", Description: "Client certificate, infra-side trust"},
				},
			},
			{
				Question:    "Which languages need first-class SDKs?",
				Header:      "SDKs",
				MultiSelect: true,
				Options: []askuserquestion.Option{
					{Label: "Go", Description: "for backend services"},
					{Label: "TypeScript", Description: "for web + edge"},
					{Label: "Python", Description: "for ML notebooks"},
					{Label: "Rust", Description: "for systems work"},
				},
			},
		},
	}

	input, _ := json.Marshal(params)
	resp, err := tool.Run(context.Background(), fantasy.ToolCall{
		ID:    "demo_call",
		Name:  tool.Info().Name,
		Input: string(input),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "tool error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n--- tool_result -------------------------")
	if resp.IsError {
		fmt.Fprintln(os.Stderr, "is_error=true")
	}
	fmt.Println(resp.Content)
}
