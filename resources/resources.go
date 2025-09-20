package resoure

import _ "embed"

//go:embed systemd/gosong.service
var SystemdConfig string

//go:embed scripts/init.sh
var InitScript string
