# askuserquestion-go

Portable Go implementation of Claude Code's `AskUserQuestion` tool, packaged
as a reusable [`charm.land/fantasy`](https://pkg.go.dev/charm.land/fantasy)
`AgentTool`.

The tool itself is host-rendered: the model emits a structured tool call
with a list of multiple-choice questions, the host renders a picker, the
user selects, and the canonical answer string flows back as the tool
result. This module owns the **schema, validation, and canonical answer
formatting**. Hosts supply a `Resolver` that actually talks to the user
(TUI modal, web form, stdin prompt, mock, whatever).

## Why split this out

The `AskUserQuestion` semantics — 1-4 questions, 2-4 options each, an
implicit "Other" free-text option, a specific tool-result string format —
are fixed by the Claude Agent SDK. They have nothing to do with which agent
host you're running. Putting them in their own module means:

- one canonical implementation, reusable across Crush, custom CLIs, web
  agents, headless workers
- versionable schema (`v0.1.0`, etc.) without coupling to a host's
  internal package layout
- testable in isolation — table-driven validation, golden answer formats,
  no UI dependencies

## Install

```sh
go get github.com/dreamware-nz/askuserquestion-go
```

Requires Go 1.26+ (because of `charm.land/fantasy`).

## Quick start

```go
import (
    "github.com/dreamware-nz/askuserquestion-go"
)

// 1. Implement a Resolver — the seam that actually asks the user.
type myResolver struct{ /* TUI handle, channels, whatever */ }
func (r *myResolver) Ask(ctx context.Context, req askuserquestion.Request) ([]askuserquestion.Answer, error) {
    // push req onto your UI, block on the user's selection, return answers
    // in the same order as req.Questions.
}

// 2. Register the tool with your fantasy-based agent.
tool := askuserquestion.NewTool(&myResolver{})
agent.RegisterTool(tool)
```

There's a working CLI demo in [`examples/cli/`](./examples/cli) that wires
the tool to stdin/stdout — useful for sanity-checking the schema and the
canonical answer format without booting a real model.

```sh
go run ./examples/cli
```

## Tool name

By default the tool registers as `ask_user_question` (snake case, to match
common Go agent-host conventions). To match the official Claude Agent SDK
spelling exactly:

```go
tool := askuserquestion.NewToolNamed(askuserquestion.SDKToolName, resolver)
// → "AskUserQuestion"
```

## Schema

Matches the Claude Agent SDK wire format exactly:

```json
{
  "questions": [
    {
      "question": "Which auth method?",
      "header": "Auth",
      "multiSelect": false,
      "options": [
        { "label": "OAuth (Recommended)", "description": "Browser flow" },
        { "label": "API key",             "description": "Static token" }
      ]
    }
  ]
}
```

Hard limits enforced by `Validate`:

| Field           | Rule                                  |
| --------------- | ------------------------------------- |
| `questions`     | 1-4 entries                           |
| `question`      | non-empty                             |
| `header`        | non-empty (12-char cap is advisory)   |
| `options`       | 2-4 entries                           |
| `option.label`  | non-empty, unique within the question |

The 12-char header cap is **advisory only**. Validation does not reject
long headers so a misbehaving model can't poison the conversation;
renderers are expected to truncate.

## Canonical answer format

`Format(req, answers)` produces the string Anthropic expects back as the
tool result:

```
Auth method?
OAuth

Languages?
- Go
- Rust

Name?
Vincent Adultman
```

Rules:

- questions joined by a blank line, in original order
- single-select: bare label
- multi-select: each label prefixed with `- `
- `Answer.Other` (free text from the implicit "Other" option) overrides
  `Answer.Selected` and is emitted verbatim

## Resolvers

The module ships three:

- **`StaticResolver`** — canned answers, for tests and dry runs.
- **`ResolverFunc`** — adapter for ad-hoc functions.
- **`StdinResolver`** — numbered prompts on stdout, picks on stdin.
  Production hosts should write a richer UI; this exists for examples and
  smoke tests.

For real agent hosts (Crush, web servers, etc.) you'll usually want a
`Resolver` that:

1. Publishes the `Request` onto a broker / channel that the UI is
   listening on.
2. Blocks (respecting `ctx`) until the UI sends back the user's selection.
3. Returns `ErrResolverCancelled` if the user dismisses or supersedes the
   request (e.g. sends a new chat message instead of answering).

The tool surfaces `ErrResolverCancelled` as a non-error tool result with
body `[cancelled by user]` so the conversation stays valid.

## License

MIT. The tool description prose, schema, and canonical answer format
follow Anthropic's public Claude Code / Claude Agent SDK design. This
module is an independent reimplementation in Go.
