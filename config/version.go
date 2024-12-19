package config

import _ "embed"

//go:generate sh -c "printf %s $(git rev-parse HEAD) > VERSION.txt"
//go:embed VERSION.txt
var version string
