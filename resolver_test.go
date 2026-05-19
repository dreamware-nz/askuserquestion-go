package askuserquestion

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestStaticResolver(t *testing.T) {
	t.Parallel()
	want := []Answer{{Selected: []string{"x"}}}
	r := StaticResolver{Answers: want}
	got, err := r.Ask(context.Background(), Request{})
	if err != nil {
		t.Fatalf("Ask returned err: %v", err)
	}
	if len(got) != 1 || got[0].Selected[0] != "x" {
		t.Fatalf("unexpected answers: %+v", got)
	}

	sentinel := errors.New("boom")
	r2 := StaticResolver{Err: sentinel}
	if _, err := r2.Ask(context.Background(), Request{}); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel, got %v", err)
	}
}

func TestResolverFunc(t *testing.T) {
	t.Parallel()
	called := false
	r := ResolverFunc(func(_ context.Context, _ Request) ([]Answer, error) {
		called = true
		return []Answer{{Selected: []string{"ok"}}}, nil
	})
	got, err := r.Ask(context.Background(), Request{})
	if err != nil || !called || len(got) != 1 {
		t.Fatalf("unexpected: called=%v got=%+v err=%v", called, got, err)
	}
}

func TestStdinResolverSingleSelect(t *testing.T) {
	t.Parallel()
	q := Question{
		Question: "Auth?",
		Header:   "Auth",
		Options: []Option{
			{Label: "OAuth"}, {Label: "API key"},
		},
	}
	in := bytes.NewBufferString("2\n")
	out := &bytes.Buffer{}
	r := StdinResolver{In: in, Out: out}

	ans, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if err != nil {
		t.Fatalf("Ask err: %v", err)
	}
	if len(ans) != 1 || len(ans[0].Selected) != 1 || ans[0].Selected[0] != "API key" {
		t.Fatalf("unexpected answer: %+v", ans)
	}
	if !strings.Contains(out.String(), "Auth?") {
		t.Fatalf("output missing question text:\n%s", out.String())
	}
}

func TestStdinResolverMultiSelect(t *testing.T) {
	t.Parallel()
	q := Question{
		Question:    "Langs?",
		Header:      "Langs",
		MultiSelect: true,
		Options:     []Option{{Label: "Go"}, {Label: "Rust"}, {Label: "Zig"}},
	}
	r := StdinResolver{In: bytes.NewBufferString("1,3\n"), Out: io.Discard}
	ans, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if err != nil {
		t.Fatalf("Ask err: %v", err)
	}
	want := []string{"Go", "Zig"}
	if len(ans[0].Selected) != 2 || ans[0].Selected[0] != want[0] || ans[0].Selected[1] != want[1] {
		t.Fatalf("unexpected selection: %+v", ans[0].Selected)
	}
}

func TestStdinResolverOther(t *testing.T) {
	t.Parallel()
	q := Question{
		Question: "Name?",
		Header:   "Name",
		Options:  []Option{{Label: "Alice"}, {Label: "Bob"}},
	}
	r := StdinResolver{In: bytes.NewBufferString("other\nVincent\n"), Out: io.Discard}
	ans, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if err != nil {
		t.Fatalf("Ask err: %v", err)
	}
	if ans[0].Other != "Vincent" || len(ans[0].Selected) != 0 {
		t.Fatalf("unexpected answer: %+v", ans[0])
	}
}

func TestStdinResolverEmptyCancels(t *testing.T) {
	t.Parallel()
	q := Question{Question: "Q?", Header: "Q", Options: []Option{{Label: "a"}, {Label: "b"}}}
	r := StdinResolver{In: bytes.NewBufferString("\n"), Out: io.Discard}
	_, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if !errors.Is(err, ErrResolverCancelled) {
		t.Fatalf("expected ErrResolverCancelled, got %v", err)
	}
}

func TestStdinResolverRejectsMultiOnSingleSelect(t *testing.T) {
	t.Parallel()
	q := Question{Question: "Q?", Header: "Q", Options: []Option{{Label: "a"}, {Label: "b"}}}
	r := StdinResolver{In: bytes.NewBufferString("1,2\n"), Out: io.Discard}
	_, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if err == nil || !strings.Contains(err.Error(), "single-select") {
		t.Fatalf("expected single-select error, got %v", err)
	}
}

func TestStdinResolverRejectsOutOfRange(t *testing.T) {
	t.Parallel()
	q := Question{Question: "Q?", Header: "Q", Options: []Option{{Label: "a"}, {Label: "b"}}}
	r := StdinResolver{In: bytes.NewBufferString("9\n"), Out: io.Discard}
	_, err := r.Ask(context.Background(), Request{Questions: []Question{q}})
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("expected out-of-range error, got %v", err)
	}
}
