package version

import (
	"fmt"
	"runtime"
)

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func Info() string {

	return fmt.Sprintf(
		"Version: %s\nCommit: %s\nBuilt:   %s\nRuntime: %s on %s/%s",
		version,
		commit,
		buildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
