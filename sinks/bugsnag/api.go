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

type request struct {
	Url        string `json:"url"`
	Parameters string `json:"parameters"`
}

type payloadFrame struct {
	File       string `json:"file"`
	LineNumber int    `json:"lineNumber"`
	Method     string `json:"method"`
	InProject  bool   `json:"inProject"`
}

type payloadException struct {
	ErrorClass string         `json:"errorClass"`
	Message    string         `json:"message"`
	Stacktrace []payloadFrame `json:"stacktrace"`
}

type payloadEvent struct {
	PayloadVersion string             `json:"payloadVersion"`
	Exceptions     []payloadException `json:"exceptions"`
	// threads
	Context string `json:"context"`
	// groupingHash
	// severity
	// user
	App struct {
		// version
		ReleaseStage string `json:"releaseStage"`
	} `json:"app"`
	Device struct {
		//osVersion
		Hostname string `json:"hostname"`
	} `json:"device"`
	// meta data
	Metadata struct {
		Request request           `json:"request"`
		Kvs     map[string]string `json:"kvs"`
	} `json:"metaData"`
}

type payload struct {
	APIKey   string         `json:"apiKey"`
	Events   []payloadEvent `json:"events"`
	Notifier struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		URL     string `json:"url"`
	} `json:"notifier"`
}
