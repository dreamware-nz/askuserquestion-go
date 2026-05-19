package askuserquestion

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	okOpts := []Option{
		{Label: "Yes", Description: "do it"},
		{Label: "No", Description: "don't"},
	}
	okQ := Question{Question: "Proceed?", Header: "Proceed", Options: okOpts}

	tests := []struct {
		name string
		in   Params
		want error
	}{
		{
			name: "valid single question",
			in:   Params{Questions: []Question{okQ}},
			want: nil,
		},
		{
			name: "valid four questions",
			in:   Params{Questions: []Question{okQ, okQ, okQ, okQ}},
			want: nil,
		},
		{
			name: "empty questions",
			in:   Params{Questions: nil},
			want: ErrNoQuestions,
		},
		{
			name: "five questions",
			in:   Params{Questions: []Question{okQ, okQ, okQ, okQ, okQ}},
			want: ErrTooManyQuestions,
		},
		{
			name: "blank question text",
			in: Params{Questions: []Question{
				{Question: "   ", Header: "h", Options: okOpts},
			}},
			want: ErrEmptyQuestion,
		},
		{
			name: "blank header",
			in: Params{Questions: []Question{
				{Question: "q", Header: "", Options: okOpts},
			}},
			want: ErrEmptyHeader,
		},
		{
			name: "one option",
			in: Params{Questions: []Question{
				{Question: "q", Header: "h", Options: okOpts[:1]},
			}},
			want: ErrOptionCount,
		},
		{
			name: "five options",
			in: Params{Questions: []Question{
				{Question: "q", Header: "h", Options: []Option{
					{Label: "a"}, {Label: "b"}, {Label: "c"}, {Label: "d"}, {Label: "e"},
				}},
			}},
			want: ErrOptionCount,
		},
		{
			name: "blank label",
			in: Params{Questions: []Question{
				{Question: "q", Header: "h", Options: []Option{{Label: ""}, {Label: "x"}}},
			}},
			want: ErrEmptyLabel,
		},
		{
			name: "duplicate label",
			in: Params{Questions: []Question{
				{Question: "q", Header: "h", Options: []Option{{Label: "a"}, {Label: "a"}}},
			}},
			want: ErrDuplicateLabel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Validate(tt.in)
			if !errors.Is(got, tt.want) {
				t.Fatalf("Validate(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValidateAllowsLongHeader(t *testing.T) {
	t.Parallel()
	// Header length is advisory in the SDK; Validate must not reject it so
	// a misbehaving model can not poison the conversation.
	q := Question{
		Question: "?",
		Header:   "this header is much longer than twelve",
		Options:  []Option{{Label: "a"}, {Label: "b"}},
	}
	if err := Validate(Params{Questions: []Question{q}}); err != nil {
		t.Fatalf("Validate rejected long header: %v", err)
	}
}
