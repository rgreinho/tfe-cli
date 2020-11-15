package tfecli

import (
	"testing"
)

func TestParseToHCLVariables(t *testing.T) {
	testcases := []struct {
		hclVar string
		want   HCLVariable
	}{
		{`stringvar = expectedValue`, HCLVariable{key: "stringvar", value: "expectedValue"}},
	}
	for _, tc := range testcases {
		gots := SimpleParseVarFile(tc.hclVar)
		for _, got := range gots {
			if got != tc.want {
				t.Errorf("Incorrect parsing got: %s, want: %s.", got, tc.want)
			}
		}
	}
}
