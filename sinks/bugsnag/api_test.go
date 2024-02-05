package bugsnag

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type notifyHandler struct {
	PayloadChan chan *payload
}

func (h notifyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(rw, "got error in ready body: %v", err)
		return
	}

	var resp payload
	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Fprintf(rw, "got error in unmarshal: %v", err)
		return
	}

	h.PayloadChan <- &resp

	fmt.Fprintf(rw, "OK")
}
