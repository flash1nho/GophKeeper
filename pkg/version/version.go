package version

var (
	buildVersion string = "N/A"
	buildCommit  string = "N/A"
	buildDate    string = "N/A"
)

func Info() string {
	return buildVersion + " (commit: " + buildCommit + ", date: " + buildDate + ")"
}
