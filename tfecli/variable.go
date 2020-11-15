package tfecli

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

// HCLVariable defines a HCLVariable.
type HCLVariable struct {
	key   string
	value string
}

func (h HCLVariable) String() string {
	return fmt.Sprintf("%s=%s", h.key, h.value)
}

// ParseVarFile reads an HCL varfile and returns its content as a list of HCL variables.
func ParseVarFile(varFile string) ([]HCLVariable, error) {
	fileContent, err := ioutil.ReadFile(varFile)
	if err != nil {
		return []HCLVariable{}, fmt.Errorf("cannot read the file %q: %s", varFile, err)
	}

	return SimpleParseVarFile(string(fileContent)), nil
}

// SimpleParseVarFile parses a var file and returns a list of HCLVariables.
func SimpleParseVarFile(varFile string) []HCLVariable {
	re := regexp.MustCompile(`(?im)^([a-z0-9_]*)\s*=\s*`)
	matches := re.FindAllStringSubmatch(varFile, -1)
	splits := re.Split(varFile, -1)
	HCLVariables := []HCLVariable{}
	for i, match := range matches {
		v := HCLVariable{
			key:   match[1],
			value: splits[i+1],
		}
		HCLVariables = append(HCLVariables, v)
	}
	return HCLVariables
}
