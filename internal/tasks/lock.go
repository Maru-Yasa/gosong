package tasks

func init() {
	RegisterTask("unlock", Task{
		Description: "Unlock deployment",
		Steps: []Step{
			{Run: "rm -f {{.AppPath}}/.gosong.lock"},
			{Run: "echo 'Deployment unlocked'"},
		},
	})
}
