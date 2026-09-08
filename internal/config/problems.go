package config

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Problem struct {
	File    string `json:"file"`
	Line    int    `json:"line,omitempty"`
	Message string `json:"message"`
}

type Problems []Problem

func (ps Problems) Error() string {
	lines := make([]string, len(ps))
	for i, p := range ps {
		if p.Line == 0 {
			lines[i] = p.File + ": " + p.Message
			continue
		}
		lines[i] = fmt.Sprintf("%s:%d: %s", p.File, p.Line, p.Message)
	}
	return strings.Join(lines, "\n")
}

func (ps *Problems) add(file string, node *yaml.Node, name, format string, args ...any) {
	p := Problem{File: file, Message: fmt.Sprintf(format, args...)}
	if node != nil {
		p.Line = node.Line
	}
	if name != "" {
		p.Message = name + ": " + p.Message
	}
	*ps = append(*ps, p)
}

var (
	linePrefix   = regexp.MustCompile(`^line (\d+): `)
	unknownField = regexp.MustCompile(`field (\S+) not found in type \S+`)
	wrongType    = regexp.MustCompile("cannot unmarshal !!\\w+(?: `([^`]*)`)? into (\\S+)")
)

func (ps *Problems) addYAML(file string, err error) {
	msgs := []string{strings.TrimPrefix(err.Error(), "yaml: ")}
	var typeErr *yaml.TypeError
	if errors.As(err, &typeErr) {
		msgs = typeErr.Errors
	}
	for _, msg := range msgs {
		p := Problem{File: file}
		if m := linePrefix.FindStringSubmatch(msg); m != nil {
			p.Line, _ = strconv.Atoi(m[1])
			msg = msg[len(m[0]):]
		}
		msg = unknownField.ReplaceAllString(msg, `unknown setting "$1"`)
		if m := wrongType.FindStringSubmatch(msg); m != nil {
			expected := "expected " + describe(m[2])
			if m[1] != "" {
				expected += ", got " + strconv.Quote(m[1])
			}
			msg = strings.Replace(msg, m[0], expected, 1)
		}
		p.Message = msg
		*ps = append(*ps, p)
	}
}

func (ps Problems) sort() {
	slices.SortStableFunc(ps, func(a, b Problem) int {
		return cmp.Or(
			cmp.Compare(slices.Index(files, a.File), slices.Index(files, b.File)),
			cmp.Compare(a.Line, b.Line),
		)
	})
}

func describe(goType string) string {
	t := strings.TrimPrefix(goType, "*")
	switch {
	case strings.HasPrefix(t, "[]"):
		return "a list"
	case t == "int":
		return "a whole number"
	case t == "bool":
		return "true or false"
	case t == "string":
		return "text"
	}
	return "a group of settings"
}
