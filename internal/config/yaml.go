package config

import (
	"bytes"
	"errors"
	"io"

	"go.yaml.in/yaml/v3"
)

func decode(file string, data []byte, v any, problems *Problems) *yaml.Node {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(v); err != nil && !errors.Is(err, io.EOF) {
		problems.addYAML(file, err)
		var typeErr *yaml.TypeError
		if !errors.As(err, &typeErr) {
			return nil
		}
	}
	var doc yaml.Node
	if yaml.Unmarshal(data, &doc) != nil || len(doc.Content) == 0 {
		return nil
	}
	return doc.Content[0]
}

func entry(n *yaml.Node, key string) (k, v *yaml.Node) {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil, nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i], n.Content[i+1]
		}
	}
	return nil, nil
}

func value(n *yaml.Node, key string) *yaml.Node {
	_, v := entry(n, key)
	return v
}

func child(n *yaml.Node, i int) *yaml.Node {
	if n == nil || n.Kind != yaml.SequenceNode || i >= len(n.Content) {
		return nil
	}
	return n.Content[i]
}
