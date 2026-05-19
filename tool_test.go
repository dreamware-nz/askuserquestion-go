package askuserquestion

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"charm.land/fantasy"
)

func TestNewToolInfo(t *testing.T) {
	t.Parallel()
	tool := NewTool(StaticResolver{})
	info := tool.Info()
	if info.Name != ToolName {
		t.Fatalf("name = %q want %q", info.Name, ToolName)
	}
	if !strings.Contains(info.Description, "ask the user questions") {
		t.Fatalf("description missing canonical prose:\n%s", info.Description)
	}
	if _, ok := info.Parameters["questions"]; !ok {
		t.Fatalf("schema missing 'questions' property: %#v", info.Parameters)
	}
	if len(info.Required) == 0 || info.Required[0] != "questions" {
		t.Fatalf("schema required mismatch: %v", info.Required)
	}
}

func TestNewToolNamed(t *testing.T) {
	t.Parallel()
	tool := NewToolNamed(SDKToolName, StaticResolver{})
	if tool.Info().Name != SDKToolName {
		t.Fatalf("expected SDK name, got %q", tool.Info().Name)
	}
}

func TestToolRunValidation(t *testing.T) {
	t.Parallel()
	tool := NewTool(StaticResolver{})
	resp, err := tool.Run(context.Background(), fantasy.ToolCall{
		ID:    "call_1",
		Name:  ToolName,
		Input: `{"questions": []}`,
	})
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if !resp.IsError {
		t.Fatalf("expected validation error response, got %+v", resp)
	}
	if !strings.Contains(resp.Content, "non-empty array") {
		t.Fatalf("error content: %q", resp.Content)
	}
}

func TestToolRunHappyPath(t *testing.T) {
	t.Parallel()
	tool := NewTool(StaticResolver{Answers: []Answer{{Selected: []string{"OAuth"}}}})

	input := mustJSON(t, Params{Questions: []Question{{
		Question: "Auth?",
		Header:   "Auth",
		Options:  []Option{{Label: "OAuth"}, {Label: "API key"}},
	}}})

	resp, err := tool.Run(context.Background(), fantasy.ToolCall{ID: "x", Name: ToolName, Input: input})
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.IsError {
		t.Fatalf("unexpected error response: %+v", resp)
	}
	want := "Auth?\nOAuth"
	if resp.Content != want {
		t.Fatalf("content = %q want %q", resp.Content, want)
	}
}

func TestToolRunResolverCancel(t *testing.T) {
	t.Parallel()
	tool := NewTool(StaticResolver{Err: ErrResolverCancelled})

	input := mustJSON(t, Params{Questions: []Question{{
		Question: "?",
		Header:   "?",
		Options:  []Option{{Label: "a"}, {Label: "b"}},
	}}})

	resp, err := tool.Run(context.Background(), fantasy.ToolCall{ID: "x", Name: ToolName, Input: input})
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.IsError {
		t.Fatalf("cancel should not be flagged as error: %+v", resp)
	}
	if !strings.Contains(resp.Content, "cancelled") {
		t.Fatalf("content: %q", resp.Content)
	}
}

func TestToolRunResolverGenericError(t *testing.T) {
	t.Parallel()
	tool := NewTool(StaticResolver{Err: errors.New("network down")})

	input := mustJSON(t, Params{Questions: []Question{{
		Question: "?",
		Header:   "?",
		Options:  []Option{{Label: "a"}, {Label: "b"}},
	}}})

	resp, err := tool.Run(context.Background(), fantasy.ToolCall{ID: "x", Name: ToolName, Input: input})
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if !resp.IsError {
		t.Fatalf("expected IsError on resolver failure: %+v", resp)
	}
}

func TestToolRunPropagatesContextCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tool := NewTool(ResolverFunc(func(ctx context.Context, _ Request) ([]Answer, error) {
		return nil, ctx.Err()
	}))

	input := mustJSON(t, Params{Questions: []Question{{
		Question: "?",
		Header:   "?",
		Options:  []Option{{Label: "a"}, {Label: "b"}},
	}}})

	_, err := tool.Run(ctx, fantasy.ToolCall{ID: "x", Name: ToolName, Input: input})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
