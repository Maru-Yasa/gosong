package executor

import (
	"errors"
	"strings"
	"testing"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/tasks"
)

// mockExecutor records calls and responds via a callback function.
type mockExecutor struct {
	name  string
	calls []string
	fn    func(cmd string) (string, error)
}

func (m *mockExecutor) GetName() string { return m.name }
func (m *mockExecutor) RunRaw(cmd string) (string, error) {
	m.calls = append(m.calls, cmd)
	if m.fn != nil {
		return m.fn(cmd)
	}
	return "", nil
}
func (m *mockExecutor) Run(cmd, cwd string) (string, error) {
	return m.RunRaw(cmd)
}

// --- GetGreenPort ---

func TestGetGreenPort_ReturnOtherPort(t *testing.T) {
	if got := GetGreenPort([]int{3000, 3001}, 3000); got != 3001 {
		t.Errorf("got %d, want 3001", got)
	}
	if got := GetGreenPort([]int{3000, 3001}, 3001); got != 3000 {
		t.Errorf("got %d, want 3000", got)
	}
}

func TestGetGreenPort_FallbackWhenNoOther(t *testing.T) {
	// single port or active not in list: returns activePort
	if got := GetGreenPort([]int{3000}, 3000); got != 3000 {
		t.Errorf("expected fallback to activePort 3000, got %d", got)
	}
}

func TestGetGreenPort_EmptyPorts(t *testing.T) {
	if got := GetGreenPort([]int{}, 3000); got != 3000 {
		t.Errorf("expected fallback to activePort 3000, got %d", got)
	}
}

// --- GetActivePort ---

func TestGetActivePort_ReadsFromStateFile(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "3001\n", nil
	}}
	port, err := GetActivePort(exec, "/opt/myapp", "myapp", []int{3000, 3001})
	if err != nil {
		t.Fatal(err)
	}
	if port != 3001 {
		t.Errorf("got %d, want 3001", port)
	}
}

func TestGetActivePort_FallbackOnInvalidOutput(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "not-a-number\n", nil
	}}
	port, err := GetActivePort(exec, "/opt/myapp", "myapp", []int{3000, 3001})
	if err != nil {
		t.Fatal(err)
	}
	if port != 3000 {
		t.Errorf("expected fallback to ports[0]=3000, got %d", port)
	}
}

func TestGetActivePort_FallbackOnExecError(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "", errors.New("ssh error")
	}}
	port, err := GetActivePort(exec, "/opt/myapp", "myapp", []int{3000, 3001})
	if err != nil {
		t.Fatal(err)
	}
	if port != 3000 {
		t.Errorf("expected fallback to ports[0]=3000, got %d", port)
	}
}

func TestGetActivePort_EmptyPorts(t *testing.T) {
	exec := &mockExecutor{}
	port, err := GetActivePort(exec, "/opt/myapp", "myapp", []int{})
	if err != nil {
		t.Fatal(err)
	}
	if port != 0 {
		t.Errorf("expected 0 for empty ports, got %d", port)
	}
}

// --- GetLastIDFromHost ---

func TestGetLastIDFromHost_ReturnsID(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "5\n", nil
	}}
	id, err := GetLastIDFromHost(exec, "/opt/myapp")
	if err != nil {
		t.Fatal(err)
	}
	if id != 5 {
		t.Errorf("got %d, want 5", id)
	}
}

func TestGetLastIDFromHost_ReturnsZeroWhenNone(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "0\n", nil
	}}
	id, err := GetLastIDFromHost(exec, "/opt/myapp")
	if err != nil {
		t.Fatal(err)
	}
	if id != 0 {
		t.Errorf("got %d, want 0", id)
	}
}

func TestGetLastIDFromHost_ErrorOnInvalidOutput(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "not-a-number\n", nil
	}}
	_, err := GetLastIDFromHost(exec, "/opt/myapp")
	if err == nil {
		t.Fatal("expected error for non-numeric output")
	}
}

// --- RunTask ---

func TestRunTask_RunStep(t *testing.T) {
	exec := &mockExecutor{}
	err := RunTask(RunTaskParams{
		Exec:     exec,
		Cfg:      &config.ConfigRoot{},
		CfgMap:   map[string]any{"Greeting": "hello"},
		TaskName: "test",
		UTask: map[string]tasks.Task{
			"test": {Steps: []tasks.Step{{Run: "echo {{.Greeting}}"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(exec.calls) != 1 || exec.calls[0] != "echo hello" {
		t.Errorf("unexpected calls: %v", exec.calls)
	}
}

func TestRunTask_CdStep(t *testing.T) {
	called := false
	exec := &mockExecutor{
		fn: func(cmd string) (string, error) {
			called = true
			return "", nil
		},
	}
	err := RunTask(RunTaskParams{
		Exec:   exec,
		Cfg:    &config.ConfigRoot{},
		CfgMap: map[string]any{"Dir": "/tmp"},
		UTask: map[string]tasks.Task{
			"test": {Steps: []tasks.Step{
				{Cd: "{{.Dir}}"},
				{Run: "echo ok"},
			}},
		},
		TaskName: "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("run step was not executed after cd step")
	}
}

func TestRunTask_NestedTask(t *testing.T) {
	exec := &mockExecutor{}
	err := RunTask(RunTaskParams{
		Exec:     exec,
		Cfg:      &config.ConfigRoot{},
		CfgMap:   map[string]any{},
		TaskName: "outer",
		UTask: map[string]tasks.Task{
			"outer": {Steps: []tasks.Step{{Task: "inner"}}},
			"inner": {Steps: []tasks.Step{{Run: "echo from-inner"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(exec.calls) != 1 || exec.calls[0] != "echo from-inner" {
		t.Errorf("expected inner task run, got: %v", exec.calls)
	}
}

func TestRunTask_PropagatesRunError(t *testing.T) {
	exec := &mockExecutor{fn: func(cmd string) (string, error) {
		return "", errors.New("command failed")
	}}
	err := RunTask(RunTaskParams{
		Exec:     exec,
		Cfg:      &config.ConfigRoot{},
		CfgMap:   map[string]any{},
		TaskName: "test",
		UTask: map[string]tasks.Task{
			"test": {Steps: []tasks.Step{{Run: "failing-cmd"}}},
		},
	})
	if err == nil {
		t.Fatal("expected error to propagate from Run()")
	}
}

func TestRunTask_TaskNotFound(t *testing.T) {
	exec := &mockExecutor{}
	err := RunTask(RunTaskParams{
		Exec:     exec,
		Cfg:      &config.ConfigRoot{},
		CfgMap:   map[string]any{},
		TaskName: "nonexistent",
		UTask:    map[string]tasks.Task{},
	})
	if err == nil {
		t.Fatal("expected error for unknown task")
	}
}

func TestRunTask_TemplateRenderedBeforeExec(t *testing.T) {
	exec := &mockExecutor{}
	err := RunTask(RunTaskParams{
		Exec:     exec,
		Cfg:      &config.ConfigRoot{},
		CfgMap:   map[string]any{"ReleasePath": "/opt/app/releases/42"},
		TaskName: "deploy",
		UTask: map[string]tasks.Task{
			"deploy": {Steps: []tasks.Step{{Run: "cd {{.ReleasePath}} && npm install"}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(exec.calls[0], "/opt/app/releases/42") {
		t.Errorf("template not rendered in command: %q", exec.calls[0])
	}
}

// --- LocalExecutor ---

func TestLocalExecutor_GetName(t *testing.T) {
	e, err := newLocalExecutor("test-host")
	if err != nil {
		t.Fatal(err)
	}
	if e.GetName() != "test-host" {
		t.Errorf("got %q, want %q", e.GetName(), "test-host")
	}
}

func TestLocalExecutor_RunEcho(t *testing.T) {
	e, _ := newLocalExecutor("local")
	out, err := e.Run("echo hello", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("expected 'hello' in output, got %q", out)
	}
}

func TestLocalExecutor_RunRaw(t *testing.T) {
	e, _ := newLocalExecutor("local")
	out, err := e.RunRaw("echo rawtest")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "rawtest") {
		t.Errorf("expected 'rawtest' in output, got %q", out)
	}
}

func TestLocalExecutor_CommandNotFound(t *testing.T) {
	e, _ := newLocalExecutor("local")
	_, err := e.Run("this-command-does-not-exist-xyz", "")
	if err == nil {
		t.Fatal("expected error for nonexistent command")
	}
}
