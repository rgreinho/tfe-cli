package tfecli

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/hashicorp/hcl"
)

// ParseVarFile reads an HCL varfile and returns its content as a map.
func ParseVarFile(varFile string) (reflect.Value, error) {
	fileContent, err := ioutil.ReadFile(varFile)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("cannot read the file %q: %s", varFile, err)
	}

	return ParseVarFileContent(fileContent)
}

// ParseVarFileContent parses the content of an HCL varfile and returns its content as a map.
func ParseVarFileContent(bs []byte) (reflect.Value, error) {
	var out interface{}
	err := hcl.Unmarshal(bs, &out)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("cannot read the HCL content: %s", err)
	}
	outMap := reflect.ValueOf(out)
	return outMap, nil
}

func encodeIntVariable(key reflect.Value, value reflect.Value) string {
	return fmt.Sprintf(`%s=%d`, key, value.Interface())
}

func encodeFloatVariable(key reflect.Value, value reflect.Value) string {
	return fmt.Sprintf(`%s=%f`, key, value.Interface())
}

func encodeMapVariable(key reflect.Value, value reflect.Value) string {
	iter := value.MapRange()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("{"))
	for iter.Next() {
		sb.WriteString(EncodeVariable(iter.Key(), iter.Value()))
		sb.WriteString(",")
	}
	sb.WriteString("}")
	return sb.String()
}

func encodeSliceVariable(key reflect.Value, value reflect.Value) string {
	b := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		b[i] = fmt.Sprintf("%s", value.Index(i))
	}
	return fmt.Sprintf("%s=[\"%s\"]", key, strings.Join(b, "\", \""))
}

func encodeStringVariable(key reflect.Value, value reflect.Value) string {
	return fmt.Sprintf(`%s="%s"`, key, value.Interface())
}

// EncodeVariable encodes a variable defined in a tfvars file.
func EncodeVariable(key reflect.Value, value reflect.Value) string {
	concreteValue := reflect.ValueOf(value.Interface())
	switch concreteValue.Kind() {
	case reflect.Int:
		return encodeIntVariable(key, concreteValue)
	case reflect.Float64:
		return encodeFloatVariable(key, concreteValue)
	case reflect.Map:
		return encodeMapVariable(key, concreteValue)
	case reflect.Slice:
		if concreteValue.Index(0).Kind() == reflect.Map {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%s=", key))
			sb.WriteString(EncodeVariable(key, concreteValue.Index(0)))
			return sb.String()
		}
		return encodeSliceVariable(key, concreteValue)
	default:
		return encodeStringVariable(key, concreteValue)
	}
}
