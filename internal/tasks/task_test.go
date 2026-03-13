package tasks

import (
	"testing"
)

func TestRegisterAndFindTask_BuiltIn(t *testing.T) {
	builtins := []string{"fetch", "lock", "unlock", "symlink", "show_releases", "start-app", "stop-old-app"}
	for _, name := range builtins {
		t.Run(name, func(t *testing.T) {
			task, err := FindTask(name, nil)
			if err != nil {
				t.Fatalf("built-in task %q not found: %v", name, err)
			}
			if len(task.Steps) == 0 {
				t.Errorf("task %q has no steps", name)
			}
		})
	}
}

func TestFindTask_UserTaskTakesPrecedence(t *testing.T) {
	custom := Task{
		Description: "custom fetch",
		Steps:       []Step{{Run: "echo custom"}},
	}
	uTasks := map[string]Task{"fetch": custom}

	task, err := FindTask("fetch", uTasks)
	if err != nil {
		t.Fatal(err)
	}
	if task.Description != "custom fetch" {
		t.Errorf("expected user task to override built-in, got %q", task.Description)
	}
}

func TestFindTask_NotFound(t *testing.T) {
	_, err := FindTask("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for unknown task")
	}
}

func TestFindTask_UserTaskOnly(t *testing.T) {
	uTasks := map[string]Task{
		"my-custom-task": {Steps: []Step{{Run: "echo hello"}}},
	}
	task, err := FindTask("my-custom-task", uTasks)
	if err != nil {
		t.Fatal(err)
	}
	if task.Steps[0].Run != "echo hello" {
		t.Errorf("unexpected step: %q", task.Steps[0].Run)
	}
}

func TestStartAppTask_HasExpectedSteps(t *testing.T) {
	task, err := FindTask("start-app", nil)
	if err != nil {
		t.Fatal(err)
	}
	// must create dirs, start app, and report
	if len(task.Steps) < 3 {
		t.Errorf("start-app: expected at least 3 steps, got %d", len(task.Steps))
	}
	// first step must create required directories
	if task.Steps[0].Run == "" {
		t.Error("start-app step[0] should be a run step")
	}
}

func TestStopOldAppTask_HasExpectedSteps(t *testing.T) {
	task, err := FindTask("stop-old-app", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(task.Steps) < 3 {
		t.Errorf("stop-old-app: expected at least 3 steps, got %d", len(task.Steps))
	}
}

func TestRegisterTask_Custom(t *testing.T) {
	// register a unique task name to avoid polluting shared built-in map
	name := "test-only-task-" + t.Name()
	RegisterTask(name, Task{Steps: []Step{{Run: "echo registered"}}})

	task, err := FindTask(name, nil)
	if err != nil {
		t.Fatalf("task registered but not findable: %v", err)
	}
	if task.Steps[0].Run != "echo registered" {
		t.Errorf("unexpected step: %q", task.Steps[0].Run)
	}
}
