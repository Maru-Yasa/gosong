package templateutil

import (
	"bytes"
	"encoding/json"
	"text/template"
)

func RenderTemplate(input string, vars map[string]any) (string, error) {
	tmpl, err := template.New("").Option("missingkey=zero").Parse(input)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ToMap(input any) (map[string]any, error) {
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
