package bindata

import "testing"

func TestSanitize(t *testing.T) {
	sanitizeTests := []struct {
		in  string
		out string
	}{
		{`hello`, "`hello`"},
		{"hello\nworld", "`hello\nworld`"},
		{"`ello", "(\"`\" + `ello`)"},
		{"`a`e`i`o`u`", "(((\"`\" + `a`) + (\"`\" + (`e` + \"`\"))) + ((`i` + (\"`\" + `o`)) + (\"`\" + (`u` + \"`\"))))"},
		{"\xEF\xBB\xBF`s away!", "(\"\\xEF\\xBB\\xBF\" + (\"`\" + `s away!`))"},
	}

	for _, tt := range sanitizeTests {
		out := sanitize([]byte(tt.in))
		if string(out) != tt.out {
			t.Errorf("sanitize(%q):\nhave %q\nwant %q", tt.in, out, tt.out)
		}
	}
}
