package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/Maru-Yasa/gosong/internal/config"
	"github.com/Maru-Yasa/gosong/internal/executor"
	"github.com/Maru-Yasa/gosong/internal/registry"
	"github.com/Maru-Yasa/gosong/pkg/logger"
	"github.com/urfave/cli/v3"
)

const appsPath = "/var/lib/gosong/apps"

func ProcessCommand() *cli.Command {
	return &cli.Command{
		Name:  "process",
		Usage: "process <command>",
		Commands: []*cli.Command{
			ProcessStartCommand(),
			ProcessStatusCommand(),
			ProcessStopCommand(),
		},
	}
}

func ProcessStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "<app_name> [--bin <bin>] [--port <port>] [--config <file>] [--host <name>]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "bin", Aliases: []string{"b"}, Usage: "Path to binary"},
			&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Usage: "Port to run the app on"},
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "Config file path (for remote hosts)"},
			&cli.StringFlag{Name: "host", Usage: "Remote host name from config"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			log := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				return fmt.Errorf("need <app_name>")
			}
			app := c.Args().First()
			bin := c.String("bin")
			port := c.Int("port")
			args := c.Args().Slice()[1:]

			if cfgFile := c.String("config"); cfgFile != "" {
				return remoteProcessStart(cfgFile, c.String("host"), app, bin, port, args)
			}

			return localProcessStart(log, app, bin, port, args)
		},
	}
}

func ProcessStatusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "<app_name> [--config <file>] [--host <name>]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "Config file path"},
			&cli.StringFlag{Name: "host", Usage: "Remote host name from config"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			log := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				return fmt.Errorf("need <app_name>")
			}
			app := c.Args().First()

			if cfgFile := c.String("config"); cfgFile != "" {
				return remoteProcessStatus(cfgFile, c.String("host"), app)
			}

			return localProcessStatus(log, app)
		},
	}
}

func ProcessStopCommand() *cli.Command {
	return &cli.Command{
		Name:  "stop",
		Usage: "<app_name> [--config <file>] [--host <name>]",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "Config file path"},
			&cli.StringFlag{Name: "host", Usage: "Remote host name from config"},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			log := logger.NewConsoleLogger()
			if c.Args().Len() < 1 {
				return fmt.Errorf("need <app_name>")
			}
			app := c.Args().First()

			if cfgFile := c.String("config"); cfgFile != "" {
				return remoteProcessStop(cfgFile, c.String("host"), app)
			}

			return localProcessStop(log, app)
		},
	}
}

// --- Local execution (no SSH) ---

func localProcessStart(log logger.Logger, appName, bin string, port int, args []string) error {
	repo := registry.NewFileRepository(appsPath)

	// reuse saved config if app already registered
	saved, _ := repo.Find(appName)
	if saved != nil {
		if bin == "" {
			bin = saved.Bin
		}
		if port == 0 {
			port = saved.Port
		}
		if len(args) == 0 {
			args = saved.Args
		}
		// check if already running
		if saved.Status == registry.AppStateStatusRunning {
			proc, err := os.FindProcess(saved.LastPID)
			if err == nil && proc.Signal(syscall.Signal(0)) == nil {
				log.Info("app %s already running with PID %d", appName, saved.LastPID)
				return nil
			}
		}
	}

	binPath, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("binary %q not found in PATH", bin)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", appName, err)
	}

	if err := repo.Save(registry.AppState{
		Name:    appName,
		Bin:     bin,
		Args:    args,
		Port:    port,
		LastPID: cmd.Process.Pid,
		Status:  registry.AppStateStatusRunning,
	}); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to save app state: %w", err)
	}

	log.Info("started %s with PID %d on PORT %d", appName, cmd.Process.Pid, port)
	return nil
}

func localProcessStatus(log logger.Logger, appName string) error {
	repo := registry.NewFileRepository(appsPath)
	app, err := repo.Find(appName)
	if err != nil {
		log.Info("app %s not found", appName)
		return nil
	}

	status := app.Status
	proc, err := os.FindProcess(app.LastPID)
	if err == nil && proc.Signal(syscall.Signal(0)) == nil {
		status = registry.AppStateStatusRunning
	} else {
		status = registry.AppStateStatusStopped
	}

	log.Info("app %s is %s (PID: %d, PORT: %d)", app.Name, status, app.LastPID, app.Port)
	return nil
}

func localProcessStop(log logger.Logger, appName string) error {
	repo := registry.NewFileRepository(appsPath)
	app, err := repo.Find(appName)
	if err != nil || app == nil {
		log.Info("app %s not found", appName)
		return nil
	}

	proc, err := os.FindProcess(app.LastPID)
	if err != nil {
		return fmt.Errorf("process not found: %w", err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to stop %s: %w", appName, err)
	}

	app.Status = registry.AppStateStatusStopped
	_ = repo.Save(*app)

	log.Info("stopped %s (PID %d)", appName, app.LastPID)
	return nil
}

// --- Remote execution (SSH) ---

func remoteExecutor(cfgFile, hostName string) (executor.Executor, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}
	if hostName == "" {
		return nil, fmt.Errorf("--host is required when using --config")
	}
	remote, ok := cfg.Config.Remote[hostName]
	if !ok {
		return nil, fmt.Errorf("host %q not found in config", hostName)
	}
	return executor.NewExecutorFromConfig(hostName, &remote)
}

func remoteProcessStart(cfgFile, hostName, appName, bin string, port int, args []string) error {
	exec, err := remoteExecutor(cfgFile, hostName)
	if err != nil {
		return err
	}

	logFile := fmt.Sprintf("/var/log/gosong/%s.log", appName)
	pidFile := fmt.Sprintf("/var/lib/gosong/apps/%s.pid", appName)
	cmd := strings.Join(append([]string{bin}, args...), " ")
	script := fmt.Sprintf(
		"mkdir -p /var/log/gosong /var/lib/gosong/apps && PORT=%d nohup %s > %s 2>&1 & echo $! > %s && echo $!",
		port, cmd, logFile, pidFile,
	)

	out, err := exec.RunRaw(script)
	if err != nil {
		return fmt.Errorf("failed to start %s on %s: %w", appName, hostName, err)
	}
	logger.Info("[%s] started %s on port %d (PID %s)", hostName, appName, port, strings.TrimSpace(out))
	return nil
}

func remoteProcessStatus(cfgFile, hostName, appName string) error {
	exec, err := remoteExecutor(cfgFile, hostName)
	if err != nil {
		return err
	}

	pidFile := fmt.Sprintf("/var/lib/gosong/apps/%s.pid", appName)
	script := fmt.Sprintf(
		"if [ -f %s ]; then pid=$(cat %s); kill -0 $pid 2>/dev/null && echo \"running (PID $pid)\" || echo stopped; else echo 'not found'; fi",
		pidFile, pidFile,
	)

	out, err := exec.RunRaw(script)
	if err != nil {
		return err
	}
	logger.Info("[%s] %s: %s", hostName, appName, strings.TrimSpace(out))
	return nil
}

func remoteProcessStop(cfgFile, hostName, appName string) error {
	exec, err := remoteExecutor(cfgFile, hostName)
	if err != nil {
		return err
	}

	pidFile := fmt.Sprintf("/var/lib/gosong/apps/%s.pid", appName)
	script := fmt.Sprintf(
		"if [ -f %s ]; then kill $(cat %s) && rm %s && echo stopped; else echo 'not running'; fi",
		pidFile, pidFile, pidFile,
	)

	out, err := exec.RunRaw(script)
	if err != nil {
		return err
	}
	logger.Info("[%s] %s: %s", hostName, appName, strings.TrimSpace(out))
	return nil
}
