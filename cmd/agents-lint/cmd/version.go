package cmd

// version is the tool's version string (FR-CLI-06). It defaults to "dev"
// for local `go build`/`go run`; a release build overrides it at link time
// via `-ldflags "-X .../cmd.version=vX.Y.Z"` (packaging-and-docs).
var version = "dev"

func init() {
	rootCmd.Version = version
}
