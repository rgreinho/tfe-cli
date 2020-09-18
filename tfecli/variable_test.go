package tfecli

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl"
)

func TestEncodeInt(t *testing.T) {
	testcases := []struct {
		hclVar string
		want   string
	}{
		{`stringvar = "expectedValue"`, `stringvar="expectedValue"`},
		{`listvar = ["item1", "item2"]`, `listvar=["item1", "item2"]`},
		{`mapvar = {
        key1 = "value1"
        key2 = {
          key21 = "value21"
        }
      }
      `, `mapvar={key1="value1",key2={key21="value21",},}`},
	}

	for _, tc := range testcases {
		var decoded interface{}
		_ = hcl.Decode(&decoded, tc.hclVar)
		decodedMap := reflect.ValueOf(decoded)
		iter := decodedMap.MapRange()
		for iter.Next() {
			got := EncodeVariable(iter.Key(), iter.Value())
			if got != tc.want {
				t.Errorf("Incorrect parsing got: %s, want: %s.", got, tc.want)
			}

		}
	}

}
