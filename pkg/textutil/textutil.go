// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package textutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
)

// ExecuteTemplate executes a text/template template.
func ExecuteTemplate(tmpl string, args any) ([]byte, error) {
	x, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	if err := x.Execute(&b, args); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// PrefixString adds prefix to beginning of each line.
func PrefixString(prefix, text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

// IndentString add spaces to beginning of each line.
func IndentString(size int, text string) string {
	prefix := strings.Repeat(" ", size)
	return PrefixString(prefix, text)
}

// MissingString returns message if the text is empty.
func MissingString(message, text string) string {
	if text == "" {
		return message
	}
	return text
}

// TemplateFuncMap is a text/template FuncMap.
var TemplateFuncMap = template.FuncMap{
	"json": func(v any) string {
		var b bytes.Buffer
		enc := json.NewEncoder(&b)
		enc.SetEscapeHTML(false)
		if err := enc.Encode(v); err != nil {
			panic(fmt.Errorf("failed to marshal as JSON: %+v: %w", v, err))
		}
		return strings.TrimSuffix(b.String(), "\n")
	},
	"yaml": func(v any) string {
		var b bytes.Buffer
		enc := yaml.NewEncoder(&b)
		if err := enc.Encode(v); err != nil {
			panic(fmt.Errorf("failed to marshal as YAML: %+v: %w", v, err))
		}
		return "---\n" + strings.TrimSuffix(b.String(), "\n")
	},
	"indent": func(a ...any) (string, error) {
		if len(a) == 0 {
			return "", errors.New("function takes at least one string argument")
		}
		if len(a) > 2 {
			return "", errors.New("function takes at most 2 arguments")
		}
		var ok bool
		size := 2
		if len(a) > 1 {
			if size, ok = a[0].(int); !ok {
				return "", errors.New("optional first argument must be an integer")
			}
		}
		text := ""
		if text, ok = a[len(a)-1].(string); !ok {
			return "", errors.New("last argument must be a string")
		}
		return IndentString(size, text), nil
	},
	"missing": func(a ...any) (string, error) {
		if len(a) == 0 {
			return "", errors.New("function takes at least one string argument")
		}
		if len(a) > 2 {
			return "", errors.New("function takes at most 2 arguments")
		}
		var ok bool
		message := "<missing>"
		if len(a) > 1 {
			if message, ok = a[0].(string); !ok {
				return "", errors.New("optional first argument must be a string")
			}
		}
		text := ""
		if text, ok = a[len(a)-1].(string); !ok {
			return "", errors.New("last argument must be a string")
		}
		return MissingString(message, text), nil
	},
}

// FuncHelp is help for TemplateFuncMap.
var FuncHelp = []string{
	"indent <size>: add spaces to beginning of each line",
	"missing <message>: return message if the text is empty",
}
