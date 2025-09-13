package tasks

func init() {
	RegisterTask("fetch", Task{
		Description: "Fetch teh application source",
		Steps: []Step{
			{Run: "echo 'Fetching from {{.Source.Type}} -> {{.Source.Url}}'"},
			{Run: "rm -rf {{.ReleasePath}}"},
			{Run: "git clone --progress --verbose -b {{.Source.Branch}} --depth 1 {{.Source.Url}} {{.ReleasePath}}"},
		},
	})
}
