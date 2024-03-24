package utils

import (
	"encoding/json"

	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

func FromJson(d []byte) (*structpb.Value, error) {
	var res structpb.Value
	if err := json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func FromYaml(d []byte) (*structpb.Value, error) {
	var yd map[string]any
	if err := yaml.Unmarshal(d, &yd); err != nil {
		return nil, err
	}
	return structpb.NewValue(yd)
}

func FromPlain(d []byte) (*structpb.Value, error) {
	return structpb.NewValue(d)
}
