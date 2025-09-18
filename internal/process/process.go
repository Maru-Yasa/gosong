package process

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func Start(appName, bin string, port int, args ...string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println(appName, "duplicate process")
		return
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("failed start:", err)
		return
	}

	os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
	os.WriteFile(portFile, []byte(strconv.Itoa(port)), 0644)

	fmt.Printf("Started %s: PID %d on port %d\n", appName, cmd.Process.Pid, port)
}

func Stop(appName string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	data, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println(appName, "not found")
		return
	}

	pid, _ := strconv.Atoi(string(data))
	proc, err := os.FindProcess(pid)
	if err == nil {
		proc.Signal(syscall.SIGTERM)
	}
	os.Remove(pidFile)
	os.Remove(portFile)

	fmt.Printf("Stopped %s (PID %d)\n", appName, pid)
}

func Status(appName string) {
	pidFile := appName + ".pid"
	portFile := appName + ".port"

	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println(appName, "not running yet")
		return
	}
	portData, _ := os.ReadFile(portFile)

	pid, _ := strconv.Atoi(string(pidData))
	port, _ := strconv.Atoi(string(portData))

	fmt.Printf("started %s: PID %d on port %d\n", appName, pid, port)
}
