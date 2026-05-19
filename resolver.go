package askuserquestion

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Resolver is the host-supplied seam that actually asks the user. The library
// validates the model's input, builds a Request, and hands it to Ask. The
// Resolver returns one Answer per Question (in the same order) or an error.
//
// Implementations should:
//   - respect ctx cancellation (e.g. the user closed the session)
//   - return ErrResolverCancelled when the user dismisses/supersedes
//   - never return more answers than there are questions
type Resolver interface {
	Ask(ctx context.Context, req Request) ([]Answer, error)
}

// ResolverFunc adapts a plain function into a Resolver.
type ResolverFunc func(ctx context.Context, req Request) ([]Answer, error)

// Ask implements Resolver.
func (f ResolverFunc) Ask(ctx context.Context, req Request) ([]Answer, error) {
	return f(ctx, req)
}

// StaticResolver returns the same canned answers for every call. Useful for
// tests, dry runs, and golden-file snapshots.
type StaticResolver struct {
	Answers []Answer
	Err     error
}

// Ask implements Resolver.
func (s StaticResolver) Ask(_ context.Context, _ Request) ([]Answer, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	out := make([]Answer, len(s.Answers))
	copy(out, s.Answers)
	return out, nil
}

// StdinResolver is a minimal CLI resolver: it prints each question to Out and
// reads selections from In. Options are numbered starting at 1; the user
// types the number, a comma-separated list of numbers for multi-select, or
// the literal word "other" followed by free text. Empty input cancels.
//
// This resolver exists primarily for examples and headless smoke tests.
// Production hosts should implement a richer UI (TUI modal, web form, etc.).
type StdinResolver struct {
	In  io.Reader
	Out io.Writer
}

// Ask implements Resolver.
func (s StdinResolver) Ask(ctx context.Context, req Request) ([]Answer, error) {
	if s.In == nil || s.Out == nil {
		return nil, fmt.Errorf("askuserquestion: StdinResolver requires both In and Out")
	}
	r := bufio.NewReader(s.In)
	answers := make([]Answer, 0, len(req.Questions))

	for i, q := range req.Questions {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		fmt.Fprintf(s.Out, "\n[%d/%d] %s\n%s\n", i+1, len(req.Questions), q.Header, q.Question)
		for j, opt := range q.Options {
			if opt.Description != "" {
				fmt.Fprintf(s.Out, "  %d) %s — %s\n", j+1, opt.Label, opt.Description)
			} else {
				fmt.Fprintf(s.Out, "  %d) %s\n", j+1, opt.Label)
			}
		}
		fmt.Fprintf(s.Out, "  o) Other (type free text)\n")
		if q.MultiSelect {
			fmt.Fprint(s.Out, "> select (comma-separated): ")
		} else {
			fmt.Fprint(s.Out, "> select: ")
		}

		line, err := readLine(r)
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return nil, ErrResolverCancelled
		}

		ans, err := parseStdinAnswer(q, line, r, s.Out)
		if err != nil {
			return nil, err
		}
		ans.Question = q.Question
		answers = append(answers, ans)
	}

	return answers, nil
}

func parseStdinAnswer(q Question, line string, r *bufio.Reader, out io.Writer) (Answer, error) {
	if strings.EqualFold(line, "o") || strings.EqualFold(line, "other") {
		fmt.Fprint(out, "  other: ")
		other, err := readLine(r)
		if err != nil {
			return Answer{}, err
		}
		other = strings.TrimSpace(other)
		if other == "" {
			return Answer{}, ErrResolverCancelled
		}
		return Answer{Other: other}, nil
	}

	picks := strings.Split(line, ",")
	if !q.MultiSelect && len(picks) > 1 {
		return Answer{}, fmt.Errorf("askuserquestion: question is single-select, got %d picks", len(picks))
	}

	selected := make([]string, 0, len(picks))
	for _, p := range picks {
		p = strings.TrimSpace(p)
		n, err := strconv.Atoi(p)
		if err != nil {
			return Answer{}, fmt.Errorf("askuserquestion: invalid selection %q", p)
		}
		if n < 1 || n > len(q.Options) {
			return Answer{}, fmt.Errorf("askuserquestion: selection %d out of range", n)
		}
		selected = append(selected, q.Options[n-1].Label)
	}
	return Answer{Selected: selected}, nil
}

func readLine(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil && (err != io.EOF || line == "") {
		return "", err
	}
	return line, nil
}
