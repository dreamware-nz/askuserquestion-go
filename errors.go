package askuserquestion

import "errors"

// Validation errors. Wrapped by Validate so callers can check with errors.Is.
var (
	// ErrNoQuestions means the model called the tool with an empty questions
	// array. The SDK requires at least one question.
	ErrNoQuestions = errors.New("askuserquestion: questions must be a non-empty array")
	// ErrTooManyQuestions means more than 4 questions were supplied.
	ErrTooManyQuestions = errors.New("askuserquestion: maximum 4 questions allowed")
	// ErrEmptyQuestion means a Question.Question field was empty.
	ErrEmptyQuestion = errors.New("askuserquestion: each question must have non-empty question text")
	// ErrEmptyHeader means a Question.Header field was empty.
	ErrEmptyHeader = errors.New("askuserquestion: each question must have a non-empty header")
	// ErrOptionCount means options length was outside [2, 4].
	ErrOptionCount = errors.New("askuserquestion: each question must have 2-4 options")
	// ErrEmptyLabel means an Option.Label field was empty.
	ErrEmptyLabel = errors.New("askuserquestion: each option must have a non-empty label")
	// ErrDuplicateLabel means a single question listed the same label twice.
	ErrDuplicateLabel = errors.New("askuserquestion: option labels within a question must be unique")
)

// ErrResolverCancelled is returned by Resolvers when the user dismisses or
// supersedes the request (for example by sending a fresh chat message before
// answering). The tool surfaces this as a non-error tool result so the
// conversation stays valid.
var ErrResolverCancelled = errors.New("askuserquestion: resolver cancelled by user")
