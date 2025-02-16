package assets

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"

	"sigs.k8s.io/yaml"
)

var TemplateFuncs template.FuncMap = template.FuncMap{
	"toYAML":     MarshalYAML,
	"indent":     Indent,
	"nindent":    NIndent,
	"indentNext": IndentNext,
	"toBytes":    ToBytes,
	"toBase64":   ToBase64,
	"map":        MakeMap,
	"repeat":     Repeat,
}

func MarshalYAML(v any) (string, error) {
	bytes, err := yaml.Marshal(v)
	return strings.TrimSpace(string(bytes)), err
}

func Indent(spaceCount int, s string) string {
	spaces := strings.Repeat(" ", spaceCount)
	return spaces + strings.Replace(s, "\n", "\n"+spaces, -1)
}

func NIndent(spaceCount int, s string) string {
	return "\n" + Indent(spaceCount, s)
}

func IndentNext(spaceCount int, s string) string {
	parts := strings.SplitAfterN(s, "\n", 2)
	if len(parts) == 1 {
		return parts[0]
	}
	return parts[0] + Indent(spaceCount, parts[1])
}

func Repeat(s string, count int) string {
	var sb strings.Builder
	sb.Grow(len(s) * count)
	for i := 0; i < count; i++ {
		sb.WriteString(s)
	}
	return sb.String()
}

func ToBytes(s string) []byte {
	return []byte(s)
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func MakeMap(kvs ...any) (map[any]any, error) {
	count := len(kvs)
	if count%2 != 0 {
		return nil, fmt.Errorf("map length %d isn't dividable into tuples", count)
	}

	m := make(map[any]any, count%2)
	for i := 0; i+1 < count; i += 2 {
		m[kvs[i]] = kvs[i+1]
	}

	return m, nil
}

func RenderTemplate(tmpl *template.Template, inputs any) ([]byte, error) {
	// We always want correctness. (Accidentally missing a key might have side effects.)
	tmpl.Option("missingkey=error")

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, inputs)
	if err != nil {
		return nil, fmt.Errorf("can't execute template %q: %w", tmpl.Name(), err)
	}

	return buf.Bytes(), nil
}
