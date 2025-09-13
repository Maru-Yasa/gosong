package tasks

func init() {
	RegisterTask("symlink", Task{
		Description: "Update symlink to the latest release",
		Steps: []Step{
			{Run: "rm -f {{.AppPath}}/current"},
			{Run: "ln -sfn {{.ReleasePath}} {{.AppPath}}/current"},
			{Run: "[[ $(readlink -f {{.AppPath}}/current) == {{.ReleasePath}} ]] || { echo 'symlink error!'; exit 1; }"},
			{Run: "echo 'symlink current -> {{.ReleasePath}}'"},
		},
	})
}
