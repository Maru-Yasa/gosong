package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/pkg/unixsocket"
)

type ProcessCommandAction string

const (
	ProcessStop        ProcessCommandAction = "daemon-stop"
	ProcessActionStart ProcessCommandAction = "app-start"
	ProcessActionStop  ProcessCommandAction = "app-stop"
	ProcessActionPing  ProcessCommandAction = "ping"
)

type Process struct {
	sockFile string
	server   *unixsocket.Server
}

var ProcessCfg config.ProcessConfig = config.ProcessConfig{
	AgentPath:   "/opt/gosong",
	ProcessPath: "/var/lib/gosong/apps",
}

func New() *Process {
	return &Process{
		sockFile: "/tmp/gosong.sock",
	}
}

func (d *Process) Run() error {
	var err error
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
		Start(rawCmd["app"], rawCmd["bin"], port, args...)
		return "started\n", nil
	case ProcessActionStop:
		Stop(rawCmd["app"])
		return fmt.Sprintf("stopped %s \n", rawCmd["app"]), nil
	case ProcessStop:
		// We can't actually exit here because it's running in a goroutine
		// Instead, we'll return a response and let the caller handle the exit
		return "terminating\n", nil
	case ProcessActionPing:
		return "pong\n", nil
	default:
		return "unknown command\n", nil
	}
}

func Start(appName, bin string, port int, args ...string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("DEBUG: pidFile exists, skipping start")
		return
	}

	fmt.Println("DEBUG bin :", bin)
	fmt.Println("DEBUG port :", port)
	fmt.Println("DEBUG args:", args)

	if _, err := os.Stat(bin); err != nil {
		fmt.Println("DEBUG bin not found:", err)
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("DEBUG running command:", bin, strings.Join(args, " "))

	if err := cmd.Start(); err != nil {
		fmt.Printf("ERROR failed to start: %s\n", err)
		return
	}

	fmt.Println("DEBUG process started with PID:", cmd.Process.Pid)

	_ = os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	_ = os.WriteFile(portFile, []byte(strconv.Itoa(port)), 0644)
}

func Stop(appName string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	pid, _ := strconv.Atoi(string(data))
	proc, err := os.FindProcess(pid)

	if err != nil {
		fmt.Printf("%v", err)
	}

	proc.Signal(syscall.SIGTERM)
	os.Remove(pidFile)
	os.Remove(portFile)
}

func Status(appName string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return
	}
	_, _ = os.ReadFile(portFile)

	_, _ = strconv.Atoi(string(pidData))
	// port, _ := strconv.Atoi(string(portData))
}
