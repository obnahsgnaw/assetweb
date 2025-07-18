package version

import (
	"fmt"
	"runtime"
)

var (
	name = "-"
	// The version is in string format and follow the ([v|V]Major.Minor.Patch[-Prerelease][+BuildMetadata])
	version = "v0.0.0"
	// gitCommit is the git sha1
	gitCommit = "-"
	// gitTreeState is the state of the git tree
	gitTreeState = "-"
	// buildDate
	buildDate = "-"
)

// BuildInfo describes to compile time information.
type BuildInfo struct {
	// Version is the current semver.
	Name string `json:"name,omitempty"`
	// Version is the current semver.
	Version string `json:"version,omitempty"`
	// GitCommit is the git sha1.
	GitCommit string `json:"gitCommit,omitempty"`
	// GitTreeState is the state of the git tree.
	// It is either clean or dirty.
	GitTreeState string `json:"gitTreeState,omitempty"`
	// BuildDate is the build date.
	BuildDate string `json:"buildDate,omitempty"`
	// GoVersion is the version of the Go compiler used.
	GoVersion string `json:"goVersion,omitempty"`
	// Compiler is the go compiler that built.
	Compiler string `json:"compiler,omitempty"`
	// Platform is the OS on which it is running.
	Platform string `json:"platform,omitempty"`
}

func (b BuildInfo) String() string {
	return fmt.Sprintf(
		"%s version: %s \nBuild Date: %s, Go Version: %s, Compiler: %s, Platform: %s \n",
		b.Name,
		b.Version,
		b.BuildDate,
		b.GoVersion,
		b.Compiler,
		b.Platform,
	)
}

// Version returns the version
func Version() string {
	return version
}

// Info returns build info
func Info() BuildInfo {
	v := BuildInfo{
		Name:         name,
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     runtime.GOOS + "/" + runtime.GOARCH,
	}

	return v
}
