package daikin

import (
	"bytes"
	"net"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gorilla/schema"
)

type BasicInfo struct {
	Ret         string
	Type        string
	Reg         string
	Dst         int
	Ver         string
	Pow         int
	Err         int
	Location    int
	Name        string
	Icon        int
	Method      string
	Port        int
	Id          string
	Pw          string
	LpwFlag     int
	AdpKind     int
	Pv          int
	Cpv         int
	Led         int
	EnSetzone   int
	Mac         string
	AdpMode     string
	UDPAddress  net.UDPAddr
	HTTPAddress net.TCPAddr
}

func mapBytes(target interface{}, body []byte) error {
	vals, err := url.ParseQuery(strings.Replace(string(body), ",", "&", -1))
	if err != nil {
		return err
	}
	return mapValues(target, vals)
}

func mapValues(target interface{}, vals map[string][]string) error {

	vals2 := map[string][]string{}

	for k, v := range vals {
		vals2[camel(k)] = v
	}

	decoder := schema.NewDecoder()
	decoder.Decode(target, vals2)

	return nil
}

var camelingRegex = regexp.MustCompile("[0-9A-Za-z]+")

func camel(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		}
	}
	return upperFirst(string(bytes.Join(chunks, nil)))
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
