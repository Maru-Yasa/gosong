package templateutil

import (
	"testing"
)

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	got, err := RenderTemplate("hello {{.Name}}", map[string]any{"Name": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestRenderTemplate_MultipleVars(t *testing.T) {
	got, err := RenderTemplate("{{.A}}-{{.B}}", map[string]any{"A": "foo", "B": "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "foo-bar" {
		t.Errorf("got %q, want %q", got, "foo-bar")
	}
}

func TestRenderTemplate_IntVar(t *testing.T) {
	got, err := RenderTemplate("port {{.Port}}", map[string]any{"Port": 3000})
	if err != nil {
		t.Fatal(err)
	}
	if got != "port 3000" {
		t.Errorf("got %q, want %q", got, "port 3000")
	}
}

func TestRenderTemplate_EmptyTemplate(t *testing.T) {
	got, err := RenderTemplate("", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestRenderTemplate_MissingKey_NoError(t *testing.T) {
	// missingkey=zero: missing keys must not return an error
	_, err := RenderTemplate("value={{.Missing}}", map[string]any{})
	if err != nil {
		t.Fatalf("expected no error with missingkey=zero, got: %v", err)
	}
}

func TestRenderTemplate_InvalidSyntax(t *testing.T) {
	_, err := RenderTemplate("{{.Unclosed", map[string]any{})
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestRenderTemplate_NestedVar(t *testing.T) {
	type Inner struct{ Name string }
	got, err := RenderTemplate("{{.App.Name}}", map[string]any{"App": Inner{Name: "myapp"}})
	if err != nil {
		t.Fatal(err)
	}
	if got != "myapp" {
		t.Errorf("got %q, want %q", got, "myapp")
	}
}

func TestToMap_BasicStruct(t *testing.T) {
	type Foo struct {
		Name string
		Port int
	}
	m, err := ToMap(Foo{Name: "myapp", Port: 3000})
	if err != nil {
		t.Fatal(err)
	}
	if m["Name"] != "myapp" {
		t.Errorf("Name: got %v, want %v", m["Name"], "myapp")
	}
	if m["Port"] != 3000 {
		t.Errorf("Port: got %v, want %v", m["Port"], 3000)
	}
}

func TestToMap_PointerStruct(t *testing.T) {
	type Foo struct{ Name string }
	m, err := ToMap(&Foo{Name: "ptr-test"})
	if err != nil {
		t.Fatal(err)
	}
	if m["Name"] != "ptr-test" {
		t.Errorf("Name: got %v, want %v", m["Name"], "ptr-test")
	}
}

func TestToMap_UsesFieldName_NotYamlTag(t *testing.T) {
	type Cfg struct {
		AppPath string `yaml:"app_path"`
	}
	m, err := ToMap(Cfg{AppPath: "/opt/myapp"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := m["AppPath"]; !ok {
		t.Error("expected key 'AppPath' (field name), not yaml tag 'app_path'")
	}
	if _, ok := m["app_path"]; ok {
		t.Error("did not expect yaml tag key 'app_path'")
	}
}

func TestToMap_EmptyStruct(t *testing.T) {
	type Empty struct{}
	m, err := ToMap(Empty{})
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}
