package daikin

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func post(host, path string, values url.Values) error {
	resp, err := http.PostForm("http://"+host+path, values)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Post to %s failed. Code:%d Message:%s", host+path, resp.StatusCode, string(bs))
	}

	return nil
}

func get(host, path string) (url.Values, error) {
	resp, err := http.Get("http://" + host + path)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("Post to %s failed. Code:%d Message:%s", host+path, resp.StatusCode, string(bs))
		return nil, err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bs)

	return url.ParseQuery(strings.Replace(body, ",", "&", -1))
}
