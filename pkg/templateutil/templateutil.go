package templateutil

import (
	"bytes"
	"reflect"
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
	out := make(map[string]any)
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		out[field.Name] = value
	}
	return out, nil
}
