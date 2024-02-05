package bugsnag

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/pchchv/health/stack"
)

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

// Notify will send the error and stack trace to Bugsnag. Note that this doesn't take advantage of all of Bugsnag's capabilities.
func Notify(config *Config, jobName string, eventName string, err error, trace *stack.Trace, kvs map[string]string) error {
	// Make a struct that serializes to the JSON needed for the API request to bugsnag
	p := newPayload(config, jobName, eventName, err, trace, kvs)
	// JSON serialize it
	data, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}

	// Post it to the server:
	client := http.Client{}
	resp, err := client.Post(config.Endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "OK" {
		return errors.New("response from bugsnag wasn't 'OK'")
	}

	return nil
}

func newPayload(config *Config, jobName string, eventName string, err error, trace *stack.Trace, kvs map[string]string) *payload {
	except := payloadException{
		ErrorClass: eventName,
		Message:    err.Error(),
	}
	for _, frame := range trace.Frames() {
		pf := payloadFrame{
			File:       frame.File,
			LineNumber: frame.LineNumber,
			Method:     frame.Package + ":" + frame.Name,
			InProject:  !frame.IsSystemPackage,
		}
		except.Stacktrace = append(except.Stacktrace, pf)
	}

	evt := payloadEvent{
		PayloadVersion: "2",
		Exceptions:     []payloadException{except},
		Context:        jobName,
	}
	evt.App.ReleaseStage = config.ReleaseStage
	evt.Device.Hostname = config.Hostname
	evt.Metadata.Kvs = kvs

	if requestUrl, requestUrlExists := kvs["request"]; requestUrlExists {
		evt.Metadata.Request.Url = requestUrl
	}

	if formData, formDataExists := kvs["formdata"]; formDataExists {
		evt.Metadata.Request.Parameters = formData
	}

	p := payload{
		APIKey: config.APIKey,
		Events: []payloadEvent{evt},
	}
	p.Notifier.Name = "health"
	p.Notifier.Version = "1.0"
	p.Notifier.URL = "https://www.github.com/gocraft/health"

	return &p
}
