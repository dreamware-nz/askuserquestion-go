package askuserquestion

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"charm.land/fantasy"
)

//go:embed description.md
var description string

// Description returns the canonical tool description shipped with the SDK.
// Exposed so hosts that build their own ToolInfo can reuse the exact prose
// Claude was trained against.
func Description() string { return description }

// NewTool builds a fantasy.AgentTool using the default tool name (ToolName)
// and the given Resolver. The returned tool validates inputs against the
// Claude Agent SDK schema, hands valid requests to the Resolver, and formats
// the user's answers using Format.
//
// On validation failure the tool returns a non-error ToolResponse with
// IsError=true so the model can self-correct on the next turn. On resolver
// cancellation (ErrResolverCancelled) it returns a non-error response with
// a stable "[cancelled]" body so the conversation stays valid.
func NewTool(r Resolver) fantasy.AgentTool {
	return NewToolNamed(ToolName, r)
}

// NewToolNamed is like NewTool but lets the host override the tool name.
// Pass SDKToolName ("AskUserQuestion") for strict Claude Agent SDK parity.
func NewToolNamed(name string, r Resolver) fantasy.AgentTool {
	if r == nil {
		r = ResolverFunc(func(_ context.Context, _ Request) ([]Answer, error) {
			return nil, errors.New("askuserquestion: no resolver configured")
		})
	}
	return fantasy.NewAgentTool(
		name,
		description,
		func(ctx context.Context, p Params, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
			if err := Validate(p); err != nil {
				return fantasy.NewTextErrorResponse(err.Error()), nil
			}
			req := Request{
				ToolCallID: call.ID,
				Questions:  p.Questions,
			}
			answers, err := r.Ask(ctx, req)
			if err != nil {
				if errors.Is(err, ErrResolverCancelled) {
					return fantasy.NewTextResponse("[cancelled by user]"), nil
				}
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return fantasy.ToolResponse{}, err
				}
				return fantasy.NewTextErrorResponse(fmt.Sprintf("resolver error: %v", err)), nil
			}
			return fantasy.NewTextResponse(Format(req, answers)), nil
		},
	)
}
