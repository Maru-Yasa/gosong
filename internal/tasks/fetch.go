package tasks

func init() {
	RegisterTask("fetch", Task{
		Description: "Fetch teh application source",
		Steps: []Step{
			{Run: "echo 'Fetching from {{.Source.Type}} -> {{.Source.Url}}'"},
			{Run: "mkdir -p {{.ReleasePath}}"},
			{Run: "git clone --progress --verbose -b {{.Source.Branch}} --depth 1 {{.Source.Url}} {{.ReleasePath}}"},
			{Run: "rm -rf {{.ReleasePath}}/.git"},
		},
	})
}
