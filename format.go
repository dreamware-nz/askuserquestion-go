package askuserquestion

import "strings"

// Format renders the user's answers into the canonical tool-result string
// used by the Claude Agent SDK. The format is:
//
//	<question 1>
//	- <label>
//	- <label>
//
//	<question 2>
//	<label>
//
// Rules:
//   - questions are joined by a blank line in the original Question order
//   - if Answer.Other is set, its text is emitted verbatim as the body
//   - if Question.MultiSelect is true, each selected label is prefixed "- "
//   - if Question.MultiSelect is false, the single label is emitted bare
//   - answers are paired by index; missing answers are skipped
//
// Format is deterministic and side-effect-free; it is the only canonical
// answer formatter in this package.
func Format(req Request, answers []Answer) string {
	parts := make([]string, 0, len(req.Questions))
	for i, q := range req.Questions {
		if i >= len(answers) {
			break
		}
		body := answerBody(q, answers[i])
		if body == "" {
			parts = append(parts, q.Question)
			continue
		}
		parts = append(parts, q.Question+"\n"+body)
	}
	return strings.Join(parts, "\n\n")
}

func answerBody(q Question, a Answer) string {
	if strings.TrimSpace(a.Other) != "" {
		return a.Other
	}
	if len(a.Selected) == 0 {
		return ""
	}
	if !q.MultiSelect && len(a.Selected) == 1 {
		return a.Selected[0]
	}
	lines := make([]string, len(a.Selected))
	for i, s := range a.Selected {
		lines[i] = "- " + s
	}
	return strings.Join(lines, "\n")
}
