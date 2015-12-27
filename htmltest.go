package main

import (
	"io/ioutil"
	"testing"

	"github.com/rtfb/go-html-transform/css/selector"
	"golang.org/x/net/html"
)

func CssSelect(t *testing.T, node *html.Node, query string) []*html.Node {
	chain, err := selector.Selector(query)
	if err != nil {
		t.Fatalf("Error: query=%q, err=%s", query, err.Error())
	}
	return chain.Find(node)
}

func curlParam(url string, method func(string) (*http.Response, error)) string {
	if r, err := method(url); err == nil {
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err == nil {
			return string(b)
		}
		println(err.Error())
	} else {
		println(err.Error())
	}
	return ""
}

func Curl(url string) string {
	return curlParam(url, tclientGet)
}

func CurlPost(url string) string {
	return curlParam(url, tclientPostForm)
}
