package healthd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/braintree/manners"
)

// serveJson will start a server on the hostPort and serve any path the Jsonified data.
// Each successive HTTP request will return the next data.
// If there is only one data, it will be returned on each request.
func serveJson(hostPort string, data ...interface{}) func() bool {
	var curData = 0
	var f http.HandlerFunc
	f = func(rw http.ResponseWriter, r *http.Request) {
		d := data[curData]
		curData = (curData + 1) % len(data)
		jsonData, err := json.MarshalIndent(d, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(rw, string(jsonData))
	}

	go manners.ListenAndServe(hostPort, f)
	time.Sleep(10 * time.Millisecond)

	return manners.Close
}
