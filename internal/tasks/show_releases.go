package tasks

func init() {
	var script string = `
		current_id=$(readlink -f {{.AppPath}}/current 2>/dev/null | xargs basename)

		ls -1 {{.AppPath}}/releases | grep -E '^[0-9]+$' | sort -n -r | while read id; do
			if [ "$id" = "$current_id" ]; then
				echo "$id (current)"
			else
				echo "$id"
			fi
		done
	`

	RegisterTask("show_releases", Task{
		Description: "Show releases",
		Steps: []Step{
			{Run: script},
		},
	})
}
