package agent

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/Maru-Yasa/gosong/internal/process"
	"github.com/Maru-Yasa/gosong/pkg/proto/kvproto"
)

type AgentCommandAction string

const (
	AgentStop        AgentCommandAction = "daemon-stop"
	AgentActionStart AgentCommandAction = "app-start"
	AgentActionStop  AgentCommandAction = "app-stop"
	AgentActionPing  AgentCommandAction = "ping"
)

type AgentCommand struct {
	Action  AgentCommandAction
	AppName string
	Port    uint8
	Bin     string
	Args    []string
}

type Agent struct {
	Network  string
	SockFile string
}

func New() *Agent {
	return &Agent{
		Network:  "unix",
		SockFile: "/tmp/gosong.sock",
	}
}

func (d *Agent) Run() error {
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

func (d *Agent) handleConnection(conn net.Conn) error {
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

	dCmd := AgentCommand{
		Action:  AgentCommandAction(rawCmd["action"]),
		AppName: rawCmd["app"],
		Port:    uint8(port),
		Bin:     rawCmd["bin"],
		Args:    args,
	}

	fmt.Printf("%+v\n", dCmd)

	switch dCmd.Action {
	case AgentActionStart:
		process.Start(dCmd.AppName, dCmd.Bin, int(dCmd.Port), dCmd.Args...)
	case AgentActionStop:
		process.Stop(dCmd.AppName)
	case AgentStop:
		fmt.Println("Terminating")
		os.Exit(0)
	case AgentActionPing:
		_, _ = conn.Write([]byte("pong\n"))
	}

	return nil
}
