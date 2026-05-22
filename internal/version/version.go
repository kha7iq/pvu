package version

import (
	"fmt"
	"runtime"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func Info() string {
	return fmt.Sprintf("version=%s commit=%s buildDate=%s go=%s %s/%s",
		Version,
		Commit,
		BuildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}

func Short() string {
	return Version
}
