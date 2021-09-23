package ast

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

func FromYAMLReader(reader io.Reader) (*Command, error) {
	decoder := yaml.NewDecoder(reader)

	cmde := &Command{}
	err := decoder.Decode(cmde)
	if err != nil {
		return nil, fmt.Errorf("failed decoding ast: %w", err)
	}
	return cmde, nil
}

func ToYAMLWriter(cmde *Command, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)

	err := encoder.Encode(cmde)
	if err != nil {
		return fmt.Errorf("failed to encode ast: %w", err)
	}
	return nil
}

func FromYAMLString(astYAML string) (*Command, error) {
	cmde := &Command{}
	err := yaml.Unmarshal([]byte(astYAML), cmde)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling ast: %w", err)
	}
	return cmde, nil
}

func ToYAMLString(cmde *Command) (string, error) {
	astYAML, err := yaml.Marshal(cmde)
	if err != nil {
		return "", fmt.Errorf("failed to marshaling ast: %w", err)
	}
	return string(astYAML), nil
}
