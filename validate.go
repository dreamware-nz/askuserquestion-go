package askuserquestion

import "strings"

// MaxQuestions is the upper bound on questions per call (Claude Agent SDK).
const MaxQuestions = 4

// MinOptions / MaxOptions bound the options-per-question range.
const (
	MinOptions = 2
	MaxOptions = 4
)

// RecommendedHeaderLen is the soft cap on Header length from the SDK prose.
// It is advisory: Validate does not reject headers longer than this so a
// misbehaving model does not break the conversation, but renderers are
// expected to truncate to this width.
const RecommendedHeaderLen = 12

// Validate enforces the structural rules from the Claude Agent SDK schema:
//
//   - 1 <= len(Questions) <= MaxQuestions
//   - each Question has non-empty Question and Header
//   - MinOptions <= len(Options) <= MaxOptions
//   - each Option has a non-empty Label
//   - Option labels are unique within a question
//
// Validate returns one of the sentinel errors from errors.go so callers can
// branch on errors.Is. The header length cap is advisory and not checked.
func Validate(p Params) error {
	if len(p.Questions) == 0 {
		return ErrNoQuestions
	}
	if len(p.Questions) > MaxQuestions {
		return ErrTooManyQuestions
	}
	for _, q := range p.Questions {
		if strings.TrimSpace(q.Question) == "" {
			return ErrEmptyQuestion
		}
		if strings.TrimSpace(q.Header) == "" {
			return ErrEmptyHeader
		}
		if len(q.Options) < MinOptions || len(q.Options) > MaxOptions {
			return ErrOptionCount
		}
		seen := make(map[string]struct{}, len(q.Options))
		for _, o := range q.Options {
			if strings.TrimSpace(o.Label) == "" {
				return ErrEmptyLabel
			}
			if _, dup := seen[o.Label]; dup {
				return ErrDuplicateLabel
			}
			seen[o.Label] = struct{}{}
		}
	}
	return nil
}
