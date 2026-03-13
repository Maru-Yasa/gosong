package tasks

func init() {
	RegisterTask("start-app", Task{
		Description: "Start the application on the green port (blue-green deployment)",
		Steps: []Step{
			{Run: "mkdir -p /var/log/gosong /var/lib/gosong/apps"},
			{Run: "PORT={{.GreenPort}} nohup {{.App.Start}} > /var/log/gosong/{{.App.Name}}.log 2>&1 & echo $! > /var/lib/gosong/apps/{{.App.Name}}.new.pid"},
			{Run: "echo 'Started {{.App.Name}} on port {{.GreenPort}}'"},
		},
	})

	RegisterTask("stop-old-app", Task{
		Description: "Stop the previous app instance and persist new active port state",
		Steps: []Step{
			{Run: "[ -f /var/lib/gosong/apps/{{.App.Name}}.pid ] && kill $(cat /var/lib/gosong/apps/{{.App.Name}}.pid) 2>/dev/null || true"},
			{Run: "mv /var/lib/gosong/apps/{{.App.Name}}.new.pid /var/lib/gosong/apps/{{.App.Name}}.pid"},
			{Run: "echo {{.GreenPort}} > /var/lib/gosong/apps/{{.App.Name}}.current_port"},
			{Run: "echo 'Stopped old {{.App.Name}} instance (was on port {{.ActivePort}})'"},
		},
	})
}
