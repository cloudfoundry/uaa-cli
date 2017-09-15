package version

var Version string
var Commit string

func VersionString() string {
	return "UAA CLI " + Version + " " + Commit
}
