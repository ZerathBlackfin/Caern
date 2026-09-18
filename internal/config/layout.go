package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Layout struct {
	Blocks []LayoutBlock `json:"blocks"`
}

type LayoutBlock struct {
	ID    string       `json:"id"`
	Tiles []LayoutTile `json:"tiles"`
}

type LayoutTile struct {
	ID   string `json:"id"`
	Size string `json:"size"`
}

func (s *Store) SaveLayout(l Layout) error {
	path := filepath.Join(s.dir, servicesFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.SequenceNode {
		return errors.New("services.yaml is not a list")
	}

	root := doc.Content[0]
	items, err := l.items(root)
	if err != nil {
		return err
	}
	root.Content = items

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	if err := replaceFile(path, spaceItems(buf.Bytes())); err != nil {
		return err
	}
	s.reload()
	return nil
}

func (l Layout) items(root *yaml.Node) ([]*yaml.Node, error) {
	taken := map[string]bool{}
	find := func(id string) (*yaml.Node, error) {
		if taken[id] {
			return nil, fmt.Errorf("%q is used twice", id)
		}
		taken[id] = true

		head, tail, nested := strings.Cut(id, ".")
		i, err := strconv.Atoi(head)
		if err != nil || i < 0 || i >= len(root.Content) {
			return nil, fmt.Errorf("unknown tile %q, reload the page", id)
		}
		node := root.Content[i]
		if !nested {
			return node, nil
		}
		j, err := strconv.Atoi(tail)
		list := value(node, "items")
		if err != nil || list == nil || j < 0 || j >= len(list.Content) {
			return nil, fmt.Errorf("unknown tile %q, reload the page", id)
		}
		return list.Content[j], nil
	}

	var items []*yaml.Node
	for _, b := range l.Blocks {
		tiles := make([]*yaml.Node, 0, len(b.Tiles))
		for _, t := range b.Tiles {
			node, err := find(t.ID)
			if err != nil {
				return nil, err
			}
			if err := setSize(node, t.Size); err != nil {
				return nil, err
			}
			tiles = append(tiles, node)
		}
		if b.ID == "" {
			items = append(items, tiles...)
			continue
		}
		section, err := find(b.ID)
		if err != nil {
			return nil, err
		}
		setKeyNode(section, "items", &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: tiles})
		items = append(items, section)
	}
	return items, nil
}

func setSize(node *yaml.Node, size string) error {
	if size == "" || size == "1x1" {
		removeKey(node, "size")
		return nil
	}
	if _, _, err := parseSize(size); err != nil {
		return err
	}
	setKeyNode(node, "size", &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: size})
	return nil
}

func setKeyNode(node *yaml.Node, key string, val *yaml.Node) {
	if _, old := entry(node, key); old != nil {
		val.HeadComment, val.LineComment, val.FootComment = old.HeadComment, old.LineComment, old.FootComment
		*old = *val
		return
	}
	node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, val)
}

func removeKey(node *yaml.Node, key string) {
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			node.Content = slices.Delete(node.Content, i, i+2)
			return
		}
	}
}

func spaceItems(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	out := make([]string, 0, len(lines))
	for i, line := range lines {
		if i > 0 && strings.HasPrefix(line, "- ") {
			start := len(out)
			for start > 0 && strings.HasPrefix(out[start-1], "#") {
				start--
			}
			if start > 0 && out[start-1] != "" {
				out = slices.Insert(out, start, "")
			}
		}
		out = append(out, line)
	}
	return []byte(strings.Join(out, "\n"))
}

func replaceFile(path string, data []byte) error {
	mode := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
