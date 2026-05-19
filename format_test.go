package askuserquestion

import "testing"

func TestFormat(t *testing.T) {
	t.Parallel()

	req := Request{
		Questions: []Question{
			{Question: "Auth method?", Header: "Auth", MultiSelect: false},
			{Question: "Languages?", Header: "Langs", MultiSelect: true},
			{Question: "Name?", Header: "Name", MultiSelect: false},
		},
	}

	t.Run("single multi other", func(t *testing.T) {
		t.Parallel()
		got := Format(req, []Answer{
			{Selected: []string{"OAuth"}},
			{Selected: []string{"Go", "Rust"}},
			{Other: "Vincent Adultman"},
		})
		want := "Auth method?\nOAuth\n\nLanguages?\n- Go\n- Rust\n\nName?\nVincent Adultman"
		if got != want {
			t.Fatalf("Format mismatch\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("multiselect single answer still uses dashes", func(t *testing.T) {
		t.Parallel()
		got := Format(Request{Questions: []Question{
			{Question: "Pick one or more", Header: "Pick", MultiSelect: true},
		}}, []Answer{{Selected: []string{"only"}}})
		want := "Pick one or more\n- only"
		if got != want {
			t.Fatalf("Format mismatch\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("missing answers are skipped", func(t *testing.T) {
		t.Parallel()
		got := Format(req, []Answer{{Selected: []string{"OAuth"}}})
		want := "Auth method?\nOAuth"
		if got != want {
			t.Fatalf("Format mismatch\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("other overrides selected", func(t *testing.T) {
		t.Parallel()
		got := Format(Request{Questions: []Question{
			{Question: "Pick", Header: "Pick"},
		}}, []Answer{{Selected: []string{"ignored"}, Other: "actually this"}})
		want := "Pick\nactually this"
		if got != want {
			t.Fatalf("Format mismatch\n got: %q\nwant: %q", got, want)
		}
	})

	t.Run("empty answer keeps question bare", func(t *testing.T) {
		t.Parallel()
		got := Format(Request{Questions: []Question{{Question: "Q?", Header: "H"}}}, []Answer{{}})
		want := "Q?"
		if got != want {
			t.Fatalf("Format mismatch\n got: %q\nwant: %q", got, want)
		}
	})
}
