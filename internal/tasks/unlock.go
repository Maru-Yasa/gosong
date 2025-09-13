package tasks

func init() {
	RegisterTask("lock", Task{
		Description: "Lock deployment to prevent concurrent runs",
		Steps: []Step{
			{Run: "if [ -f {{.AppPath}}/.gosong.lock ]; then echo '[lock] Already locked!'; exit 1; fi"},
			{Run: "echo 'locked' > {{.AppPath}}/.gosong.lock"},
			{Run: "echo 'Deployment locked succesfully'"},
		},
	})
}
