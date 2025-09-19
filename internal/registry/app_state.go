package registry

type AppStateStatus string

const (
	AppStateStatusRunning AppStateStatus = "running"
	AppStateStatusStopped AppStateStatus = "stopped"
	AppStateStatusError   AppStateStatus = "error"
)

type AppState struct {
	Name    string         `json:"name"`
	Bin     string         `json:"bin"`
	Args    []string       `json:"args"`
	Port    int            `json:"port"`
	LastPID int            `json:"last_pid"`
	Status  AppStateStatus `json:"status"`
}
