package tasks

func init() {
	RegisterTask("start-app", Task{
		Description: "Start the application on the green port (blue-green deployment)",
		Steps: []Step{
			{Run: "mkdir -p {{.AppPath}}/logs {{.AppPath}}/state"},
			{Run: "PORT={{.GreenPort}} nohup {{.App.Start}} > {{.AppPath}}/logs/{{.App.Name}}.log 2>&1 & echo $! > {{.AppPath}}/state/{{.App.Name}}.new.pid"},
			{Run: "echo 'Started {{.App.Name}} on port {{.GreenPort}}'"},
		},
	})

	RegisterTask("stop-old-app", Task{
		Description: "Stop the previous app instance and persist new active port state",
		Steps: []Step{
			{Run: "[ -f {{.AppPath}}/state/{{.App.Name}}.pid ] && kill $(cat {{.AppPath}}/state/{{.App.Name}}.pid) 2>/dev/null || true"},
			{Run: "mv {{.AppPath}}/state/{{.App.Name}}.new.pid {{.AppPath}}/state/{{.App.Name}}.pid"},
			{Run: "echo {{.GreenPort}} > {{.AppPath}}/state/{{.App.Name}}.current_port"},
			{Run: "echo 'Stopped old {{.App.Name}} instance (was on port {{.ActivePort}})'"},
		},
	})
}
