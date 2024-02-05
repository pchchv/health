package bugsnag

type notifyHandler struct {
	PayloadChan chan *payload
}
