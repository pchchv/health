package bugsnag

type Config struct {
	// Your Bugsnag API key, e.g. "c9d60ae4c7e70c4b6c4ebd3e8056d2b8".
	APIKey string
	// The Endpoint to notify about crashes.
	Endpoint string
	// The current release stage.
	ReleaseStage string
	// The currently running version of the app.
	// This is used to filter errors in the Bugsnag dasboard.
	AppVersion string
	// The hostname of the current server. This defaults to the return value of
	// os.Hostname() and is graphed in the Bugsnag dashboard.
	Hostname string
}
