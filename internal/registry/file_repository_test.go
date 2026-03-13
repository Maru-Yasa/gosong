package registry

import (
	"testing"
)

func newTestRepo(t *testing.T) *FileRepository {
	t.Helper()
	return NewFileRepository(t.TempDir())
}

func TestFileRepository_SaveAndFind(t *testing.T) {
	repo := newTestRepo(t)
	app := AppState{
		Name:    "myapp",
		Bin:     "node",
		Args:    []string{"server.js"},
		Port:    3000,
		LastPID: 12345,
		Status:  AppStateStatusRunning,
	}

	if err := repo.Save(app); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Find("myapp")
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	if got.Name != app.Name {
		t.Errorf("Name: got %q, want %q", got.Name, app.Name)
	}
	if got.Port != app.Port {
		t.Errorf("Port: got %d, want %d", got.Port, app.Port)
	}
	if got.LastPID != app.LastPID {
		t.Errorf("LastPID: got %d, want %d", got.LastPID, app.LastPID)
	}
	if got.Status != AppStateStatusRunning {
		t.Errorf("Status: got %q, want %q", got.Status, AppStateStatusRunning)
	}
}

func TestFileRepository_Find_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	_, err := repo.Find("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent app")
	}
}

func TestFileRepository_Save_Overwrites(t *testing.T) {
	repo := newTestRepo(t)
	app := AppState{Name: "myapp", Status: AppStateStatusRunning, LastPID: 100}
	_ = repo.Save(app)

	app.Status = AppStateStatusStopped
	app.LastPID = 0
	if err := repo.Save(app); err != nil {
		t.Fatal(err)
	}

	got, _ := repo.Find("myapp")
	if got.Status != AppStateStatusStopped {
		t.Errorf("Status: got %q, want %q", got.Status, AppStateStatusStopped)
	}
}

func TestFileRepository_FindAll_Empty(t *testing.T) {
	repo := newTestRepo(t)
	apps, err := repo.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 0 {
		t.Errorf("expected empty slice, got %d items", len(apps))
	}
}

func TestFileRepository_FindAll_Multiple(t *testing.T) {
	repo := newTestRepo(t)
	_ = repo.Save(AppState{Name: "app1", Status: AppStateStatusRunning})
	_ = repo.Save(AppState{Name: "app2", Status: AppStateStatusStopped})

	apps, err := repo.FindAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Errorf("expected 2 apps, got %d", len(apps))
	}
}

func TestFileRepository_Delete(t *testing.T) {
	repo := newTestRepo(t)
	_ = repo.Save(AppState{Name: "myapp"})

	if err := repo.Delete("myapp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := repo.Find("myapp")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestFileRepository_Delete_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	err := repo.Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error when deleting nonexistent app")
	}
}

func TestFileRepository_ArgsPreserved(t *testing.T) {
	repo := newTestRepo(t)
	args := []string{"--port", "3000", "--env", "production"}
	_ = repo.Save(AppState{Name: "myapp", Args: args})

	got, _ := repo.Find("myapp")
	if len(got.Args) != len(args) {
		t.Fatalf("Args len: got %d, want %d", len(got.Args), len(args))
	}
	for i, a := range args {
		if got.Args[i] != a {
			t.Errorf("Args[%d]: got %q, want %q", i, got.Args[i], a)
		}
	}
}
