package daikin

import "github.com/davecgh/go-spew/spew"

func post(host, path string, query map[string]string) error {
	qs := ""
	for k, v := range query {
		qs += "&" + k + "=" + v
	}
	qs = qs[1:]

	spew.Dump("POST", query, qs)

	return nil
}
