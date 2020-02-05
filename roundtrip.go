package agilecrm

import "net/http"

type rt struct {
	user string
	pass string
	orig http.RoundTripper
}

// RoundTrip ...
func (r rt) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(r.user, r.pass)

	return r.orig.RoundTrip(req)
}
