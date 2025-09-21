package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/registry"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/Maru-Yasa/gosong/pkg/unixsocket"
)

type ProcessCommandAction string

const (
	ProcessStop         ProcessCommandAction = "daemon-stop"
	ProcessActionStart  ProcessCommandAction = "app-start"
	ProcessActionStatus ProcessCommandAction = "app-status"
	ProcessActionStop   ProcessCommandAction = "app-stop"
	ProcessActionPing   ProcessCommandAction = "ping"
)

type Process struct {
	sockFile string
	server   *unixsocket.Server
	repo     registry.Repository
}

var ProcessCfg config.ProcessConfig = config.ProcessConfig{
	AgentPath:    "/opt/gosong",
	AppsPath:     "/var/lib/gosong/apps",
	SockFilePath: "/tmp/gosong.sock",
}

func New() *Process {
	return &Process{
		sockFile: ProcessCfg.SockFilePath,
		repo:     registry.NewFileRepository(ProcessCfg.AppsPath),
	}
}

// run the daemon that run unix socket server
func (d *Process) Run() error {
	var err error
	logger.Info("Starting gosong daemon")

	// revive the registered apps
	d.reviveApps()

	// start the unix socket server
	d.server, err = unixsocket.NewServer(d.sockFile)
	if err != nil {
		return err
	}

	err = d.server.Start()
	if err != nil {
		return err
	}

	return d.server.Accept(d.handleCommand)
}

func (d *Process) reviveApps() error {
	logger.Info("Reviving registered apps")
	apps, err := d.repo.FindAll()
	if err != nil {
		return err
	}

	for _, app := range apps {
		if app.Status == registry.AppStateStatusRunning {
			// check if the process is still running
			proc, err := os.FindProcess(app.LastPID)
			if err == nil {
				// send signal 0 to check if the process is running
				err = proc.Signal(syscall.Signal(0))
				if err == nil {
					// process is still running, skip starting it again
					continue
				}
			}

			// start the app again
			logger.Info("Starting app %s...", app.Name)
			result, err := d.start(app.Name, app.Bin, app.Port, app.Args...)
			if err != nil {
				fmt.Printf("failed to start app %s: %s\n", app.Name, err.Error())
			}

			logger.Info(result)
		}
	}

	return nil
}

func (d *Process) handleCommand(rawCmd map[string]string) (string, error) {
	fmt.Printf("%+v\n", rawCmd)

	// convert port string as integer
	port, err := strconv.Atoi(rawCmd["port"])
	if err != nil {
		port = 0
	}

	// split args
	args := []string{}
	if rawCmd["args"] != "" {
		args = strings.Split(rawCmd["args"], ",")
	}

	action := ProcessCommandAction(rawCmd["action"])

	switch action {
	case ProcessActionStart:
		result, err := d.start(rawCmd["app"], rawCmd["bin"], port, args...)
		if err != nil {
			return fmt.Sprintf("failed to start app: %s\n", err.Error()), nil
		}

		return result, nil
	case ProcessActionStop:
		return d.stop(rawCmd["app"])
	case ProcessActionStatus:
		return d.status(rawCmd["app"])
	case ProcessStop:
		// we cant actually exit here because its running in a goroutine
		// ill use systemd to manage the daemon process instead wkwk
		return "terminating\n", nil
	case ProcessActionPing:
		return "pong\n", nil
	default:
		return "unknown command\n", nil
	}
}

func (d *Process) start(appName, bin string, port int, args ...string) (string, error) {
	// check if app already running
	app, _ := d.repo.Find(appName)

	if app != nil && app.Status == registry.AppStateStatusRunning {
		// check if app actually running
		proc, err := os.FindProcess(app.LastPID)
		if err == nil {
			// send signal 0 to check if the process is running
			err = proc.Signal(syscall.Signal(0))
			if err == nil {
				// process is still running, skip starting it again
				return fmt.Sprintf("app %s already running with PID %d\n", app.Name, app.LastPID), nil
			}
		}
	}

	if app != nil {
		// use saved bin, args, port if not provided
		bin = app.Bin
		args = app.Args
		port = app.Port
	}

	// check if binary exists
	//if _, err := os.Stat(bin); err != nil {
	//	return "an error occured when searching app binary", fmt.Errorf("binary %s does not exist\n", bin)
	//}

	// look for the command path
	cPath, err := exec.LookPath(bin)
	if err != nil {
		return "an error occured when starting the app", fmt.Errorf("binary not found in PATH")
	}

	// some settings before start the process
	cmd := exec.Command(cPath, args...)
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// start the process
	if err := cmd.Start(); err != nil {
		return "an error occured when starting the app", err
	}

	// update app state
	if err := d.repo.Save(registry.AppState{
		Name:    appName,
		Bin:     bin,
		Args:    args,
		Port:    port,
		LastPID: cmd.Process.Pid,
		Status:  registry.AppStateStatusRunning,
	}); err != nil {
		// kill the process if saving state failed
		_ = cmd.Process.Kill()
		return "an error occured when saving app state", err
	}

	return fmt.Sprintf("started %s with PID %d and %d PORT \n", appName, cmd.Process.Pid, port), nil
}

func (d *Process) stop(appName string) (string, error) {
	// check if app is exists
	app, _ := d.repo.Find(appName)

	// check if app is running
	if app == nil || app.Status != registry.AppStateStatusRunning {
		return fmt.Sprintf("app %s is not running\n", appName), nil
	}

	// find the process
	proc, err := os.FindProcess(app.LastPID)
	if err != nil {
		return "an error occured when trying to find the process", err
	}

	// send SIGTERM to the process
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return "an error occured when trying to stop the process", err
	}

	// update app state
	app.Status = registry.AppStateStatusStopped
	if err := d.repo.Save(*app); err != nil {
		return "an error occured when saving app state", err
	}

	return fmt.Sprintf("stopped %s \n", appName), nil
}

func (d *Process) status(appName string) (string, error) {
	// check if app is exists
	app, _ := d.repo.Find(appName)

	if app == nil {
		return fmt.Sprintf("app %s not found\n", appName), nil
	}

	return fmt.Sprintf("app %s is %s (PID: %d, PORT: %d)\n", app.Name, app.Status, app.LastPID, app.Port), nil
}
