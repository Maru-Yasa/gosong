package daemon

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/Maru-Yasa/gosong/pkg/proto/kvproto"
)

type ProcessCommandAction string

const (
	ProcessStop        ProcessCommandAction = "daemon-stop"
	ProcessActionStart ProcessCommandAction = "app-start"
	ProcessActionStop  ProcessCommandAction = "app-stop"
	ProcessActionPing  ProcessCommandAction = "ping"
)

type ProcessCommand struct {
	Action  ProcessCommandAction
	AppName string
	Port    uint8
	Bin     string
	Args    []string
}

type Process struct {
	Network  string
	SockFile string
}

func New() *Process {
	return &Process{
		Network:  "unix",
		SockFile: "/tmp/gosong.sock",
	}
}

func (d *Process) Run() error {
	// delete previous unix socket
	_ = os.Remove(d.SockFile)

	listener, err := net.Listen(d.Network, d.SockFile)

	if err != nil {
		return fmt.Errorf("can't connect unix socker %s", d.SockFile)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go d.handleConnection(conn)
	}

}

func (d *Process) handleConnection(conn net.Conn) error {
	// this should handle connection from unix socket
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)

	// lazy to use buffer so convert to string
	msg := string(buf[:n])

	rawCmd := kvproto.Decode(msg)

	// convert port string as integer
	port, err := strconv.Atoi(rawCmd["port"])

	if err != nil {
		port = 0
	}

	defer conn.Close()

	// split args
	args := []string{}
	args = strings.Split(rawCmd["args"], ",")

	dCmd := ProcessCommand{
		Action:  ProcessCommandAction(rawCmd["action"]),
		AppName: rawCmd["app"],
		Port:    uint8(port),
		Bin:     rawCmd["bin"],
		Args:    args,
	}

	switch dCmd.Action {
	case ProcessActionStart:
		Start(dCmd.AppName, dCmd.Bin, int(dCmd.Port), dCmd.Args...)
	case ProcessActionStop:
		Stop(dCmd.AppName)
	case ProcessStop:
		os.Exit(0)
	case ProcessActionPing:
		_, _ = conn.Write([]byte("pong\n"))
	}

	return nil
}

func Start(appName, bin string, port int, args ...string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	if _, err := os.Stat(pidFile); err == nil {
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return
	}

	os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	os.WriteFile(portFile, []byte(strconv.Itoa(port)), 0644)
}

func Stop(appName string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	data, err := os.ReadFile(pidFile)
	if err != nil {
		return
	}

	pid, _ := strconv.Atoi(string(data))
	proc, err := os.FindProcess(pid)
	if err == nil {
		proc.Signal(syscall.SIGTERM)
	}
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
