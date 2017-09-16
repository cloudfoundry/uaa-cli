package version

var Version string
var Commit string

func VersionString() string {
	return Version + " " + Commit
}
