// Package askuserquestion implements Claude Code's AskUserQuestion tool as a
// portable Go library that plugs into any charm.land/fantasy-based agent host.
//
// The tool is host-rendered: the model emits a structured tool call, the host
// renders a picker UI, the user selects, and the canonical answer string is
// returned as the tool result. This package owns the schema, validation, and
// canonical formatting; hosts supply a Resolver that actually talks to the
// user (TUI modal, stdin prompt, web form, mock, whatever).
package askuserquestion

// ToolName is the canonical tool name exposed to the model.
//
// The original Claude Agent SDK calls it "AskUserQuestion". We expose snake
// case as the default to match Crush's naming conventions; hosts that want
// the SDK-compatible spelling can pass a custom name to NewToolNamed.
const ToolName = "ask_user_question"

// SDKToolName is the spelling used by Anthropic's Claude Agent SDK.
const SDKToolName = "AskUserQuestion"

// Params is the top-level input schema for the tool.
//
// JSON tags match the Claude Agent SDK wire format exactly so models trained
// against the SDK can drive this implementation without prompt adjustments.
type Params struct {
	// Questions is the ordered list of questions to ask. 1-4 entries.
	Questions []Question `json:"questions" description:"Questions to ask the user (1-4 questions)"`
	// Context is an optional short hint shown to the user above the questions
	// so they can see what work the questions belong to (e.g. project name,
	// repo path, current task). It is opaque metadata: the library does not
	// validate its content, and it never appears in the canonical answer
	// string returned to the model. Hosts that render a UI should display it
	// near the top of the picker; resolvers that don't render UI may ignore it.
	Context string `json:"context,omitempty" description:"Optional short hint shown above the questions to remind the user what work the questions belong to (e.g. project name, repo path, current task). Free text, ~200 chars or less."`
}

// Question is a single question with multiple-choice options.
type Question struct {
	// Question is the full prompt shown to the user.
	Question string `json:"question" description:"The complete question to ask the user"`
	// Header is a very short label (<= 12 chars) used as a column header or
	// compact title in the picker UI. Examples: "Auth method", "Library".
	Header string `json:"header" description:"Very short label (max 12 chars). Examples: \"Auth method\", \"Library\""`
	// Options are the available choices. 2-4 entries.
	Options []Option `json:"options" description:"The available choices (2-4 options)"`
	// MultiSelect allows the user to pick more than one option.
	MultiSelect bool `json:"multiSelect" description:"Set to true to allow multiple selections"`
}

// Option is one choice within a Question.
type Option struct {
	// Label is the short display text for this option.
	Label string `json:"label" description:"The display text for this option"`
	// Description explains what selecting this option means.
	Description string `json:"description" description:"Explanation of what this option means"`
}

// Answer is the user's response to a single Question.
//
// Exactly one of Selected or Other should be populated. Hosts that allow the
// user to pick "Other" set Other to the typed text and leave Selected empty.
type Answer struct {
	// Question is the original question text. Used to pair answers back to
	// questions when order is unstable.
	Question string
	// Selected is the list of option labels the user chose. For single-select
	// questions this has at most one element.
	Selected []string
	// Other is free-text the user typed when they picked the implicit "Other"
	// option. If set, Selected is ignored by Format.
	Other string
}

// Request is what the tool hands to a Resolver. It carries the model-supplied
// questions plus the upstream tool-call ID so hosts that route through a
// pubsub broker or stream-json child can correlate the reply.
type Request struct {
	// ToolCallID is the fantasy.ToolCall.ID for this invocation.
	ToolCallID string
	// Questions is the validated question set.
	Questions []Question
	// Context is the optional hint copied from Params.Context. Resolvers that
	// render a UI should surface it near the top of the picker; resolvers
	// without a UI may ignore it.
	Context string
}
