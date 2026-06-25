package bots

import "testing"

// TestLikeEscaper verifies that LIKE/ILIKE wildcard characters supplied by the
// user are escaped so they are matched literally instead of being interpreted
// as pattern metacharacters. The repository wraps the escaped term in `%...%`
// itself, so the escaper must only neutralise `\`, `%` and `_`.
func TestLikeEscaper(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "plain text is unchanged", input: "hello", want: "hello"},
		{name: "percent is escaped", input: "50%", want: `50\%`},
		{name: "underscore is escaped", input: "a_b", want: `a\_b`},
		{name: "backslash is escaped first", input: `a\b`, want: `a\\b`},
		{name: "mixed metacharacters", input: `%_\`, want: `\%\_\\`},
		{name: "empty string", input: "", want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := likeEscaper.Replace(tc.input)
			if got != tc.want {
				t.Fatalf("likeEscaper.Replace(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
